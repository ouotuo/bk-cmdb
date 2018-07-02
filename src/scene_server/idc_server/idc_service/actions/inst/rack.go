/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by idclicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package inst

import (
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/validator"
	"configcenter/src/source_controller/api/auditlog"
	"configcenter/src/source_controller/api/metadata"
	api "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

var rack = &rackAction{}

type rackAction struct {
	base.BaseAction
}

func init() {

	// init action
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/rack/{idc_id}", Params: nil, Handler: rack.CreateRack})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/rack/{idc_id}/{rack_id}", Params: nil, Handler: rack.DeleteRack})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/rack/{idc_id}/{rack_id}", Params: nil, Handler: rack.UpdateRack})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/rack/search/{owner_id}/{idc_id}", Params: nil, Handler: rack.SearchRack})

	// rack cc interface
	rack.CreateAction()
}

// CreateModule
func (cli *rackAction) CreateRack(req *restful.Request, resp *restful.Response) {

	blog.Debug("create rack")

	// get language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	user := util.GetActionUser(req)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		//create default module
		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%v", err)
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		js, err := simplejson.NewJson(value)
		if nil != err {
			blog.Error("failed to unmarshal the data , error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		input, jsonErr := js.Map()
		if nil != jsonErr {
			blog.Error("failed to unmarshal the data , error info is %s", jsonErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		_, rackOK := input[common.BKRackNameField]
		if !rackOK {
			blog.Errorf("not rack '%s'", common.BKRackNameField)
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKRackNameField)
		}

		ownerID, ownerOK := input[common.BKOwnerIDField]
		if !ownerOK {
			blog.Error("'%s' field must be rackted", common.BKOwnerIDField)
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKOwnerIDField)
		}
		//
		//_, parentOK := input[common.BKInstParentStr]
		//if !parentOK {
		//	blog.Error("'%s' field must be rackted", common.BKInstParentStr)
		//	return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKInstParentStr)
		//}

		idcID, convErr := strconv.Atoi(req.PathParameter("idc_id"))
		if nil != convErr {
			blog.Error("failed to convert the idcid to int, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "idc_id")
		}

		tmpID, ok := ownerID.(string)
		if !ok {
			blog.Error("'OwnerID' must be a string value")
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKOwnerIDField)
		}

		input[common.BKIdcIDField] = idcID
		// check
		valid := validator.NewValidMapWithKeyFields(tmpID, common.BKInnerObjIDRack, cli.CC.ObjCtrl(), []string{common.BKInstParentStr, common.BKOwnerIDField}, forward, defErr)
		_, err = valid.ValidMap(input, common.ValidCreate, 0)
		if nil != err {
			blog.Error("failed to valid the input data, error info is %s", err.Error())
			return http.StatusBadRequest, "", err
		}

		// create
		input[common.BKDefaultField] = 0
		input[common.CreateTimeField] = util.GetCurrentTimeStr()

		inputJSON, jsErr := json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		cModuleURL := cli.CC.ObjCtrl() + "/object/v1/insts/rack"
		moduleRes, err := httpcli.ReqHttp(req, cModuleURL, common.HTTPCreate, inputJSON)
		if nil != err {
			blog.Error("failed to create the rack, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcRackCreateFailed)
		}

		{
			// save change log
			instID := gjson.Get(moduleRes, "data."+common.BKRackIDField).Int()
			ownerID := fmt.Sprint(input[common.BKOwnerIDField])
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDRack)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcRackCreateFailed)
			}

			curData, retStrErr := inst.getInstDetail(req, int(instID), common.BKInnerObjIDRack, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcRackCreateFailed)
			}
			auditContent := metadata.Content{
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditRackLog(instID, auditContent, "create rack", ownerID, fmt.Sprint(idcID), user, auditoplog.AuditOpTypeAdd)
		}

		return http.StatusOK, moduleRes, nil
	}, resp)

}

// DeleteRack delete rack by conditions
func (cli *rackAction) DeleteRack(req *restful.Request, resp *restful.Response) {

	blog.Debug("delete rack")

	// get language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	user := util.GetActionUser(req)
	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		idcID, convErr := strconv.Atoi(req.PathParameter("idc_id"))
		if nil != convErr {
			blog.Error("the idcid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "idc_id")
		}

		rackID, convErr := strconv.Atoi(req.PathParameter("rack_id"))
		if nil != convErr {
			blog.Error("the rackid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "rack_id")
		}

		// check wether it can be delete
		rstOk, rstErr := hasHost(req, cli.CC.HostCtrl(), map[string][]int{
			common.BKIdcIDField: []int{idcID},
			common.BKRackIDField: []int{rackID},
		})
		if nil != rstErr {
			blog.Error("failed to check rack wether it has hosts, error info is %s", rstErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrTopoHasHostCheckFailed)
		}

		if !rstOk {
			blog.Error("failed to delete rack, because of it has some hosts")
			return http.StatusBadRequest, "", defErr.Error(common.CCErrTopoHasHost)
		}

		// take snapshot before operation
		ownerID := idc.getOwnerIDByIdcID(req, idcID)
		if ownerID == "" {
			blog.Errorf("owner id not found")
		}
		preData, retStrErr := inst.getInstDetail(req, rackID, common.BKInnerObjIDRack, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		//delete rack
		input := make(map[string]interface{})
		input[common.BKIdcIDField] = idcID
		input[common.BKRackIDField] = rackID

		uURL := cli.CC.ObjCtrl() + "/object/v1/insts/rack"

		inputJSON, jsErr := json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		_, err := httpcli.ReqHttp(req, uURL, common.HTTPDelete, []byte(inputJSON))
		if nil != err {
			blog.Error("failed to delete the rack, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcRackDeleteFailed)
		}

		//delete module
		input = make(map[string]interface{})
		input[common.BKIdcIDField] = idcID
		input[common.BKRackIDField] = rackID

		uURL = cli.CC.ObjCtrl() + "/object/v1/insts/pos"
		inputJSON, jsErr = json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		moduleRes, err := httpcli.ReqHttp(req, uURL, common.HTTPDelete, []byte(inputJSON))
		if nil != err {
			blog.Error("failed to delete the module, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoModuleDeleteFailed)
		}

		{
			// save change log
			instID := gjson.Get(moduleRes, "data.bk_rack_id").Int()
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDRack)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcRackDeleteFailed)
			}

			auditContent := metadata.Content{
				PreData: preData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditRackLog(instID, auditContent, "delete rack", ownerID, fmt.Sprint(idcID), user, auditoplog.AuditOpTypeDel)
		}
		return http.StatusOK, moduleRes, nil
	}, resp)
}

// UpdateRack update rack by condition
func (cli *rackAction) UpdateRack(req *restful.Request, resp *restful.Response) {
	blog.Debug("updte rack")

	// get language
	language := util.GetActionLanguage(req)

	// get error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	user := util.GetActionUser(req)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		idcID, convErr := strconv.Atoi(req.PathParameter("idc_id"))
		if nil != convErr {
			blog.Error("the idcid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "idc_id")
		}

		rackID, convErr := strconv.Atoi(req.PathParameter("rack_id"))
		if nil != convErr {
			blog.Error("the rackid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "rack_id")
		}

		//update rack
		input := make(map[string]interface{})
		condition := make(map[string]interface{})
		condition[common.BKIdcIDField] = idcID
		condition[common.BKRackIDField] = rackID

		value, readErr := ioutil.ReadAll(req.Request.Body)
		if nil != readErr {
			blog.Error("read request body failed, error:%s", readErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		js, jsErr := simplejson.NewJson([]byte(value))
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		data, jsErr := js.Map()
		data[common.BKIdcIDField] = idcID
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}
		valid := validator.NewValidMapWithKeyFields(common.BKDefaultOwnerID, common.BKInnerObjIDRack, cli.CC.ObjCtrl(), []string{common.BKInstParentStr, common.BKOwnerIDField, common.BKRackNameField}, forward, defErr)
		_, err := valid.ValidMap(data, common.ValidUpdate, rackID)
		if nil != err {
			blog.Error("failed to valid the input data, error info is %s", err.Error())
			return http.StatusBadRequest, "", err
		}

		// take snapshot before operation
		ownerID := idc.getOwnerIDByIdcID(req, idcID)
		if ownerID == "" {
			blog.Errorf("owner id not found")
		}
		preData, retStrErr := inst.getInstDetail(req, rackID, common.BKInnerObjIDRack, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		input["condition"] = condition
		input["data"] = data

		uURL := cli.CC.ObjCtrl() + "/object/v1/insts/rack"

		inputJSON, jsErr := json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		moduleRes, err := httpcli.ReqHttp(req, uURL, "PUT", []byte(inputJSON))
		if nil != err {
			blog.Error("failed to delete the rack, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcRackUpdateFailed)
		}

		{
			// save change log
			instID := rackID //gjson.Get(moduleRes, "data.bk_rack_id").Int()
			//ownerID := fmt.Sprint(input[common.BKOwnerIDField])
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDRack)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcRackUpdateFailed)
			}

			curData, retStrErr := inst.getInstDetail(req, int(instID), common.BKInnerObjIDRack, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcRackUpdateFailed)
			}
			auditContent := metadata.Content{
				PreData: preData,
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditRackLog(instID, auditContent, "update rack", ownerID, fmt.Sprint(idcID), user, auditoplog.AuditOpTypeModify)
		}
		return http.StatusOK, moduleRes, nil
	}, resp)

}

// SearfhPost
func (cli *rackAction) SearchRack(req *restful.Request, resp *restful.Response) {
	blog.Debug("search rack")
	// get language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {

		ownerID := req.PathParameter("owner_id")
		idcID, convErr := strconv.Atoi(req.PathParameter("idc_id"))
		if nil != convErr {
			blog.Error("the idcid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "idc_id")
		}

		value, readErr := ioutil.ReadAll(req.Request.Body)
		if nil != readErr {
			blog.Error("read request body failed, error:%s", readErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		var js params.SearchParams
		err := json.Unmarshal([]byte(value), &js)
		if nil != err {
			blog.Error("failed to unmarshal the data , error info is %s", err.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		condition := params.ParseAppSearchParams(js.Condition)

		condition[common.BKIdcIDField] = idcID

		page := js.Page

		searchParams := make(map[string]interface{})
		searchParams["condition"] = condition
		searchParams["fields"] = strings.Join(js.Fields, ",")
		searchParams["start"] = page["start"]
		searchParams["limit"] = page["limit"]
		searchParams["sort"] = page["sort"]

		//search
		sURL := cli.CC.ObjCtrl() + "/object/v1/insts/rack/search"
		inputJSON, _ := json.Marshal(searchParams)
		moduleRes, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, []byte(inputJSON))
		if nil != err {
			blog.Error("failed to select the rack, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcRackSelectFailed)
		}

		// replace the association id into name
		retStr, retStrErr := inst.getInstDetails(req, common.BKInnerObjIDRack, ownerID, moduleRes, map[string]interface{}{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  "",
		})
		if common.CCSuccess != retStrErr {
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		return http.StatusOK, retStr["data"], nil

	}, resp)

}
