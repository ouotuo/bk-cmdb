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
)

var app = &apptestAction{}

type apptestAction struct {
	base.BaseAction
}

func init() {
	blog.Info("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/testapi/test", Params: nil, Handler: app.ApiTest})
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
