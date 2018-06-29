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

var pos = &posAction{}

type posAction struct {
	base.BaseAction
}

func init() {

	// init action
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/pos/{idc_id}/{rack_id}", Params: nil, Handler: pos.CreatePos})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/pos/{idc_id}/{rack_id}/{pos_id}", Params: nil, Handler: pos.DeletePos})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/pos/{idc_id}/{rack_id}/{pos_id}", Params: nil, Handler: pos.UpdatePos})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/pos/search/{owner_id}/{idc_id}/{rack_id}", Params: nil, Handler: pos.SearchPos})

	// rack cc interface
	pos.CreateAction()
}

// CreatePos
func (cli *posAction) CreatePos(req *restful.Request, resp *restful.Response) {

	blog.Debug("create pos")

	// get language
	language := util.GetActionLanguage(req)
	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}
	user := util.GetActionUser(req)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {

		//create default pos
		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		js, err := simplejson.NewJson(value)
		if nil != err {
			blog.Error("the input json is invalid, error info is %s", err.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		input, jsErr := js.Map()
		if nil != jsErr {
			blog.Error("the input json is invalid, error info is %s", jsErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		rackID, convErr := strconv.Atoi(req.PathParameter("rack_id"))
		if nil != convErr {
			blog.Error("the rackid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "rack_id")
		}

		idcID, convErr := strconv.Atoi(req.PathParameter("idc_id"))
		if nil != convErr {
			blog.Error("the idcid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "idc_id")
		}

		if _, ok := input[common.BKOwnerIDField]; !ok {
			blog.Error("not rack %s", common.BKOwnerIDField)
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKOwnerIDField)
		}

		if _, ok := input[common.BKPosNameField]; !ok {
			blog.Error("not rack PosName")
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKPosNameField)
		}

		//if _, ok := input[common.BKInstParentStr]; !ok {
		//	blog.Error("not rack %s", common.BKInstParentStr)
		//	return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKInstParentStr)
		//}

		tmpID, ok := input[common.BKOwnerIDField].(string)
		if !ok {
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedString, common.BKOwnerIDField)
		}

		// create
		input[common.BKRackIDField] = rackID
		input[common.BKIdcIDField] = idcID
		// check
		valid := validator.NewValidMapWithKeyFields(tmpID, common.BKInnerObjIDPos, cli.CC.ObjCtrl(), []string{common.BKOwnerIDField, common.BKInstParentStr}, forward, defErr)
		_, err = valid.ValidMap(input, common.ValidCreate, 0)
		if nil != err {
			blog.Error("failed to valide, error is %s", err.Error())
			return http.StatusBadRequest, "", err
		}

		// create
		input[common.BKDefaultField] = 0
		input[common.CreateTimeField] = util.GetCurrentTimeStr()

		inputJSON, jsErr := json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to marshal the json, error is info %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		cPosURL := cli.CC.ObjCtrl() + "/object/v1/insts/pos"

		posRes, err := httpcli.ReqHttp(req, cPosURL, common.HTTPCreate, inputJSON)
		if nil != err {
			blog.Error("failed to create the pos, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcPosCreateFailed)
		}

		{
			// save change log
			instID := gjson.Get(posRes, "data."+common.BKPosIDField).Int()
			if instID == 0 {
				blog.Errorf("inst id not found")
			}
			ownerID := idc.getOwnerIDByIdcID(req, idcID)
			if ownerID == "" {
				blog.Errorf("owner id not found")
			}
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDPos)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, "", defErr.Error(attErr)
			}

			curData, retStrErr := inst.getInstDetail(req, int(instID), common.BKInnerObjIDPos, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcPosCreateFailed)
			}
			auditContent := metadata.Content{
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditPosLog(instID, auditContent, "create pos", ownerID, fmt.Sprint(idcID), user, auditoplog.AuditOpTypeAdd)
		}
		return http.StatusOK, posRes, nil
	}, resp)

}

// DeletePos delete pos by condition
func (cli *posAction) DeletePos(req *restful.Request, resp *restful.Response) {

	blog.Debug("delete pos")

	// get the language
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

		posID, convErr := strconv.Atoi(req.PathParameter("pos_id"))
		if nil != convErr {
			blog.Error("the posid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "pos_id")
		}

		// check wether it can be delete
		rstOk, rstErr := hasHost(req, cli.CC.HostCtrl(), map[string][]int{
			common.BKIdcIDField:    []int{idcID},
			common.BKPosIDField: []int{posID},
			common.BKRackIDField:    []int{rackID},
		})
		if nil != rstErr {
			blog.Error("failed to check pos wether it has hosts, error info is %s", rstErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrTopoHasHostCheckFailed)
		}

		if !rstOk {
			blog.Error("failed to delete pos, because of it has some hosts")
			return http.StatusBadRequest, "", defErr.Error(common.CCErrTopoHasHost)
		}

		// take snapshot before operation
		ownerID := idc.getOwnerIDByIdcID(req, idcID)
		if ownerID == "" {
			blog.Errorf("owner id not found")
		}
		preData, retStrErr := inst.getInstDetail(req, posID, common.BKInnerObjIDPos, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		//delete pos
		input := make(map[string]interface{})
		input[common.BKIdcIDField] = idcID
		input[common.BKRackIDField] = rackID
		input[common.BKPosIDField] = posID

		uURL := cli.CC.ObjCtrl() + "/object/v1/insts/pos"

		inputJSON, jsErr := json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		posRes, err := httpcli.ReqHttp(req, uURL, "DELETE", []byte(inputJSON))
		if nil != err {
			blog.Error("failed to delete the pos, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcPosDeleteFailed)
		}

		{
			// save change log
			instID := gjson.Get(posRes, "data.bk_pos_id").Int()
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDPos)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcPosDeleteFailed)
			}
			auditContent := metadata.Content{
				PreData: preData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditPosLog(instID, auditContent, "delete pos", ownerID, fmt.Sprint(idcID), user, auditoplog.AuditOpTypeDel)
		}
		return http.StatusOK, posRes, nil

	}, resp)

}

// UpdatePos
func (cli *posAction) UpdatePos(req *restful.Request, resp *restful.Response) {
	blog.Debug("update pos")

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

		posID, _ := strconv.Atoi(req.PathParameter("pos_id"))
		if nil != convErr {
			blog.Error("the posid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "pos_id")
		}

		//update pos
		input := make(map[string]interface{})
		condition := make(map[string]interface{})
		condition[common.BKIdcIDField] = idcID
		condition[common.BKRackIDField] = rackID
		condition[common.BKPosIDField] = posID

		value, readErr := ioutil.ReadAll(req.Request.Body)
		if nil != readErr {
			blog.Error("failed to read the http request, error info is %s", readErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		js, err := simplejson.NewJson([]byte(value))
		if nil != err {
			blog.Error("failed to create simplejson, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		data, jsErr := js.Map()
		if nil != jsErr {
			blog.Error("failed to unmarshal data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		data[common.BKIdcIDField] = idcID
		data[common.BKRackIDField] = rackID
		valid := validator.NewValidMapWithKeyFields(common.BKDefaultOwnerID, common.BKInnerObjIDPos, cli.CC.ObjCtrl(), []string{common.BKOwnerIDField, common.BKInstParentStr, common.BKPosNameField}, forward, defErr)
		_, err = valid.ValidMap(data, common.ValidUpdate, posID)

		if nil != err {
			blog.Error("failed to valid the input , error is %s", err.Error())
			return http.StatusBadRequest, "", err
		}

		// take snapshot before operation
		ownerID := idc.getOwnerIDByIdcID(req, idcID)
		if ownerID == "" {
			blog.Errorf("owner id not found")
		}
		preData, retStrErr := inst.getInstDetail(req, posID, common.BKInnerObjIDPos, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		input["condition"] = condition
		input["data"] = data
		uURL := cli.CC.ObjCtrl() + "/object/v1/insts/pos"
		inputJSON, jsErr := json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		posRes, err := httpcli.ReqHttp(req, uURL, "PUT", []byte(inputJSON))

		if nil != err {
			blog.Error("failed to update the pos, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcPosUpdateFailed)
		}

		{
			// save change log
			instID := posID
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDPos)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcPosCreateFailed)
			}

			curData, retStrErr := inst.getInstDetail(req, instID, common.BKInnerObjIDPos, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcPosCreateFailed)
			}
			auditContent := metadata.Content{
				PreData: preData,
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditPosLog(instID, auditContent, "update pos", ownerID, fmt.Sprint(idcID), user, auditoplog.AuditOpTypeModify)
		}

		return http.StatusOK, posRes, nil

	}, resp)

	return

}

// SearfhPos search poss
func (cli *posAction) SearchPos(req *restful.Request, resp *restful.Response) {
	blog.Debug("search pos")

	// get the language
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
		rackID, convErr := strconv.Atoi(req.PathParameter("rack_id"))
		if nil != convErr {
			blog.Error("the rackid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "rack_id")
		}

		value, readErr := ioutil.ReadAll(req.Request.Body)
		if nil != readErr {
			blog.Error("failed to read the http request, error info is %s", readErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		var js params.SearchParams
		err := json.Unmarshal([]byte(value), &js)
		if nil != err {
			blog.Error("failed to unmarshal the input, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		condition := params.ParseAppSearchParams(js.Condition)

		condition[common.BKIdcIDField] = idcID
		condition[common.BKRackIDField] = rackID

		page := js.Page

		searchParams := make(map[string]interface{})
		searchParams["condition"] = condition
		searchParams["fields"] = strings.Join(js.Fields, ",")
		searchParams["start"] = page["start"]
		searchParams["limit"] = page["limit"]
		searchParams["sort"] = page["sort"]

		//search
		sURL := cli.CC.ObjCtrl() + "/object/v1/insts/pos/search"
		inputJSON, jsErr := json.Marshal(searchParams)
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		posRes, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, []byte(inputJSON))
		if nil != err {
			blog.Error("failed to update the pos, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrIdcPosSelectFailed)
		}

		// replace the association id to name
		retStr, retStrErr := inst.getInstDetails(req, common.BKInnerObjIDPos, ownerID, posRes, map[string]interface{}{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  "",
		})
		if common.CCSuccess != retStrErr {
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		return http.StatusOK, retStr["data"], nil
	}, resp)

	return
}
