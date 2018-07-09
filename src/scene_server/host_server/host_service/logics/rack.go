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
	"encoding/json"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

//GetRackIDByCond get rackid by cond
func GetRackIDByCond(req *restful.Request, objURL string, cond []interface{}) ([]int, error) {
	rackIDArr := make([]int, 0)
	condition := make(map[string]interface{})
	condition["fields"] = common.BKRackIDField
	condition["sort"] = common.BKRackIDField
	condition["start"] = 0
	condition["limit"] = 0
	condc := make(map[string]interface{})
	parse.ParseCommonParams(cond, condc)
	condition["condition"] = condc
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/rack/search"
	blog.Infof("GetRackIDByCond url :%s content:%s", url, string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetrackIDByCond return :%s", string(reply))
	if err != nil {
		return rackIDArr, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	rackData := output["data"]
	rackResult, ok := rackData.(map[string]interface{})
	if !ok {
		return rackIDArr, nil
	}
	rackInfo, ok := rackResult["info"].([]interface{})
	if !ok {
		return rackIDArr, nil
	}
	for _, i := range rackInfo {
		rack := i.(map[string]interface{})
		rackID, _ := rack[common.BKRackIDField].(json.Number).Int64()
		rackIDArr = append(rackIDArr, int(rackID))
	}
	return rackIDArr, nil
}

//GetRackMapByCond get rackmap by cond
func GetRackMapByCond(req *restful.Request, fields string, objURL string, cond interface{}) (map[int]interface{}, error) {
	rackMap := make(map[int]interface{})
	condition := make(map[string]interface{})
	condition["fields"] = fields
	condition["sort"] = common.BKModuleIDField
	condition["start"] = 0
	condition["limit"] = 0
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/rack/search"
	blog.Info("GetRackMapByCond url :%s", url)
	blog.Info("GetRackMapByCond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetRackMapByCond return :%s", string(reply))
	if err != nil {
		return rackMap, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	rackData := output["data"]
	rackResult := rackData.(map[string]interface{})
	rackInfo := rackResult["info"].([]interface{})
	for _, i := range rackInfo {
		rack := i.(map[string]interface{})
		rackID, _ := rack[common.BKRackIDField].(json.Number).Int64()
		rackMap[int(rackID)] = i
	}
	return rackMap, nil
}
