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
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	parse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/source_controller/common/commondata"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

//get posid by cond
func GetPosIDByCond(req *restful.Request, objURL string, cond []interface{}) ([]int, error) {
	posIDArr := make([]int, 0)
	condition := make(map[string]interface{})
	condition["fields"] = common.BKPosIDField
	condition["sort"] = common.BKPosIDField
	condition["start"] = 0
	condition["limit"] = 0
	condc := make(map[string]interface{})
	parse.ParseCommonParams(cond, condc)
	condition["condition"] = condc
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/pos/search"
	blog.Infof("GetPosIDByCond url :%s content:%s", url, string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetPosIDByCond return :%s", string(reply))
	if err != nil {
		return posIDArr, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	posData := output["data"]
	posResult, ok := posData.(map[string]interface{})
	if !ok {
		return posIDArr, nil
	}
	posInfo, ok := posResult["info"].([]interface{})
	if !ok {
		return posIDArr, nil
	}
	for _, i := range posInfo {
		pos := i.(map[string]interface{})
		posID, _ := pos[common.BKPosIDField].(json.Number).Int64()
		posIDArr = append(posIDArr, int(posID))
	}
	return posIDArr, nil
}

//get posmap by cond
func GetPosMapByCond(req *restful.Request, fields string, objURL string, cond interface{}) (map[int]interface{}, error) {
	posMap := make(map[int]interface{})
	condition := make(map[string]interface{})
	condition["fields"] = fields
	condition["sort"] = common.BKPosIDField
	condition["start"] = 0
	condition["limit"] = 0
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/pos/search"
	blog.Info("GetPosMapByCond url :%s", url)
	blog.Info("GetPosMapByCond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetPosMapByCond return :%s", string(reply))
	if err != nil {
		return posMap, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	posData := output["data"]
	posResult := posData.(map[string]interface{})
	posInfo := posResult["info"].([]interface{})
	for _, i := range posInfo {
		pos := i.(map[string]interface{})
		posID, _ := pos[common.BKPosIDField].(json.Number).Int64()
		posMap[int(posID)] = i
	}
	return posMap, nil
}

//GetPosByPosID  get pos by pos id
func GetPosByPosID(req *restful.Request, appID int, posID int, hostAddr string) ([]interface{}, error) {
	URL := hostAddr + "/object/v1/insts/pos/search"
	params := make(map[string]interface{})

	conditon := make(map[string]interface{})
	conditon[common.BKIdcIDField] = appID
	conditon[common.BKPosIDField] = posID
	params["condition"] = conditon
	params["sort"] = common.BKPosIDField
	params["start"] = 0
	params["limit"] = 1
	params["fields"] = common.BKPosIDField
	isSuccess, errMsg, data := GetHttpResult(req, URL, common.HTTPSelectPost, params)
	if !isSuccess {
		blog.Error("get idle pos error, params:%v, error:%s", params, errMsg)
		return nil, errors.New(errMsg)
	}
	dataStrArry := data.(map[string]interface{})
	dataInfo, ok := dataStrArry["info"].([]interface{})
	if !ok {
		blog.Error("get idle pos error, params:%v, error:%s", params, errMsg)
		return nil, errors.New(errMsg)
	}

	return dataInfo, nil
}

//GetSinglePosID  get single pos id
func GetSinglePosID(req *restful.Request, conds interface{}, hostAddr string) (int, error) {
	//posURL := "http://" + cc.ObjCtrl + "/object/v1/insts/pos/search"
	url := hostAddr + "/object/v1/insts/pos/search"
	params := make(map[string]interface{})

	params["condition"] = conds
	params["sort"] = common.BKPosIDField
	params["start"] = 0
	params["limit"] = 1
	params["fields"] = common.BKPosIDField
	isSuccess, errMsg, data := GetHttpResult(req, url, common.HTTPSelectPost, params)
	if !isSuccess {
		blog.Error("get idle pos error, params:%v, error:%s", params, errMsg)
		return 0, errors.New(errMsg)
	}
	dataInterface := data.(map[string]interface{})
	info := dataInterface["info"].([]interface{})
	if 1 != len(info) {
		blog.Error("not find pos error, params:%v, error:%s", params, errMsg)
		return 0, errors.New("获取集群，返回数据格式错误")
	}
	row := info[0].(map[string]interface{})
	posID, _ := util.GetIntByInterface(row[common.BKPosIDField])

	if 0 == posID {
		blog.Error("not find pos error, params:%v, error:%s", params, errMsg)
		return 0, errors.New("获取集群信息失败")
	}

	return posID, nil
}

// NewHostSyncValidPos
// 1. check pos is exist,
// 2. multiple posID  Check whether  all the pos default 0
func NewHostSyncValidPos(req *restful.Request, appID int, posID []int, objAddr string) ([]int, error) {
	if 0 == len(posID) {
		return nil, fmt.Errorf("pos id number must be > 1")
	}

	conds := common.KvMap{
		common.BKIdcIDField:    appID,
		common.BKPosIDField: common.KvMap{common.BKDBIN: posID},
	}

	condition := new(commondata.ObjQueryInput)
	condition.Sort = common.BKPosIDField
	condition.Start = 0
	condition.Limit = 0
	condition.Condition = conds
	bodyContent, err := json.Marshal(condition)
	if nil != err {
		return nil, fmt.Errorf("query pos parameters not json")
	}
	url := objAddr + "/object/v1/insts/pos/search"
	blog.Info("NewHostSyncValidPos url :%s", url)
	blog.Info("NewHostSyncValidPos content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("NewHostSyncValidPos return :%s", string(reply))

	js, err := simplejson.NewJson([]byte(reply))
	if nil != err {
		blog.Errorf("NewHostSyncValidPos get pos reply not json,  url:%s, params:%s, reply:%s", string(bodyContent), url, reply)
		return nil, fmt.Errorf("get moduel reply not json, reply:%s", reply)
	}

	posInfos, err := js.Get("data").Get("info").Array()
	if nil != err {
		blog.Errorf("NewHostSyncValidPos get pos reply not foound data.info,  url:%s, params:%s, reply:%s", string(bodyContent), url, reply)
		return nil, fmt.Errorf("get moduel reply not found info from data, reply:%s", reply)
	}

	// only pos  and pos exist return true
	if 1 == len(posID) && 1 == len(posInfos) {
		return posID, nil
	}

	// use pos id is exist
	posIDMap := make(map[int64]bool)
	for _, id := range posID {
		posIDMap[int64(id)] = false
	}

	// multiple pos all pos default = 0
	for _, pos := range posInfos {
		posMap, ok := pos.(map[string]interface{})
		if !ok {
			blog.Errorf("NewHostSyncValidPos item not map[string]interface{},  pos:%v", pos)
			return nil, fmt.Errorf("item not map[string]interface{},  pos:%v", pos)
		}
		moduelDefault, err := util.GetInt64ByInterface(posMap[common.BKDefaultField])
		if nil != err {
			blog.Errorf("NewHostSyncValidPos item not found default,  pos:%v", pos)
			return nil, fmt.Errorf("pos information not found default,  pos:%v", pos)
		}
		if 0 != moduelDefault {
			return nil, fmt.Errorf("multiple pos cannot appear system pos")
		}
		posID, err := util.GetInt64ByInterface(posMap[common.BKPosIDField])
		if nil != err {
			blog.Errorf("NewHostSyncValidPos item not found pos id,  pos:%v", pos)
			return nil, fmt.Errorf("pos information not found pos id,  pos:%v", pos)
		}
		posIDMap[posID] = true //pos id  exist db

	}

	var dbPosID []int
	var notExistPosId []string
	for id, exist := range posIDMap {
		if exist {
			dbPosID = append(dbPosID, int(id))
		} else {
			notExistPosId = append(notExistPosId, strconv.FormatInt(id, 10))
		}
	}
	if 0 < len(notExistPosId) {
		return nil, fmt.Errorf("pos id %s not found", strings.Join(notExistPosId, ","))
	}

	return dbPosID, nil
}
