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
	"configcenter/src/common/http/httpclient"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bitly/go-simplejson"
)

// getObjFieldIDs get the values of properyID and properyName
func GetObjFieldIDs(objID, url string, header http.Header) (common.KvMap, error) {
	conds := common.KvMap{common.BKObjIDField: objID, common.BKOwnerIDField: common.BKDefaultOwnerID, "page": common.KvMap{"skip": 0, "limit": common.BKNoLimit}}
	result, err := httpRequest(url, conds, header)
	if nil != err {
		blog.Errorf("get %s fields error:%s", objID, err.Error())
		return nil, err
	}
	blog.Info("get %s fields  url:%s", objID, url)
	blog.Info("get %s fields return:%s", objID, result)
	js, _ := simplejson.NewJson([]byte(result))
	hostFields, _ := js.Map()
	fields, _ := hostFields["data"].([]interface{})
	ret := common.KvMap{}

	for _, field := range fields {
		mapField, _ := field.(map[string]interface{})

		fieldName, _ := mapField[common.BKPropertyNameField].(string)

		blog.Debug("fieldName:%v", fieldName)
		fieldId, _ := mapField[common.BKPropertyIDField].(string)
		propertyType, _ := mapField[common.BKPropertyTypeField].(string)

		blog.Debug("fieldId:%v", fieldId)
		ret[fieldId] = common.KvMap{"name": fieldName, "type": propertyType, "require": mapField[common.BKIsRequiredField]}
	}

	return ret, nil
}

// AutoInputV3Field fields required to automatically populate the current object v3
func AutoInputV3Field(params common.KvMap, objId, url string, header http.Header) (common.KvMap, error) {
	appFields, err := GetObjFieldIDs(objId, url+"/topo/v1/objectattr/search", header)
	if nil != err {
		blog.Error("CreateApp error:%s", err.Error())
		//		converter.RespFailV2(common.CC_Err_Comm_APP_Create_FAIL, common.CC_Err_Comm_APP_Create_FAIL_STR, resp)

		return nil, errors.New("CC_Err_Comm_APP_Create_FAIL_STR")
	}
	for fieldId, item := range appFields {
		mapItem, _ := item.(common.KvMap)
		_, ok := params[fieldId]
		if !ok {
			strType, _ := mapItem["type"].(string)
			if common.FieldTypeLongChar == strType || common.FieldTypeSingleChar == strType {
				params[fieldId] = ""
			} else {
				params[fieldId] = nil

			}
		}
	}

	return params, nil
}

// httpRequest http request
func httpRequest(url string, body interface{}, header http.Header) (string, error) {
	params, _ := json.Marshal(body)
	blog.Info("input:%s", string(params))
	httpClient := httpclient.NewHttpClient()
	httpClient.SetHeader("Content-Type", "application/json")
	httpClient.SetHeader("Accept", "application/json")
	reply, err := httpClient.POST(url, header, params)
	return string(reply), err
}
