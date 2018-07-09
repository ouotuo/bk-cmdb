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

package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/language"
	idcParse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"encoding/json"
	"errors"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

//GetIdcIDByCond get idcid by cond
func GetIdcIDByCond(req *restful.Request, objURL string, cond []interface{}) ([]int, error) {
	idcIDArr := make([]int, 0)
	condition := make(map[string]interface{})
	condition["fields"] = common.BKIdcIDField
	condition["sort"] = common.BKIdcIDField
	condition["start"] = 0
	condition["limit"] = 1000000
	condc := make(map[string]interface{})
	idcParse.ParseCommonParams(cond, condc)
	condition["condition"] = condc
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/" + common.BKInnerObjIDIdc + "/search"
	blog.Info("GetIdcIDByCond url :%s", url)
	blog.Info("GetIdcIDByCond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetIdcIDByCond return :%s", string(reply))
	if err != nil {
		return idcIDArr, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	idcData := output["data"]
	idcResult := idcData.(map[string]interface{})
	idcInfo := idcResult["info"].([]interface{})
	for _, i := range idcInfo {
		idc := i.(map[string]interface{})
		idcID, _ := idc[common.BKIdcIDField].(json.Number).Int64()
		idcIDArr = append(idcIDArr, int(idcID))
	}
	return idcIDArr, nil
}

//GetIdcMapByCond get idcmap by cond
func GetIdcMapByCond(req *restful.Request, fields string, objURL string, cond interface{}) (map[int]interface{}, error) {
	idcMap := make(map[int]interface{})
	condition := make(map[string]interface{})
	condition["fields"] = fields
	condition["sort"] = common.BKIdcIDField
	condition["start"] = 0
	condition["limit"] = 0
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/" + common.BKInnerObjIDIdc + "/search"
	blog.Info("GetIdcMapByCond url :%s", url)
	blog.Info("GetIdcMapByCond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetIdcMapByCond return :%s", string(reply))
	if err != nil {
		blog.Errorf("GetIdcMapByCond params:%s  error:%v", string(bodyContent), err)
		return idcMap, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	idcData := output["data"]
	idcResult := idcData.(map[string]interface{})
	idcInfo := idcResult["info"].([]interface{})
	for _, i := range idcInfo {
		idc := i.(map[string]interface{})
		idcID, _ := idc[common.BKIdcIDField].(json.Number).Int64()
		idcMap[int(idcID)] = i
	}
	return idcMap, nil
}

//GetSingleIdc  get single idc
func GetSingleIdc(req *restful.Request, objURL string, cond interface{}) (map[string]interface{}, error) {
	condition := make(map[string]interface{})
	condition["sort"] = common.BKIdcIDField
	condition["start"] = 0
	condition["limit"] = 1
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/" + common.BKInnerObjIDIdc + "/search"
	fmt.Println("GetSingleIdc", url, string(bodyContent))

	blog.Info("GetOneIdc url :%s", url)
	blog.Info("GetOneIdc content :%s", string(bodyContent))

	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	fmt.Println("GetSingleIdc", url, string(reply))
	blog.Info("GetOneIdc return :%s", string(reply))
	if err != nil {
		blog.Info("GetOneIdc return http request error:%s", string(reply))
		return nil, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	idcData := output["data"]
	idcResult := idcData.(map[string]interface{})
	idcInfo := idcResult["info"].([]interface{})
	for _, i := range idcInfo {
		idc, _ := i.(map[string]interface{})
		return idc, nil
	}
	return nil, nil
}

//GetIdcInfo get idc info
func GetIdcInfo(req *restful.Request, fields string, conditon map[string]interface{}, objAddr string, defLang language.DefaultCCLanguageIf) (map[string]interface{}, error) {
	//moduleURL := "http://" + cc.ObjCtrl + "/object/v1/insts/module/search"
	URL := objAddr + "/object/v1/insts/" + common.BKInnerObjIDIdc + "/search"
	params := make(map[string]interface{})
	params["condition"] = conditon
	params["sort"] = common.BKIdcIDField
	params["start"] = 0
	params["limit"] = 1
	params["fields"] = fields

	blog.Info("get idclication info  url:%s", URL)
	blog.Info("get idclication info  url:%v", params)
	isSuccess, errMsg, data := GetHttpResult(req, URL, common.HTTPSelectPost, params)
	if !isSuccess {
		blog.Error("get idclication info  error, params:%v, error:%s", params, errMsg)
		return nil, errors.New(errMsg)
	}
	dataInterface := data.(map[string]interface{})
	info := dataInterface["info"].([]interface{})
	if 1 != len(info) {
		blog.Error("not idclication info error, params:%v, error:%s", params, errMsg)
		return nil, errors.New(defLang.Languagef("idc_not_exist")) //"业务不存在")
	}
	row := info[0].(map[string]interface{})

	if 0 == len(row) {
		blog.Error("not idclication info error, params:%v, error:%s", params, errMsg)
		return nil, errors.New(defLang.Languagef("idc_not_exist")) //"业务不存在")
	}

	return row, nil
}

//GetDefaultIdcID get default biz id
func GetDefaultIdcID(req *restful.Request, ownerID, fields, objAddr string, defLang language.DefaultCCLanguageIf) (int, error) {
	conds := make(map[string]interface{})
	conds[common.BKOwnerIDField] = ownerID
	conds[common.BKDefaultField] = common.DefaultIdcFlag
	idcinfo, err := GetIdcInfo(req, fields, conds, objAddr, defLang)
	if nil != err {
		blog.Errorf("get default idc info error:%v", err.Error())
		return 0, err
	}
	return util.GetIntByInterface(idcinfo[common.BKIdcIDField])
}

//GetDefaultIdcID get supplier ID
func GetDefaultIdcIDBySupplierID(req *restful.Request, supplierID int, fields, objAddr string, defLang language.DefaultCCLanguageIf) (int, error) {
	conds := make(map[string]interface{})
	conds[common.BKSupplierIDField] = supplierID
	conds[common.BKDefaultField] = common.DefaultIdcFlag
	idcinfo, err := GetIdcInfo(req, fields, conds, objAddr, defLang)
	if nil != err {
		blog.Errorf("get default idc info error:%v", err.Error())
		return 0, err
	}
	return util.GetIntByInterface(idcinfo[common.BKIdcIDField])
}

// IsExistHostIDInIdc  is host exsit in idc
func IsExistHostIDInIdc(CC *api.APIResource, req *restful.Request, idcID int, hostID int, defLang language.DefaultCCLanguageIf) (bool, error) {
	conds := common.KvMap{common.BKIdcIDField: idcID, common.BKHostIDField: hostID}
	url := CC.HostCtrl() + "/host/v1/meta/hosts/pos/search"
	isSucess, errmsg, data := GetHttpResult(req, url, common.HTTPSelectPost, conds)
	blog.Info("IsExistHostIDInIdc request url:%s, params:{idcid:%d, hostid:%d}", url, idcID, hostID)
	blog.Info("IsExistHostIDInIdc res:%v,%s, %v", isSucess, errmsg, data)
	if !isSucess {
		return false, errors.New(defLang.Languagef("host_search_module_fail_with_errmsg", errmsg)) //"获取主机关系失败;" + errmsg)
	}
	//数据为空
	if nil == data {
		return false, nil
	}
	ids, ok := data.([]interface{})
	if !ok {
		return false, errors.New(defLang.Languagef("host_search_module_fail_with_errmsg", errmsg)) //"获取主机关系失败;" + errmsg)
	}

	if len(ids) > 0 {
		return true, nil
	}
	return false, nil

}
