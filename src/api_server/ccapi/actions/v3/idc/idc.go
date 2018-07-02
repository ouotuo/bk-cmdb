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

package idc

import (
	"configcenter/src/api_server/ccapi/actions/v3"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/scene_server/api"

	"github.com/emicklei/go-restful"

	"fmt"
)

var idc = &idcAction{}

type idcAction struct {
	base.BaseAction
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/idc1/test/api", Params: nil, Handler: idc.TestApi, Version: v3.APIVersion})

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/idc1/model/{owner_id}", Params: nil, Handler: idc.SelectTopoModel, Version: v3.APIVersion})
	// set cc api interface
	idc.CreateAction()
}
func (cli *idcAction) TestApi(req *restful.Request, resp *restful.Response) {
	blog.Info("TestApi TestApi TestApiTestApi ")
	senceCLI := api.NewClient(idc.CC.IdcAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaIdc(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, "0"),
		resp)
}

// SelectTopoInst search the inst idc tree
func (cli *idcAction) SelectIdcoInst(req *restful.Request, resp *restful.Response) {

	blog.Info("select idc inst ")
	ownerID := req.PathParameter("owner_id")
	appID := req.PathParameter("app_id")
	senceCLI := api.NewClient(idc.CC.IdcAPI())

	cli.CallResponse(
		senceCLI.ReForwardSelectMetaIdcInst(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, fmt.Sprintf("%s?level=%s", url, req.QueryParameter("level")), method)
		}, ownerID, appID),
		resp)
}
// SelectTopoModel search the main line object idc tree
func (cli *idcAction) SelectTopoModel(req *restful.Request, resp *restful.Response) {

	blog.Info("select idc model")

	ownerID := req.PathParameter("owner_id")

	senceCLI := api.NewClient(idc.CC.IdcAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaIdc(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID),
		resp)

}
// SelectSet search some sets
func (cli *idcAction) SelectIdc(req *restful.Request, resp *restful.Response) {

	blog.Info("select idc")

	ownerID := req.PathParameter("owner_id")
	appID := req.PathParameter("app_id")

	senceCLI := api.NewClient(idc.CC.IdcAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaIdc1(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, appID),
		resp)

}
