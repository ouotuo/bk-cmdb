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

var rack = &rackAction{}

type rackAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/rack/{idc_id}", Params: nil, Handler: rack.CreateRack, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/rack/{idc_id}/{rack_id}", Params: nil, Handler: rack.DeleteRack, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/rack/{idc_id}/{rack_id}", Params: nil, Handler: rack.UpdateRack, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/rack/search/{owner_id}/{idc_id}", Params: nil, Handler: rack.SelectRack, Version: v3.APIVersion})

	// init
	rack.CreateAction()
}

// CreateRack create a rack
func (cli *rackAction) CreateRack(req *restful.Request, resp *restful.Response) {

	blog.Info("create rack")
	idcID := req.PathParameter("idc_id")

	senceCLI := api.NewClient(pos.CC.IdcAPI())
	cli.CallResponse(
		senceCLI.ReForwardCreateMetaRack(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, idcID),
		resp)

}

// DeleteRack delete a rack
func (cli *rackAction) DeleteRack(req *restful.Request, resp *restful.Response) {

	blog.Info("delete object")

	idcID := req.PathParameter("idc_id")
	rackID := req.PathParameter("rack_id")

	senceCLI := api.NewClient(pos.CC.IdcAPI())
	cli.CallResponse(
		senceCLI.ReForwardDeleteMetaRack(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, idcID, rackID),
		resp)
}

// UpdateRack update a rack
func (cli *rackAction) UpdateRack(req *restful.Request, resp *restful.Response) {

	blog.Info("update rack")

	idcID := req.PathParameter("idc_id")
	rackID := req.PathParameter("rack_id")

	senceCLI := api.NewClient(pos.CC.IdcAPI())
	cli.CallResponse(
		senceCLI.ReForwardUpdateMetaRack(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, idcID, rackID),
		resp)

}

// SelectRack search some racks
func (cli *rackAction) SelectRack(req *restful.Request, resp *restful.Response) {

	blog.Info("select rack")

	ownerID := req.PathParameter("owner_id")
	idcID := req.PathParameter("idc_id")

	senceCLI := api.NewClient(pos.CC.IdcAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaRack(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, idcID),
		resp)

}
