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
	"configcenter/src/common/http/httpclient"
	"configcenter/src/scene_server/api"

	"github.com/emicklei/go-restful"
)

var pos = &posAction{}

type posAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/pos/{idc_id}/{rack_id}", Params: nil, Handler: pos.CreatePos, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/pos/{idc_id}/{rack_id}/{pos_id}", Params: nil, Handler: pos.DeletePos, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/pos/{idc_id}/{rack_id}/{pos_id}", Params: nil, Handler: pos.UpdatePos, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/pos/search/{owner_id}/{idc_id}/{rack_id}", Params: nil, Handler: pos.SelectPos, Version: v3.APIVersion})

	// rack cc api interface
	pos.CreateAction()

}

// CreatePos create a pos of the rack
func (cli *posAction) CreatePos(req *restful.Request, resp *restful.Response) {

	blog.Info("create object")

	idcID := req.PathParameter("idc_id")
	rackID := req.PathParameter("rack_id")

	senceCLI := api.NewClient(pos.CC.IdcAPI())

	cli.CallResponse(
		senceCLI.ReForwardCreateMetaPos(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, idcID, rackID),
		resp)

}

// DeletePos delete the pos
func (cli *posAction) DeletePos(req *restful.Request, resp *restful.Response) {

	blog.Info("delete pos")

	idcID := req.PathParameter("idc_id")
	rackID := req.PathParameter("rack_id")
	posID := req.PathParameter("pos_id")

	senceCLI := api.NewClient(pos.CC.IdcAPI())
	cli.CallResponse(
		senceCLI.ReForwardDeleteMetaPos(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, idcID, rackID, posID),
		resp)

}

// UpdatePos update the pos information
func (cli *posAction) UpdatePos(req *restful.Request, resp *restful.Response) {

	blog.Info("update pos")

	idcID := req.PathParameter("idc_id")
	rackID := req.PathParameter("rack_id")
	posID := req.PathParameter("pos_id")

	senceCLI := api.NewClient(pos.CC.IdcAPI())

	cli.CallResponse(
		senceCLI.ReForwardUpdateMetaPos(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, idcID, rackID, posID),
		resp)

}

// SelectPos search the pos detail information
func (cli *posAction) SelectPos(req *restful.Request, resp *restful.Response) {

	blog.Info("select pos ")
	ownerID := req.PathParameter("owner_id")
	idcID := req.PathParameter("idc_id")
	rackID := req.PathParameter("rack_id")

	senceCLI := api.NewClient(pos.CC.IdcAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaPos(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, idcID, rackID),
		resp)
}
