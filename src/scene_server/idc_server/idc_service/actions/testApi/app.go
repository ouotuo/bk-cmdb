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

package testApi

import (
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"net/http"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"configcenter/src/common/util"
	//api "configcenter/src/source_controller/api/object"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	httpcli "configcenter/src/common/http/httpclient"

)

var app = &apptestAction{}

type apptestAction struct {
	base.BaseAction
}

func init() {
	blog.Info("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/testapi/test", Params: nil, Handler: app.ApiTest})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/idc/{idc_id}", Params: nil, Handler: app.CreateIdc})
	// create CC object
	app.CreateAction()
}

//delete application
func (cli *apptestAction) ApiTest(req *restful.Request, resp *restful.Response) {
	cli.CallResponseEx(func() (int, interface{}, error) {

		//new feature, app not allow deletion
		blog.Info("app not allow deletion")
		return http.StatusOK, "11111111111111111", nil
	}, resp)
}
func (cli *apptestAction) CreateIdc(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	//user := util.GetActionUser(req)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {
		//forward := &api.ForwardParam{Header: req.Request.Header}
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
		//appID, convErr := strconv.Atoi(req.PathParameter("app_id"))
		//if nil != convErr {
		//	blog.Error("failed to convert the appid to int, error info is %s", convErr.Error())
		//	return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "app_id")
		//}
		// create
		input[common.BKDefaultField] = 0
		input[common.CreateTimeField] = util.GetCurrentTimeStr()

		inputJSON, jsErr := json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		cModuleURL := cli.CC.ObjCtrl() + "/object/v1/insts/idc"
		moduleRes, err := httpcli.ReqHttp(req, cModuleURL, common.HTTPCreate, inputJSON)
		if nil != err {
			blog.Error("failed to create the set, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoSetCreateFailed)
		}
		blog.Info("moduleRes:",moduleRes)



		return http.StatusOK, moduleRes, nil
	}, resp)

}