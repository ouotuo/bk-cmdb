/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package instdata

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"encoding/json"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/util"
	dcCommon "configcenter/src/scene_server/datacollection/common"
	eventtypes "configcenter/src/scene_server/event_server/types"
	"configcenter/src/source_controller/common/commondata"
	"configcenter/src/source_controller/common/eventdata"
	"configcenter/src/source_controller/common/instdata"

	"github.com/rs/xid"
	"io"
)

var host *hostAction = &hostAction{}
var SwitchTable = "cc_SwitchBase"

type hostAction struct {
	base.BaseAction
}

//AddHost add host to resource
func (cli *hostAction) AddHost(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		objType := common.BKInnerObjIDHost
		instdata.DataH = cli.CC.InstCli
		value, _ := ioutil.ReadAll(req.Request.Body)
		js, _ := simplejson.NewJson([]byte(value))
		input, _ := js.Map()
		blog.Info("create object type:%s,data:%v", objType, input)
		input[common.CreateTimeField] = time.Now()
		var idName string
		ID, err := instdata.CreateObject(objType, input, &idName)
		if err != nil {
			blog.Error("create object type:%s,data:%v error:%v", objType, input, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostCreateInst)
		}

		// record event
		originData := map[string]interface{}{}
		if err := instdata.GetObjectByID(objType, nil, ID, originData, ""); err != nil {
			blog.Error("create event error:%v", err)
		} else {
			ec := eventdata.NewEventContextByReq(req)
			err := ec.InsertEvent(eventtypes.EventTypeInstData, "host", eventtypes.EventActionCreate, originData, nil)
			if err != nil {
				blog.Error("create event error:%v", err)
			}
		}

		info := make(map[string]int)
		info[idName] = ID
		return http.StatusOK, info, nil
	}, resp)
}

//GetHostByID get host detail
func (cli *hostAction) GetHostByID(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		pathParams := req.PathParameters()
		hostID, _ := strconv.Atoi(pathParams["bk_host_id"])
		var result interface{}
		condition := make(map[string]interface{})
		condition[common.BKHostIDField] = hostID
		fields := make([]string, 0)
		err := cli.CC.InstCli.GetOneByCondition("cc_HostBase", fields, condition, &result)
		if err != nil {
			blog.Error("get GetHostByID err %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDBSelectFailed)
		}
		return http.StatusOK, result, nil
	}, resp)
}

//GetHosts batch search host
func (cli *hostAction) GetHosts(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		objType := common.BKInnerObjIDHost
		instdata.DataH = cli.CC.InstCli

		value, err := ioutil.ReadAll(req.Request.Body)
		var dat commondata.ObjQueryInput
		err = json.Unmarshal([]byte(value), &dat)
		if err != nil {
			blog.Error("get object type:%s,input:%v error:%v", objType, value, err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		fields := dat.Fields
		condition := util.ConvParamsTime(dat.Condition)
		start := dat.Start
		limit := dat.Limit
		sort := dat.Sort
		fieldArr := strings.Split(fields, ",")
		result := make([]map[string]interface{}, 0)
		count, err := instdata.GetCntByCondition(objType, condition)
		if err != nil {
			blog.Error("get object type:%s,input:%s error:%v", objType, value, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostSelectInst)
		}
		err = instdata.GetObjectByCondition(defLang, objType, fieldArr, condition, &result, sort, start, limit)
		if err != nil {
			blog.Error("get object type:%s,input:%v error:%v", objType, value, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostSelectInst)
		}
		info := make(map[string]interface{})
		info["count"] = count
		info["info"] = result
		return http.StatusOK, info, nil
	}, resp)
}

//GetHostSnap get host snap
func (cli *hostAction) GetHostSnap(req *restful.Request, resp *restful.Response) {
	redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		hostID := req.PathParameter("bk_host_id")

		var result string
		err := redisCli.Get(dcCommon.RedisSnapKeyPrefix + hostID).Scan(&result)
		if err != nil {
			statuscode := 0
			err := redisCli.Get(dcCommon.RedisSnapKeyChannelStatus).Scan(&statuscode)
			if err != nil {
				blog.Error("get host snapshot error,input:%v error:%v", hostID, err)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetSnapshot)
			}

			if statuscode != common.CCSuccess {
				return http.StatusInternalServerError, nil, defErr.Error(statuscode)
			}
			blog.Error("get host snapshot error,input:%v error:%v", hostID, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetSnapshot)
		}

		return http.StatusOK, common.KvMap{"data": result}, nil
	}, resp)
}

//Add add new SwitchAdd
func (u *hostAction) SwitchAdd(req *restful.Request, resp *restful.Response) {

	value, _ := ioutil.ReadAll(req.Request.Body)

	language := util.GetActionLanguage(req)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	params := make(map[string]interface{})
	if err := json.Unmarshal([]byte(value), &params); nil != err {
		blog.Error("fail to unmarshal json, error information is %s, msg:%s", err.Error(), string(value))
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}

	//mogo generate id for product
	xidDevice := xid.New()
	params[common.CreateTimeField] = time.Now()
	params[common.LastTimeField] = ""
	_, err := u.CC.InstCli.Insert(SwitchTable, params)

	if err != nil {
		blog.Error("create switch api  error:data:%v error:%v", params, err)
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommDBInsertFailed, defErr.Error(common.CCErrCommDBInsertFailed).Error(), resp)
		return
	}

	info := make(map[string]interface{})
	info["id"] = xidDevice.String()

	rsp, _ := u.CC.CreateAPIRspStr(common.CCSuccess, info)
	io.WriteString(resp, rsp)
}

//Update update user api content
func (u *hostAction) SwitchUpdate(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	value, _ := ioutil.ReadAll(req.Request.Body)
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(value), &data); nil != err {
		blog.Error("fail to unmarshal json, error information is %s, msg:%s", err.Error(), string(value))
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}

	sw, _ := data["data"]
	swData := sw.(map[string]interface{})

	condition, _ := data["condition"]
	swData[common.LastTimeField] = time.Now()
	condition = condition.(map[string]interface{})
	params := make(map[string]interface{})
	params[common.BKBindIpField] = swData[common.BKBindIpField]

	rowCount, err := u.CC.InstCli.GetCntByCondition(SwitchTable, params)
	if nil != err {
		blog.Error("query user api fail, error information is %s, params:%v", err.Error(), params)
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommDBSelectFailed, defErr.Error(common.CCErrCommDBSelectFailed).Error(), resp)
		return
	}
	if 1 != rowCount {
		blog.Info("host user api not permissions or not exists, params:%v", params)
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommNotFound, defErr.Error(common.CCErrCommNotFound).Error(), resp)
		return
	}
	err = u.CC.InstCli.UpdateByCondition(SwitchTable, swData, condition)
	if nil != err {
		blog.Error("updata user api fail, error information is %s, params:%v", err.Error(), params)
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommDBUpdateFailed, defErr.Errorf(common.CCErrCommDBUpdateFailed).Error(), resp)
		return
	}
	rsp, _ := u.CC.CreateAPIRspStr(common.CCSuccess, nil)
	io.WriteString(resp, rsp)
	return

}

//GetHosts batch search host
func (cli *hostAction) GetSwitchs(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		objType := common.BKInnerObjIDSwitchHost
		instdata.DataH = cli.CC.InstCli

		value, err := ioutil.ReadAll(req.Request.Body)
		var dat commondata.ObjQueryInput
		err = json.Unmarshal([]byte(value), &dat)
		if err != nil {
			blog.Error("cccget object type:%s,input:%v error:%v", objType, value, err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		blog.Info("test for get switch is ", dat)
		fields := dat.Fields
		condition := util.ConvParamsTime(dat.Condition)
		start := dat.Start
		limit := dat.Limit
		sort := dat.Sort
		fieldArr := strings.Split(fields, ",")
		result := make([]map[string]interface{}, 0)
		count, err := instdata.GetCntByCondition(objType, condition)
		if err != nil {
			blog.Error("aaget object type:%s,input:%s error:%v", objType, value, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostSelectInst)
		}
		err = instdata.GetObjectByCondition(defLang, objType, fieldArr, condition, &result, sort, start, limit)
		if err != nil {
			blog.Error("bbget object type:%s,input:%v error:%v", objType, value, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostSelectInst)
		}
		info := make(map[string]interface{})
		info["count"] = count
		info["info"] = result
		return http.StatusOK, info, nil
	}, resp)
}

//GetHosts batch search host
func (cli *hostAction) GetSwitchsPort(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Errorf("read http request boody error:%s", err.Error())
			cli.ResponseFailed(common.CCErrCommHTTPReadBodyFailed, defErr.Error(common.CCErrCommHTTPReadBodyFailed).Error(), resp)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostSelectInst)

		}
		var dat commondata.ObjQueryInput
		err = json.Unmarshal([]byte(value), &dat)
		if err != nil {
			blog.Error("json unmarshal failed,input:%v error:%v", string(value), err)
			cli.ResponseFailed(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostSelectInst)

		}
		fields := dat.Fields
		condition := dat.Condition
		start := dat.Start
		limit := dat.Limit
		sort := dat.Sort
		fieldArr := strings.Split(fields, ",")
		result := make([]map[string]interface{}, 0)
		instdata.DataH = cli.CC.InstCli
		err = instdata.GetObjectByCondition(defLang, common.BKInnerObjIDSwitchHost, fieldArr, condition, &result, sort, start, limit)
		if nil != err {
			blog.Error("get data from data  error:%s", err.Error())
			cli.ResponseFailed(common.CCErrCommDBSelectFailed, defErr.Error(common.CCErrCommDBSelectFailed).Error(), resp)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostSelectInst)

		}
		count, err := instdata.GetCntByCondition(common.BKInnerObjIDSwitchHost, dat.Condition)
		if nil != err {
			blog.Error("get data from data  error:%s", err.Error())
			cli.ResponseFailed(common.CCErrCommDBSelectFailed, defErr.Error(common.CCErrCommDBSelectFailed).Error(), resp)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostSelectInst)

		}
		info := make(map[string]interface{})
		info["count"] = count
		info["info"] = result
		return http.StatusOK, info, nil
	}, resp)
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/host/{bk_host_id}", Params: nil, Handler: host.GetHostByID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/search", Params: nil, Handler: host.GetHosts})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/insts", Params: nil, Handler: host.AddHost})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/host/snapshot/{bk_host_id}", Params: nil, Handler: host.GetHostSnap})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/switch/add", Params: nil, Handler: host.SwitchAdd})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/switch/update", Params: nil, Handler: host.SwitchUpdate})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/switch/get", Params: nil, Handler: host.GetSwitchs})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/switch/port", Params: nil, Handler: host.GetSwitchsPort})

	// create cc object
	host.CreateAction()
}
