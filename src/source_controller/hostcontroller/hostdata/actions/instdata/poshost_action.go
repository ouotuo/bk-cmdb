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

package instdata

import (
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/common/eventdata"
	"configcenter/src/source_controller/hostcontroller/hostdata/logics"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"
)

type PosHostConfigParams struct {
	IdclicationID int   `json:"bk_idc_id"`
	HostID        int   `json:"bk_host_id"`
	PosID      []int `json:"bk_pos_id"`
}

var (
	posBaseTaleName = "cc_PosBase"
)

var posHostConfigActionCli *posHostConfigAction = &posHostConfigAction{}

// HostAction
type posHostConfigAction struct {
	base.BaseAction
}

//AddPosHostConfig add pos host config
func (cli *posHostConfigAction) AddPosHostConfig(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		cc := api.NewAPIResource()
		//instdata.DataH = cc.InstCli
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		params := PosHostConfigParams{}
		if err := json.Unmarshal([]byte(value), &params); nil != err {
			blog.Error("fail to unmarshal json, error information is %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		hostID := params.HostID

		//add new relation ship
		ec := eventdata.NewEventContextByReq(req)
		for _, posID := range params.PosID {
			_, err := logics.AddSingleHostPosRelation(ec, cc, hostID, posID, params.IdclicationID)
			if nil != err {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostTransferIdc)
			}
		}

		return http.StatusOK, nil, nil
	}, resp)
}

//DelDefaultPosHostConfig delete default pos host config
func (cli *posHostConfigAction) DelDefaultPosHostConfig(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		cc := api.NewAPIResource()
		//instdata.DataH = cc.InstCli
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		params := PosHostConfigParams{}
		if err = json.Unmarshal([]byte(value), &params); nil != err {
			blog.Error("fail to unmarshal json, error information is %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		defaultPosIDs, err := logics.GetDefaultPosIDs(cc, params.IdclicationID)
		if nil != err {
			blog.Errorf("defaultPosIds idcID:%d, error:%v", params.IdclicationID, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrGetPos)
		}

		hostID := params.HostID

		//delete default host pos relation
		ec := eventdata.NewEventContextByReq(req)
		for _, defaultPosID := range defaultPosIDs {
			_, err := logics.DelSingleHostPosRelation(ec, cc, hostID, defaultPosID, params.IdclicationID)
			if nil != err {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrDelDefaultPosHostConfig)
			}
		}

		return http.StatusOK, nil, nil
	}, resp)
}

//DelPosHostConfig delete pos host config
func (cli *posHostConfigAction) DelPosHostConfig(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		cc := api.NewAPIResource()
		//instdata.DataH = cc.InstCli
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		params := PosHostConfigParams{}
		if err = json.Unmarshal([]byte(value), &params); nil != err {
			blog.Error("fail to unmarshal json, error information is %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		getPosParams := make(map[string]interface{}, 2)
		getPosParams[common.BKHostIDField] = params.HostID
		getPosParams[common.BKIdcIDField] = params.IdclicationID
		posIDs, err := logics.GetPosIDsByHostID(cc, getPosParams) //params.HostID, params.IdclicationID)
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrGetOriginHostModuelRelationship)
		}

		ec := eventdata.NewEventContextByReq(req)
		for _, posID := range posIDs {
			_, err := logics.DelSingleHostPosRelation(ec, cc, params.HostID, posID, params.IdclicationID)
			if nil != err {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrDelOriginHostModuelRelationship)
			}
		}

		return http.StatusOK, nil, nil
	}, resp)
}

//GetHostPossIDs get host pos ids
func (cli *posHostConfigAction) GetHostPossIDs(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		cc := api.NewAPIResource()
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		params := PosHostConfigParams{}
		if err = json.Unmarshal([]byte(value), &params); nil != err {
			blog.Error("fail to unmarshal json, error information is %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		posIDs, err := logics.GetPosIDsByHostID(cc, map[string]interface{}{common.BKIdcIDField: params.IdclicationID, common.BKHostIDField: params.HostID}) //params.HostID, params.IdclicationID)
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrGetPos)
		}
		return http.StatusOK, posIDs, nil
	}, resp)
}

//AssignHostToIdc assign host to idc
func (cli *posHostConfigAction) AssignHostToIdc(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		type paramsStruct struct {
			IdclicationID      int   `json:"bk_idc_id"`
			HostID             []int `json:"bk_host_id"`
			PosID           int   `json:"bk_pos_id"`
			OwnerIdclicationID int   `json:"bk_owner_idc_id"`
			OwnerPosID      int   `json:"bk_owner_pos_id"`
		}

		cc := api.NewAPIResource()
		ec := eventdata.NewEventContextByReq(req)
		//instdata.DataH = cc.InstCli
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		params := paramsStruct{}
		if err := json.Unmarshal([]byte(value), &params); nil != err {
			blog.Error("fail to unmarshal json, error information is %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		getPosParams := make(map[string]interface{})
		for _, hostID := range params.HostID {
			//delete relation in default idc pos
			_, err := logics.DelSingleHostPosRelation(ec, cc, hostID, params.OwnerPosID, params.OwnerIdclicationID)
			if nil != err {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTransferHostFromPool)
			}
			getPosParams[common.BKHostIDField] = hostID
			posIDs, err := logics.GetPosIDsByHostID(cc, getPosParams)
			if nil != err {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrGetPos)
			}
			//delete from empty pos, no relation
			if 0 < len(posIDs) {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAlreadyAssign)
			}

			//add new host
			_, err = logics.AddSingleHostPosRelation(ec, cc, hostID, params.PosID, params.IdclicationID)
			if nil != err {
			}
		}
		return http.StatusOK, nil, nil
	}, resp)
	return

}

//GetPossHostConfig  get pos host config
func (cli *posHostConfigAction) GetPossHostConfig(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		var params = make(map[string][]int)
		cc := api.NewAPIResource()
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		if err = json.Unmarshal([]byte(value), &params); nil != err {
			blog.Error("fail to unmarshal json, error information is %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		query := make(map[string]interface{})
		for key, val := range params {
			conditon := make(map[string]interface{})
			conditon[common.BKDBIN] = val
			query[key] = conditon
		}
		fields := []string{common.BKIdcIDField, common.BKHostIDField, common.BKSetIDField, common.BKPosIDField}
		var result []interface{}
		err = cc.InstCli.GetMutilByCondition("cc_IdcHostConfig", fields, query, &result, common.BKHostIDField, 0, 100000)
		if err != nil {
			blog.Error("fail to get pos host config %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDBSelectFailed)
		}
		return http.StatusOK, result, nil
	}, resp)
}



func init() {
	posHostConfigActionCli.CreateAction()
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/meta/hosts/poss/search", Params: nil, Handler: posHostConfigActionCli.GetHostPossIDs})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/meta/hosts/poss", Params: nil, Handler: posHostConfigActionCli.AddPosHostConfig})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/meta/hosts/poss", Params: nil, Handler: posHostConfigActionCli.DelPosHostConfig})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/meta/hosts/defaultposs", Params: nil, Handler: posHostConfigActionCli.DelDefaultPosHostConfig})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/meta/hosts/assign", Params: nil, Handler: posHostConfigActionCli.AssignHostToIdc})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/meta/hosts/pos/config/search", Params: nil, Handler: posHostConfigActionCli.GetPossHostConfig})
}
