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

package idc

import (
	"configcenter/src/api_server/ccapi/actions/v3"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"io"
	"configcenter/src/scene_server/api"
	"github.com/emicklei/go-restful"
	"configcenter/src/common/http/httpclient"
)

var app = &appAction{}

type appAction struct {
	base.BaseAction
}

// CreateApp create idclication
func (cli *appAction) CreateApp(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	url := cli.CC.IdcAPI() + "/idc/v1/idc/" + ownerID
	//	req.Request.URL.Path = "/idc/v1/idc/" + ownerID
	blog.Info("Create App url:%s", req.Request.URL.Path)
	//	httpcli.ProxyRestHttp(req, resp, url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPCreate)
	io.WriteString(resp, rsp)
}

// DeleteApp delete idclication
func (cli *appAction) DeleteApp(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	idcID := pathParams["idc_id"]
	url := cli.CC.IdcAPI() + "/idc/v1/idc/" + ownerID + "/" + idcID
	//req.Request.URL.Path = "/idc/v1/idc/" + ownerID + "/" + idcID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPDelete)
	io.WriteString(resp, rsp)
}

// UpdateApp update idclication
func (cli *appAction) UpdateApp(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	idcID := pathParams["idc_id"]
	url := cli.CC.IdcAPI() + "/idc/v1/idc/" + ownerID + "/" + idcID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPUpdate)
	io.WriteString(resp, rsp)
}

// UpdateAppDataStatus update idclication data status
func (cli *appAction) UpdateAppDataStatus(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	idcID := pathParams["idc_id"]
	flag := pathParams["flag"]
	url := cli.CC.IdcAPI() + "/idc/v1/idc/status/" + flag + "/" + ownerID + "/" + idcID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPUpdate)
	io.WriteString(resp, rsp)
}

// SearchApp search idclication
func (cli *appAction) SearchApp(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	url := cli.CC.IdcAPI() + "/idc/v1/idc/search/" + ownerID
	blog.Info("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&")
	blog.Info(url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

// GetInternalTopo get internal idc
func (cli *appAction) GetInternalTopo(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	idcID := pathParams["idc_id"]
	url := cli.CC.IdcAPI() + "/idc/v1/idc/internal/" + ownerID + "/" + idcID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectGet)
	io.WriteString(resp, rsp)
}
func (cli *appAction) TestApi(req *restful.Request, resp *restful.Response) {
	blog.Info("TestApi TestApi TestApiTestApi ")
	senceCLI := api.NewClient(idc.CC.IdcAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaIdc(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, "0"),
		resp)
}
func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/idc/test/app", Params: nil, Handler: app.TestApi, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/idc/{owner_id}", Params: nil, Handler: app.CreateApp, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/idc/{owner_id}/{idc_id}", Params: nil, Handler: app.DeleteApp, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/idc/{owner_id}/{idc_id}", Params: nil, Handler: app.UpdateApp, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/idc/search/{owner_id}", Params: nil, Handler: app.SearchApp, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/idc/status/{flag}/{owner_id}/{idc_id}", Params: nil, Handler: app.UpdateAppDataStatus, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/idc/internal/{owner_id}/{idc_id}", Params: nil, Handler: app.GetInternalTopo, Version: v3.APIVersion})
	// set cc api interface
	app.CreateAction()
}
