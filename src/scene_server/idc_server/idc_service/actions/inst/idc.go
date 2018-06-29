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
	"configcenter/src/common/errors"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	sencecommon "configcenter/src/scene_server/common"
	"configcenter/src/scene_server/validator"
	"configcenter/src/source_controller/api/auditlog"
	"configcenter/src/source_controller/api/metadata"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tidwall/gjson"

	"io/ioutil"
	"strconv"
	"strings"

	api "configcenter/src/source_controller/api/object"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

var idc = &idcAction{}

type idcAction struct {
	base.BaseAction
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/idc/{owner_id}", Params: nil, Handler: idc.CreateIdc})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/idc/{owner_id}/{idc_id}", Params: nil, Handler: idc.DeleteIdc})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/idc/{owner_id}/{idc_id}", Params: nil, Handler: idc.UpdateIdc})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/idc/status/{flag}/{owner_id}/{idc_id}", Params: nil, Handler: idc.UpdateIdcDataStatus})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/idc/search/{owner_id}", Params: nil, Handler: idc.SearchIdc})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/idc/default/{owner_id}/search", Params: nil, Handler: idc.GetDefaultIdc})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/idc/default/{owner_id}", Params: nil, Handler: idc.CreateDefaultIdc})

	// create CC object
	idc.CreateAction()
}

//delete idclication
func (cli *idcAction) DeleteIdc(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get error code in language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {

		//new feature, idc not allow deletion
		blog.Error("idc not allow deletion")
		return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcDeleteFailed)

		forward := &api.ForwardParam{Header: req.Request.Header}
		pathParams := req.PathParameters()
		idcID, _ := strconv.Atoi(pathParams["idc_id"])
		ownerID, _ := pathParams["owner_id"]
		user := sencecommon.GetUserFromHeader(req)

		// check wether it can be delete
		rstOk, rstErr := hasHost(req, cli.CC.HostCtrl(), map[string][]int{common.BKIdcIDField: []int{idcID}})
		if nil != rstErr {
			blog.Error("failed to check idc wether it has hosts, error info is %s", rstErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoHasHostCheckFailed)
		}

		if !rstOk {
			blog.Error("failed to delete idc, because of it has some hosts")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoHasHostCheckFailed)
		}

		// take snapshot before operation
		preData, retStrErr := inst.getInstDetail(req, idcID, common.BKInnerObjIDIdc, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditTakeSnapshotFaile)
		}
		idcData, ok := preData.(map[string]interface{})
		if false == ok {
			blog.Error("failed to get idc detail")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcDeleteFailed)
		}
		idcNameI, ok := idcData[common.BKIdcNameField]
		if false == ok {
			blog.Error("failed to get idc detail")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcDeleteFailed)
		}
		bkIdcName, ok := idcNameI.(string)
		if false == ok {
			blog.Error("failed to get idc detail")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcDeleteFailed)
		}
		if "" == bkIdcName {
		}
		//delete idc
		input := make(map[string]interface{})
		input[common.BKIdcIDField] = idcID
		dIdcURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDIdc
		inputJSON, _ := json.Marshal(input)
		blog.Info("delete idc url:%s", dIdcURL)
		_, err := httpcli.ReqHttp(req, dIdcURL, common.HTTPDelete, []byte(inputJSON))
		if nil != err {
			blog.Error("delete idc error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcDeleteFailed)
		}
		{
			// save change log
			instID, _ := strconv.Atoi(fmt.Sprint(idcID))
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDIdc)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditSaveLogFaile)
			}

			auditContent := metadata.Content{
				PreData: preData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditObjLog(instID, auditContent, "delete idc", common.BKInnerObjIDIdc, ownerID, "0", user, auditoplog.AuditOpTypeDel)
		}
		//delete set in idc
		setInput := make(map[string]interface{})
		setInput[common.BKIdcIDField] = idcID
		inputSetJSON, _ := json.Marshal(setInput)
		dSetURL := cli.CC.ObjCtrl() + "/object/v1/insts/rack"
		_, err = httpcli.ReqHttp(req, dSetURL, common.HTTPDelete, []byte(inputSetJSON))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcDeleteFailed)
		}
		//delete module in idc
		moduleInput := make(map[string]interface{})
		moduleInput[common.BKIdcIDField] = idcID
		inputModuleJSON, _ := json.Marshal(moduleInput)
		dModuleURL := cli.CC.ObjCtrl() + "/object/v1/insts/pos"
		_, err = httpcli.ReqHttp(req, dModuleURL, common.HTTPDelete, []byte(inputModuleJSON))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoModuleDeleteFailed)
		}
		return http.StatusOK, nil, nil
	}, resp)
}

// update  idclication data status
func (cli *idcAction) UpdateIdcDataStatus(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get error code in language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		pathParams := req.PathParameters()
		idcID, _ := strconv.Atoi(pathParams["idc_id"])
		ownerID, _ := pathParams["owner_id"]
		flag, _ := pathParams["flag"]
		if flag != fmt.Sprint(common.DataStatusDisabled) && flag != fmt.Sprint(common.DataStatusEnable) {
			blog.Error("input params error:")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		user := sencecommon.GetUserFromHeader(req)
		//update idc
		input := make(map[string]interface{})
		data := make(map[string]interface{})
		condition := make(map[string]interface{})
		condition[common.BKIdcIDField] = idcID
		condition[common.BKOwnerIDField] = ownerID
		data[common.BKDataStatusField] = flag

		// take snapshot before operation
		preData, retStrErr := inst.getInstDetail(req, idcID, common.BKInnerObjIDIdc, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditTakeSnapshotFaile)
		}

		input["condition"] = condition
		input["data"] = data
		uIdcURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDIdc
		inputJSON, _ := json.Marshal(input)
		_, err := httpcli.ReqHttp(req, uIdcURL, common.HTTPUpdate, []byte(inputJSON))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcUpdateFailed)
		}

		{
			// save change log
			instID, _ := strconv.Atoi(fmt.Sprint(idcID))
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDIdc)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditTakeSnapshotFaile)
			}

			curData, retStrErr := inst.getInstDetail(req, instID, common.BKInnerObjIDIdc, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditSaveLogFaile)
			}

			auditContent := metadata.Content{
				PreData: preData,
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditObjLog(instID, auditContent, "update idc", common.BKInnerObjIDIdc, ownerID, "0", user, auditoplog.AuditOpTypeModify)
		}

		return http.StatusOK, nil, nil
	}, resp)

}

//update idclication
func (cli *idcAction) UpdateIdc(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get error code in language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		pathParams := req.PathParameters()
		idcID, _ := strconv.Atoi(pathParams["idc_id"])
		ownerID, _ := pathParams["owner_id"]
		user := sencecommon.GetUserFromHeader(req)
		//update idc
		input := make(map[string]interface{})
		condition := make(map[string]interface{})
		condition[common.BKIdcIDField] = idcID
		condition[common.BKOwnerIDField] = ownerID
		value, err := ioutil.ReadAll(req.Request.Body)
		js, err := simplejson.NewJson([]byte(value))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		data, _ := js.Map()
		valid := validator.NewValidMap(common.BKDefaultOwnerID, common.BKInnerObjIDIdc, cli.CC.ObjCtrl(), forward, defErr)
		_, err = valid.ValidMap(data, common.ValidUpdate, idcID)
		if nil != err {
			blog.Errorf("UpdateIdc vaild error:%s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommFieldNotValid)
		}

		// take snapshot before operation
		preData, retStrErr := inst.getInstDetail(req, idcID, common.BKInnerObjIDIdc, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditTakeSnapshotFaile)
		}

		idcData, ok := preData.(map[string]interface{})
		if false == ok {
			blog.Error("failed to get idc detail")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcUpdateFailed)
		}
		idcNameI, ok := idcData[common.BKIdcNameField]
		if false == ok {
			blog.Error("failed to get idc detail")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcUpdateFailed)
		}
		bkIdcName, ok := idcNameI.(string)
		if false == ok {
			blog.Error("failed to get idc detail")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcUpdateFailed)
		}
		if "" == bkIdcName {

		}

		input["condition"] = condition
		input["data"] = data
		uIdcURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDIdc
		inputJSON, _ := json.Marshal(input)
		_, err = httpcli.ReqHttp(req, uIdcURL, common.HTTPUpdate, []byte(inputJSON))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcUpdateFailed)
		}

		{
			// save change log
			instID, _ := strconv.Atoi(fmt.Sprint(idcID))
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDIdc)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditTakeSnapshotFaile)
			}

			curData, retStrErr := inst.getInstDetail(req, instID, common.BKInnerObjIDIdc, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditSaveLogFaile)
			}

			auditContent := metadata.Content{
				PreData: preData,
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditObjLog(instID, auditContent, "update idc", common.BKInnerObjIDIdc, ownerID, "0", user, auditoplog.AuditOpTypeModify)
		}

		return http.StatusOK, nil, nil
	}, resp)

}

func (cli *idcAction) getOwnerIDByIdcID(req *restful.Request, idcID int) (ownerID string) {
	condition := map[string]interface{}{}
	condition[common.BKIdcIDField] = idcID
	sIdcURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDIdc + "/search"
	inputJSON, _ := json.Marshal(map[string]interface{}{"condition": condition})
	idcInfo, err := httpcli.ReqHttp(req, sIdcURL, common.HTTPSelectPost, []byte(inputJSON))
	if nil != err {
		blog.Error("search idc error: %v", err)
		return
	}
	return gjson.Get(idcInfo, "data.info.0."+common.BKOwnerIDField).String()
}

//search idclication
func (cli *idcAction) SearchIdc(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID, _ := pathParams["owner_id"]
		value, _ := ioutil.ReadAll(req.Request.Body)
		var js params.SearchParams
		err := json.Unmarshal([]byte(value), &js)
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		var condition map[string]interface{}
		if 1 == js.Native {
			condition = js.Condition
		} else {
			condition = params.ParseAppSearchParams(js.Condition)
		}

	//	condition[common.BKOwnerIDField] = ownerID
		condition[common.BKDefaultField] = 0
		page := js.Page
		searchParams := make(map[string]interface{})
		searchParams["condition"] = condition
		searchParams["fields"] = strings.Join(js.Fields, ",")
		searchParams["start"] = page["start"]
		searchParams["limit"] = page["limit"]
		searchParams["sort"] = page["sort"]
		//search idc
		sIdcURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDIdc + "/search"
		inputJSON, _ := json.Marshal(searchParams)
		idcInfo, err := httpcli.ReqHttp(req, sIdcURL, common.HTTPSelectPost, []byte(inputJSON))
		blog.Debug("search idc params: %s", string(inputJSON))
		if nil != err {
			blog.Error("search idc error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcSearchFailed)
		}
		blog.Debug("search idc return %v", idcInfo)
		// replace the association id to name
		retstr, retStrErr := inst.getInstDetails(req, common.BKInnerObjIDIdc, ownerID, idcInfo, map[string]interface{}{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  "",
		})
		if common.CCSuccess != retStrErr {
			blog.Error("search idc error: %v", retStrErr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcSearchFailed)
		}
		blog.Info("search idc return %v", retstr)
		return http.StatusOK, retstr["data"], nil
	}, resp)
}

//create idclication
func (cli *idcAction) CreateIdc(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		pathParams := req.PathParameters()
		ownerID := pathParams["owner_id"]
		user := sencecommon.GetUserFromHeader(req)
		value, _ := ioutil.ReadAll(req.Request.Body)
		js, err := simplejson.NewJson([]byte(value))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		input, err := js.Map()
		valid := validator.NewValidMap(common.BKDefaultOwnerID, common.BKInnerObjIDIdc, cli.CC.ObjCtrl(), forward, defErr)
		_, err = valid.ValidMap(input, common.ValidCreate, 0)
		if nil != err {
			blog.Errorf("create idc valid eror:%s, data:%v", err.Error(), string(value))
			if _, ok := err.(errors.CCErrorCoder); ok {
				return http.StatusInternalServerError, nil, err
			}
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommFieldNotValid)
		}

		input[common.BKOwnerIDField] = ownerID
		input[common.BKDefaultField] = 0
		input[common.BKSupplierIDField] = common.BKDefaultSupplierID
		idcInfoJSON, _ := json.Marshal(input)
		cIdcURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDIdc
		cIdcRes, err := httpcli.ReqHttp(req, cIdcURL, common.HTTPCreate, []byte(idcInfoJSON))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcCreateFailed)
		}
		js, err = simplejson.NewJson([]byte(cIdcRes))
		idcResData, _ := js.Map()
		idcIDInfo := idcResData["data"].(map[string]interface{})
		idcID := idcIDInfo[common.BKIdcIDField]
		{
			// save change log
			instID, _ := strconv.Atoi(fmt.Sprint(idcID))
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDIdc)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditTakeSnapshotFaile)
			}

			curData, retStrErr := inst.getInstDetail(req, instID, common.BKInnerObjIDIdc, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditSaveLogFaile)
			}
			auditContent := metadata.Content{
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditObjLog(instID, auditContent, "create idc", common.BKInnerObjIDIdc, ownerID, "0", user, auditoplog.AuditOpTypeAdd)
		}
		//create default set
		inputSetInfo := make(map[string]interface{})
		inputSetInfo[common.BKIdcIDField] = idcID
		inputSetInfo[common.BKInstParentStr] = idcID
		inputSetInfo[common.BKSetNameField] = common.DefaultResSetName
		inputSetInfo[common.BKDefaultField] = common.DefaultResSetFlag
		inputSetInfo[common.BKOwnerIDField] = ownerID
		cSetURL := cli.CC.ObjCtrl() + "/object/v1/insts/set"
		setJSONData, _ := json.Marshal(inputSetInfo)
		cSetRes, err := httpcli.ReqHttp(req, cSetURL, common.HTTPCreate, []byte(setJSONData))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoSetCreateFailed)
		}
		//create default module
		js, err = simplejson.NewJson([]byte(cSetRes))
		setResData, _ := js.Map()
		setIDInfo := setResData["data"].(map[string]interface{})
		setID := setIDInfo[common.BKSetIDField]
		inputResModuleInfo := make(map[string]interface{})
		inputResModuleInfo[common.BKSetIDField] = setID
		inputResModuleInfo[common.BKInstParentStr] = setID
		inputResModuleInfo[common.BKIdcIDField] = idcID
		inputResModuleInfo[common.BKModuleNameField] = common.DefaultResModuleName
		inputResModuleInfo[common.BKDefaultField] = common.DefaultResModuleFlag
		inputResModuleInfo[common.BKOwnerIDField] = ownerID
		cModuleURL := cli.CC.ObjCtrl() + "/object/v1/insts/module"
		resModuleJSONData, _ := json.Marshal(inputResModuleInfo)
		_, err = httpcli.ReqHttp(req, cModuleURL, common.HTTPCreate, []byte(resModuleJSONData))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoModuleCreateFailed)
		}
		inputFaultModuleInfo := make(map[string]interface{})
		inputFaultModuleInfo[common.BKSetIDField] = setID
		inputFaultModuleInfo[common.BKInstParentStr] = setID
		inputFaultModuleInfo[common.BKIdcIDField] = idcID
		inputFaultModuleInfo[common.BKModuleNameField] = common.DefaultFaultModuleName
		inputFaultModuleInfo[common.BKDefaultField] = common.DefaultFaultModuleFlag
		inputFaultModuleInfo[common.BKOwnerIDField] = ownerID
		resFaultModuleJSONData, _ := json.Marshal(inputFaultModuleInfo)
		_, err = httpcli.ReqHttp(req, cModuleURL, common.HTTPCreate, []byte(resFaultModuleJSONData))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcCreateFailed)
		}
		result := make(map[string]interface{})
		result[common.BKIdcIDField] = idcID

		return http.StatusOK, result, nil
	}, resp)
}

//get default idclication
func (cli *idcAction) GetDefaultIdc(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID, _ := pathParams["owner_id"]
		value, _ := ioutil.ReadAll(req.Request.Body)
		var js params.SearchParams
		err := json.Unmarshal([]byte(value), &js)
		if nil != err {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		condition := js.Condition
		condition[common.BKOwnerIDField] = ownerID
		condition[common.BKDefaultField] = common.DefaultAppFlag
		page := js.Page
		searchParams := make(map[string]interface{})
		searchParams["condition"] = condition
		searchParams["fields"] = strings.Join(js.Fields, ",")
		searchParams["start"] = page["start"]
		searchParams["limit"] = page["limit"]
		searchParams["sort"] = page["sort"]
		//search idc
		sIdcURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDIdc + "/search"
		inputJSON, _ := json.Marshal(searchParams)
		idcInfo, err := httpcli.ReqHttp(req, sIdcURL, common.HTTPSelectPost, []byte(inputJSON))
		blog.Info("get default idc params: %s", string(inputJSON))
		if nil != err {
			blog.Error("search idc error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcSearchFailed)
		}
		blog.Info("get default a idc return %v", idcInfo)
		idcJson, err := simplejson.NewJson([]byte(idcInfo))
		idcResData, _ := idcJson.Map()
		return http.StatusOK, idcResData["data"], nil
	}, resp)

}

//create default idclication
func (cli *idcAction) CreateDefaultIdc(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		pathParams := req.PathParameters()
		ownerID := pathParams["owner_id"]
		value, _ := ioutil.ReadAll(req.Request.Body)
		js, err := simplejson.NewJson([]byte(value))
		if nil != err {
			blog.Errorf("create default idc get params error %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		input, err := js.Map()
		valid := validator.NewValidMap(ownerID, common.BKInnerObjIDIdc, cli.CC.ObjCtrl(), forward, defErr)
		_, err = valid.ValidMap(input, common.ValidCreate, 0)
		if nil != err {
			blog.Errorf("create default idc get params error %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommFieldNotValid)
		}

		//create idclication
		input[common.BKOwnerIDField] = ownerID
		input[common.BKSupplierIDField] = common.BKDefaultSupplierID
		input[common.BKDefaultField] = common.DefaultAppFlag
		idcInfoJSON, _ := json.Marshal(input)
		cIdcURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDIdc
		cIdcRes, err := httpcli.ReqHttp(req, cIdcURL, common.HTTPCreate, []byte(idcInfoJSON))
		if nil != err {
			blog.Errorf("add default idclication error, ownerID:%s, error:%v ", ownerID, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcCreateFailed)
		}
		//create default set
		js, err = simplejson.NewJson([]byte(cIdcRes))
		idcResData, _ := js.Map()
		idcIDInfo := idcResData["data"].(map[string]interface{})
		idcID := idcIDInfo[common.BKIdcIDField]
		inputSetInfo := make(map[string]interface{})
		inputSetInfo[common.BKIdcIDField] = idcID
		inputSetInfo[common.BKInstParentStr] = idcID
		inputSetInfo[common.BKSetNameField] = common.DefaultResSetName
		inputSetInfo[common.BKDefaultField] = common.DefaultResSetFlag
		inputSetInfo[common.BKOwnerIDField] = ownerID
		cSetURL := cli.CC.ObjCtrl() + "/object/v1/insts/set"
		setJSONData, _ := json.Marshal(inputSetInfo)
		cSetRes, err := httpcli.ReqHttp(req, cSetURL, common.HTTPCreate, []byte(setJSONData))
		if nil != err {
			blog.Errorf("add default idclication Set error, ownerID:%s, error:%v ", ownerID, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcCreateFailed)
		}
		//create default module
		js, err = simplejson.NewJson([]byte(cSetRes))
		setResData, _ := js.Map()
		setIDInfo := setResData["data"].(map[string]interface{})
		setID := setIDInfo[common.BKSetIDField]
		inputResModuleInfo := make(map[string]interface{})
		inputResModuleInfo[common.BKSetIDField] = setID
		inputResModuleInfo[common.BKInstParentStr] = setID
		inputResModuleInfo[common.BKIdcIDField] = idcID
		inputResModuleInfo[common.BKModuleNameField] = common.DefaultResModuleName
		inputResModuleInfo[common.BKDefaultField] = common.DefaultResModuleFlag
		inputResModuleInfo[common.BKOwnerIDField] = ownerID
		cModuleURL := cli.CC.ObjCtrl() + "/object/v1/insts/module"
		resModuleJSONData, _ := json.Marshal(inputResModuleInfo)
		_, err = httpcli.ReqHttp(req, cModuleURL, common.HTTPCreate, []byte(resModuleJSONData))
		if nil != err {
			blog.Errorf("add default idclication module error, ownerID:%s, error:%v ", ownerID, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcCreateFailed)
		}
		inputFaultModuleInfo := make(map[string]interface{})
		inputFaultModuleInfo[common.BKSetIDField] = setID
		inputFaultModuleInfo[common.BKInstParentStr] = setID
		inputFaultModuleInfo[common.BKIdcIDField] = idcID
		inputFaultModuleInfo[common.BKModuleNameField] = common.DefaultFaultModuleName
		inputFaultModuleInfo[common.BKDefaultField] = common.DefaultFaultModuleFlag
		inputFaultModuleInfo[common.BKOwnerIDField] = ownerID
		resFaultModuleJSONData, _ := json.Marshal(inputFaultModuleInfo)
		_, err = httpcli.ReqHttp(req, cModuleURL, common.HTTPCreate, []byte(resFaultModuleJSONData))
		if nil != err {
			blog.Errorf("add default idclication module error, ownerID:%s, error info is %v ", ownerID, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcIdcCreateFailed)
		}
		result := make(map[string]interface{})
		result[common.BKIdcIDField] = idcID

		return http.StatusOK, result, nil
	}, resp)
}
