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
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"encoding/json"
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful"
)

var topo = &topoAction{}

// topoAction
type topoAction struct {
	base.BaseAction
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/idc/internal/{owner_id}/{idc_id}", Params: nil, Handler: topo.GetInternalPos})
	topo.CreateAction()
}

//get built in pos
func (cli *topoAction) GetInternalPos(req *restful.Request, resp *restful.Response) {

	// get language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		idcIDStr := req.PathParameter("idc_id")
		idcID, _ := strconv.Atoi(idcIDStr)
		rackCond := make(map[string]interface{})
		cond := make(map[string]interface{})
		cond[common.BKIdcIDField] = idcID
		//cond[common.BKDefaultField] = common.DefaultResModuleFlag
		rackCond["condition"] = cond

		//search rack
		sURL := cli.CC.ObjCtrl() + "/object/v1/insts/rack/search"
		inputJSON, _ := json.Marshal(rackCond)
		rackRes, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, []byte(inputJSON))
		blog.Info("search rack params: %s", string(inputJSON))
		blog.Info("search rack return: %s", string(rackRes))
		if nil != err {
			blog.Error("search rack error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcRackSelectFailed)
		}

		posCond := make(map[string]interface{})
		//defaultCond := make(map[string]interface{})
	//	defaultCond[common.BKDBIN] = []int{common.DefaultResModuleFlag, common.DefaultFaultModuleFlag}
		//cond[common.BKDefaultField] = defaultCond
		posCond["condition"] = cond

		//search pos
		sURL = cli.CC.ObjCtrl() + "/object/v1/insts/pos/search"
		inputJSON, _ = json.Marshal(posCond)
		posRes, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, []byte(inputJSON))
		blog.Debug("search pos params: %s", string(inputJSON))
		blog.Info("search pos return: %s", string(posRes))
		if nil != err {
			blog.Error("search pos error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcPosSelectFailed)
		}
		var rackObj api.APIRsp
		var posObj api.APIRsp
		err = json.Unmarshal([]byte(rackRes), &rackObj)
		if nil != err || !rackObj.Result {
			blog.Error("search rack error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcRackSelectFailed)
		}
		err = json.Unmarshal([]byte(posRes), &posObj)
		if nil != err || !posObj.Result {
			blog.Error("search rack error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcPosSelectFailed)
		}
		rackData := rackObj.Data.(map[string]interface{})
		rackInfo := rackData["info"].([]interface{})
		posData := posObj.Data.(map[string]interface{})
		posInfo := posData["info"].([]interface{})
		if 0 == len(rackInfo) || 0 == len(posInfo) {
			blog.Error("search rack error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrIdcPosSelectFailed)
		}

		rackResult := make(map[string]interface{})
		posResult := make([]map[string]interface{}, 0)
		for _, i := range rackInfo {
			rack := i.(map[string]interface{})
			rackResult[common.BKRackIDField] = rack[common.BKRackIDField]
			rackResult[common.BKRackNameField] = rack[common.BKRackNameField]
		}
		for _, j := range posInfo {
			posR := make(map[string]interface{})
			pos := j.(map[string]interface{})
			posR[common.BKPosIDField] = pos[common.BKPosIDField]
			posR[common.BKPosNameField] = pos[common.BKPosNameField]
			posResult = append(posResult, posR)
		}
		rackResult[common.BKInnerObjIDPos] = posResult

		return http.StatusOK, rackResult, nil
	}, resp)

}
