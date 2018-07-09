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

package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/api/auditlog"
	"configcenter/src/source_controller/api/metadata"
	"configcenter/src/source_controller/common/commondata"
	"encoding/json"
	"errors"
	"fmt"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)



type HostPosConfigLog struct {
	cur       []interface{}
	pre       []interface{}
	req       *restful.Request
	ownerID   string
	hostCtrl  string
	objCtrl   string
	auditCtrl string
	hostInfos []interface{}
	instID    []int
	desc      string
}






func (h *HostPosConfigLog) getHostPosConfig() []interface{} {

	conds := common.KvMap{common.BKHostIDField: h.instID}
	inputJson, _ := json.Marshal(conds)
	gHostURL := h.hostCtrl + "/host/v1/meta/hosts/pos/config/search"

	gHostRe, err := httpcli.ReqHttp(h.req, gHostURL, common.HTTPSelectPost, inputJson)
	blog.Infof("GetHostPosConfig, input:%s, return:%s", string(inputJson), gHostRe)
	if nil != err {
		blog.Error("getHostPosConfig info error :%v, url:%s", err, gHostURL)
		return nil
	}
	//
	js, err := simplejson.NewJson([]byte(gHostRe))

	gResult, _ := js.Get("result").Bool()
	if false == gResult {
		blog.Error("getHostPosConfig  info error :%v", err)
		return nil
	}

	//
	hostData, _ := js.Get("data").Array()
	return hostData
}

func (h *HostPosConfigLog) getInnerIP() []interface{} {

	var dat commondata.ObjQueryInput

	dat.Fields = fmt.Sprintf("%s,%s", common.BKHostIDField, common.BKHostInnerIPField)
	dat.Condition = common.KvMap{common.BKHostIDField: common.KvMap{common.BKDBIN: h.instID}}
	dat.Start = 0
	dat.Limit = common.BKNoLimit
	inputJson, _ := json.Marshal(dat)
	gHostURL := h.hostCtrl + "/host/v1/hosts/search"

	gHostRe, err := httpcli.ReqHttp(h.req, gHostURL, common.HTTPSelectPost, inputJson)
	blog.Infof("getInnerIP, input:%s, replay:%s", string(inputJson), gHostRe)
	if nil != err {
		blog.Error("GetInnerIP info error :%v, url:%s", err, gHostURL)
		return nil
	}
	//
	js, err := simplejson.NewJson([]byte(gHostRe))

	gResult, _ := js.Get("result").Bool()
	if false == gResult {
		blog.Error("GetHostDetail  info error :%v", err)
		return nil
	}

	//
	hostData, _ := js.Get("data").Get("info").Array()
	return hostData
}

func (h *HostPosConfigLog) getPoss(posIds []int) ([]interface{}, error) {
	if 0 == len(posIds) {
		return nil, nil
	}

	var dat commondata.ObjQueryInput

	dat.Fields = fmt.Sprintf("%s,%s,%s,%s,%s", common.BKPosIDField, common.BKSetIDField, common.BKPosNameField, common.BKAppIDField, common.BKOwnerIDField)
	dat.Limit = common.BKNoLimit
	dat.Start = 0
	dat.Condition = common.KvMap{common.BKPosIDField: common.KvMap{common.BKDBIN: posIds}}
	bodyContent, _ := json.Marshal(dat)
	url := h.objCtrl + "/object/v1/insts/pos/search"
	blog.Info("getPoss url :%s", url)
	blog.Info("getPoss content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(h.req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("getPoss return :%s", string(reply))
	if err != nil {
		blog.Errorf("getPoss url:%s, input:%s error:%s, ", url, string(bodyContent), err.Error())
		return nil, err
	}
	js, err := simplejson.NewJson([]byte(reply))
	if nil != err {
		blog.Errorf("getPoss url:%s, input:%s error:%s, ", url, string(bodyContent), err.Error())

		return nil, err
	}
	return js.Get("data").Get("info").Array()
}

func (h *HostPosConfigLog) getSets(setIds []int) ([]interface{}, error) {
	var dat commondata.ObjQueryInput

	dat.Fields = fmt.Sprintf("%s,%s,%s", common.BKSetNameField, common.BKSetIDField, common.BKOwnerIDField)
	dat.Limit = common.BKNoLimit
	dat.Start = 0
	dat.Condition = common.KvMap{common.BKSetIDField: common.KvMap{common.BKDBIN: setIds}}
	bodyContent, _ := json.Marshal(dat)
	url := h.objCtrl + "/object/v1/insts/set/search"
	blog.Info("getSets url :%s", url)
	blog.Info("getSets content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(h.req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("getSets return :%s", string(reply))
	if err != nil {
		blog.Errorf("getSets url:%s, input:%s error:%s, ", url, string(bodyContent), err.Error())
		return nil, err
	}
	js, err := simplejson.NewJson([]byte(reply))
	if nil != err {
		blog.Errorf("getSets url:%s, input:%s error:%s, ", url, string(bodyContent), err.Error())

		return nil, err
	}
	return js.Get("data").Get("info").Array()

}

//set host id, host id must be nil
func (h *HostPosConfigLog) SetHostID(hostID []int) error {
	if nil == h.instID {
		h.instID = hostID
		h.hostInfos = h.getInnerIP()
		return nil
	}
	return errors.New("hostID not empty")
}

func (h *HostPosConfigLog) SetDesc(desc string) {
	h.desc = desc
}

func (h *HostPosConfigLog) SaveLog(appID, user string) error {
	//gHostURL := "http://" + cli.CC.HostCtrl + "/host/v1/host/" + hostID
	h.cur = h.getHostPosConfig()

	var setIDs []int
	var posIDs []int
	preMap := make(map[int]map[int]interface{})
	curMap := make(map[int]map[int]interface{})

	for _, val := range h.pre {
		valMap, _ := val.(map[string]interface{})
		hostID, _ := util.GetIntByInterface(valMap[common.BKHostIDField])
		mID, _ := util.GetIntByInterface(valMap[common.BKPosIDField])
		sID, _ := util.GetIntByInterface(valMap[common.BKSetIDField])
		if _, ok := preMap[hostID]; false == ok {
			preMap[hostID] = make(map[int]interface{}, 0)
		}
		preMap[hostID][mID] = valMap
		setIDs = append(setIDs, sID)
		posIDs = append(posIDs, mID)
	}
	for _, val := range h.cur {
		valMap, _ := val.(map[string]interface{})
		hostID, _ := util.GetIntByInterface(valMap[common.BKHostIDField])
		mID, _ := util.GetIntByInterface(valMap[common.BKPosIDField])
		sID, _ := util.GetIntByInterface(valMap[common.BKSetIDField])
		if _, ok := curMap[hostID]; false == ok {
			curMap[hostID] = make(map[int]interface{}, 0)
		}
		curMap[hostID][mID] = valMap
		setIDs = append(setIDs, sID)
		posIDs = append(posIDs, mID)
	}
	moduels, err := h.getPoss(posIDs)
	if nil != err {
		return fmt.Errorf("HostPosConfigLog get pos error:%s", err.Error())
	}
	sets, err := h.getSets(setIDs)
	if nil != err {
		return fmt.Errorf("HostPosConfigLog get set error:%s", err.Error())
	}

	setMap := make(map[int]metadata.Ref, 0)
	for _, set := range sets {
		setInfo := set.(map[string]interface{})
		instID, _ := util.GetIntByInterface(setInfo[common.BKSetIDField])
		setMap[instID] = metadata.Ref{
			RefID:   instID,
			RefName: setInfo[common.BKSetNameField].(string),
		}
	}
	type PosRef struct {
		metadata.Ref
		Set     []interface{} `json:"set"`
		appID   interface{}
		ownerID string
	}
	posMap := make(map[int]PosRef, 0)
	for _, pos := range moduels {
		posInfo := pos.(map[string]interface{})
		mID, _ := util.GetIntByInterface(posInfo[common.BKPosIDField])
		sID, _ := util.GetIntByInterface(posInfo[common.BKSetIDField])
		posRef := PosRef{}
		posRef.Set = append(posRef.Set, setMap[sID])
		posRef.RefID = mID
		posRef.RefName = posInfo[common.BKPosNameField].(string)
		posRef.appID = posInfo[common.BKAppIDField]
		posRef.ownerID = posInfo[common.BKOwnerIDField].(string)
		posMap[mID] = posRef
	}
	posReName := "pos"
	setRefName := "set"
	headers := []metadata.Header{
		metadata.Header{PropertyID: posReName, PropertyName: "pos"},
		metadata.Header{PropertyID: setRefName, PropertyName: "app"},
		metadata.Header{PropertyID: common.BKAppIDField, PropertyName: "business ID"},
	}
	logs := []auditoplog.AuditLogExt{}

	for _, host := range h.hostInfos {
		host := host.(map[string]interface{})
		instID, _ := util.GetIntByInterface(host[common.BKHostIDField])
		log := auditoplog.AuditLogExt{ID: instID}
		log.ExtKey = host[common.BKHostInnerIPField].(string)

		prePos := make([]interface{}, 0)
		var preApp interface{}
		for posID, _ := range preMap[instID] {
			prePos = append(prePos, posMap[posID])
			preApp = posMap[posID].appID
			h.ownerID = posMap[posID].ownerID
		}

		curPos := make([]interface{}, 0)
		var curApp interface{}

		for posID, _ := range curMap[instID] {
			curPos = append(curPos, posMap[posID])
			curApp = posMap[posID].appID
			h.ownerID = posMap[posID].ownerID
		}

		log.Content = metadata.Content{
			PreData: common.KvMap{posReName: prePos, common.BKAppIDField: preApp},
			CurData: common.KvMap{posReName: curPos, common.BKAppIDField: curApp},
			Headers: headers,
		}
		logs = append(logs, log)

	}
	if "" == h.desc {
		h.desc = "host pos change"
	}
	opClient := auditlog.NewClient(h.auditCtrl)
	_, err = opClient.AuditHostsLog(logs, h.desc, h.ownerID, appID, user, auditoplog.AuditOpTypeHostModule)

	return err
}
