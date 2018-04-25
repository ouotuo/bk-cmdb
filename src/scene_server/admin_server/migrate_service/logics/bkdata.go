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
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/migrate_service/data"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

var prc2port = []string{"java:java:8008,8443", "nginx:nginx:80", "php-fpm:php-fpm:8009", "mysql:mysqld:3306", "mongodb:mongod:27017", "redis:redis-server:6379", "redis_cluster:redis-server:16379,6379",
	"zookeeper:java:2181", "kafka:java:9092", "elasticsearch:java:9300,10004", "beanstalkd:beanstalkd:6380", "influxdb:influxdb:5620,5621", "etcd:etcd:2379,2380", "consul:consul:8301,8300,8302,8500,53",
	"jumpserver:jumpserver:8080",
	"cmdb_adminserver:cmdb_adminserver:60004", "cmdb_apiserver:cmdb_apiserver:8080", "cmdb_auditcontroller:cmdb_auditcontroller:50005", "cmdb_datacollection:cmdb_datacollection:60005", "cmdb_eventserver:cmdb_eventserver:60009", "cmdb_hostcontroller:cmdb_hostcontroller:50002", "cmdb_hostserver:cmdb_hostserver:60001", "cmdb_objectcontroller:cmdb_objectcontroller:50001", "cmdb_proccontroller:cmdb_proccontroller:50003", "cmdb_procserver:cmdb_procserver:60003", "cmdb_toposerver:cmdb_toposerver:60002", "cmdb_webserver:cmdb_webserver:8083",
	"falcon-transfer:falcon-transfer:6060,8433", "falcon-hbs:falcon-hbs:6031", "falcon-judge:falcon-judge:6081,6080", "falcon-graph:falcon-graph:6071,6070", "falcon-nodata:falcon-nodata:6090", "falcon-aggregator:falcon-aggregator:6055", "falcon-api:falcon-api:8070", "falcon-alarm:falcon-alarm:9912", "falcon-task:falcon-task:8002"}

var setModuleKv = map[string]map[string]string{"CMDB配置管理": {"cmdb": "cmdb_adminserver,cmdb_apiserver,cmdb_auditcontroller,cmdb_datacollection,cmdb_eventserver,cmdb_hostcontroller,cmdb_hostserver,cmdb_objectcontroller,cmdb_proccontroller,cmdb_procserver,cmdb_toposerver,cmdb_webserver", "mongodb": "mongodb", "redis": "redis", "zookeeper": "zookeeper"},
	"监控告警平台": {"falcon+": "falcon-transfer,falcon-hbs,falcon-judge,falcon-graph,falcon-nodata,falcon-aggregator,falcon-api,falcon-alarm,falcon-task", "mysql": "mysql", "redis": "redis"},
	"跳板机":    {"jumpserver": "jumpserver"},
	"公共组件": {"mysql": "mysql", "redis": "redis", "redis_cluster": "redis_cluster", "zookeeper": "zookeeper", "kafka": "kafka", "elasticsearch": "elasticsearch",
		"nginx": "nginx", "beanstalk": "beanstalkd", "influxdb": "influxdb", "etcd": "etcd", "consul": "consul"}}

var appID int = 0
var ownerID string = common.BKDefaultOwnerID
var procAPI, topoAPI string
var procName2ID map[string]int
var appModelData map[string]interface{}
var setModelData map[string]interface{}
var moduleModelData map[string]interface{}
var procModelData map[string]interface{}

//BKAppInit  init bk app
func BKAppInit(req *restful.Request, cc *api.APIResource, ownerID string) error {
	var err error
	//get api addr
	procAPI = cc.ProcAPI()
	topoAPI = cc.TopoAPI()
	//get model module
	procModelData, err = getObjectFields(cc.TopoAPI(), req, common.BKInnerObjIDProc)
	if err != nil {
		blog.Error("get procModelData err :%v ", err)
		return err
	}
	appModelData, err = getObjectFields(cc.TopoAPI(), req, common.BKInnerObjIDApp)
	if err != nil {
		blog.Error("get appModelData err :%v ", err)
		return err
	}
	setModelData, err = getObjectFields(cc.TopoAPI(), req, common.BKInnerObjIDSet)
	if err != nil {
		blog.Error("get setModelData err :%v ", err)
		return err
	}
	moduleModelData, err = getObjectFields(cc.TopoAPI(), req, common.BKInnerObjIDModule)
	if err != nil {
		blog.Error("get moduleModelData err :%v ", err)
		return err
	}

	isExist, err := BKAppIsExist(req)
	if nil != err {
		blog.Error("get app isExist err :%v ", err)
		return err
	}

	if !isExist {
		err = addBKApp(req)
		if nil != err {
			blog.Error("add bk app err :%v ", err)
			return err
		}

		err = addBKProcess(req)
		if nil != err {
			blog.Error("add bk process err :%v ", err)
			return err
		}
	}
	return nil

}

//addBKApp add bk app
func addBKApp(req *restful.Request) error {
	appModelData[common.BKAppNameField] = common.BKAppName
	appModelData[common.BKMaintainersField] = "admin"

	if data.Distribution == common.RevisionEnterprise {
		appModelData[common.BKTimeZoneField] = "Asia/Shanghai"
		appModelData[common.BKLanguageField] = "1" //"中文"
	} else {
		delete(appModelData, common.BKTimeZoneField)
		delete(appModelData, common.BKLanguageField)
	}
	appModelData[common.BKLifeCycleField] = common.DefaultAppLifeCycleNormal

	byteParams, _ := json.Marshal(appModelData)
	url := topoAPI + "/topo/v1/app/" + ownerID
	blog.Info("migrate add bk app url :%s", url)
	blog.Info("migrate add bk app content :%s", string(byteParams))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPCreate, byteParams)
	blog.Info("migrate add bk app return :%s", string(reply))
	if err != nil {
		blog.Error("add bk app err :%v ", err)
		return err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()

	code, err := util.GetIntByInterface(output["bk_error_code"])
	if err != nil || 0 != code {
		if strings.Contains(reply, "duplicate") || strings.Contains(reply, "重复") {
			return nil
		}
		blog.Error("add bk app json err :%v ", err)
		return errors.New(reply)
	}
	data, ok := output["data"].(map[string]interface{})
	if false == ok {
		blog.Error("add bk app result data err :%v ", err)
		return errors.New("get appID error")
	}
	appID, err = util.GetIntByInterface(data[common.BKAppIDField])
	if nil != err {
		blog.Error("add bk app result data app id err :%v ", err)
		return err
	}
	return nil
}

//addBKProcess add bk process
func addBKProcess(req *restful.Request) error {
	procName2ID = make(map[string]int)
	appIDStr := strconv.Itoa(appID)

	for _, procStr := range prc2port {
		procArr := strings.Split(procStr, ":")
		procName := procArr[0]
		funcName := procArr[1]
		portStr := procArr[2]
		procModelData[common.BKProcNameField] = procName
		procModelData[common.BKFuncName] = funcName
		procModelData[common.BKPort] = portStr
		//procModelData[common.BKWorkPath] = "/data/bkee"
		byteParams, _ := json.Marshal(procModelData)
		url := procAPI + "/process/v1/" + ownerID + "/" + appIDStr
		blog.Info("migrate add process url :%s", url)
		blog.Info("migrate add process content :%s", string(byteParams))
		reply, err := httpcli.ReqHttp(req, url, common.HTTPCreate, byteParams)
		blog.Info("migrate add process return :%s", string(reply))
		if err != nil {
			blog.Error("add process err :%v ", err)
			procName2ID[procName] = 0
			continue
		}
		js, err := simplejson.NewJson([]byte(reply))
		if nil != err {
			blog.Error("add bk process data return not json err :%v ", err)
			return err
		}
		output, err := js.Map()
		if nil != err {
			blog.Error("add bk process data return not json err :%v ", err)
			return err
		}
		code, err := util.GetIntByInterface(output["bk_error_code"])
		if err != nil || 0 != code {
			blog.Error("add process code err :%v ", err)
			continue
		}
		data, ok := output["data"].(map[string]interface{})
		if false == ok {
			blog.Error("add process data err :%v ", err)
			continue
		}
		procIDi, ok := data[common.BKProcIDField]
		if false == ok {
			blog.Error("add process data process ID err :%v ", err)
			continue
		}
		procID, err := util.GetIntByInterface(procIDi)
		if nil != err {
			continue
		}
		procName2ID[procName] = procID
	}
	addSetInBKApp(req)
	return nil
}

//addSetInBKApp add set in bk app
func addSetInBKApp(req *restful.Request) {
	appIDStr := strconv.Itoa(appID)
	for setName, moduleArr := range setModuleKv {
		setModelData[common.BKSetNameField] = setName
		setModelData[common.BKAppIDField] = appID
		setModelData[common.BKOwnerIDField] = common.BKDefaultOwnerID
		setModelData[common.BKInstParentStr] = appID
		byteParams, _ := json.Marshal(setModelData)
		url := topoAPI + "/topo/v1/set" + "/" + appIDStr
		blog.Info("migrate add set url :%s", url)
		blog.Info("migrate add set content :%s", string(byteParams))
		reply, err := httpcli.ReqHttp(req, url, common.HTTPCreate, byteParams)
		blog.Info("migrate add set return :%s", string(reply))
		if err != nil {
			blog.Error("add set data err :%v ", err)
			continue
		}
		js, _ := simplejson.NewJson([]byte(reply))
		output, _ := js.Map()

		code, err := util.GetIntByInterface(output["bk_error_code"])
		if err != nil || 0 != code {
			blog.Error("add set data code err :%v ", err)
			continue
		}
		data, ok := output["data"].(map[string]interface{})
		if false == ok {
			blog.Error("add set data result err :%v ", err)
			continue
		}
		setIDi, ok := data[common.BKSetIDField]
		if false == ok {
			continue
		}
		setID, err := util.GetIntByInterface(setIDi)
		if nil != err {
			continue
		}
		// add module in set
		addModuleInSet(req, moduleArr, setID)
	}
}

//addModuleInSet add module in set
func addModuleInSet(req *restful.Request, moduleArr map[string]string, setID int) {
	appIDStr := strconv.Itoa(appID)
	for moduleName, processNameStr := range moduleArr {
		moduleModelData[common.BKModuleNameField] = moduleName
		moduleModelData[common.BKAppIDField] = appID
		moduleModelData[common.BKSetIDField] = setID
		moduleModelData[common.BKOwnerIDField] = common.BKDefaultOwnerID
		moduleModelData[common.BKInstParentStr] = setID
		setIDStr := strconv.Itoa(setID)
		byteParams, _ := json.Marshal(moduleModelData)
		url := topoAPI + "/topo/v1/module" + "/" + appIDStr + "/" + setIDStr
		blog.Info("migrate add module url :%s", url)
		blog.Info("migrate add module content :%s", string(byteParams))
		reply, err := httpcli.ReqHttp(req, url, common.HTTPCreate, byteParams)
		blog.Info("migrate add module return :%s", string(reply))
		if err != nil {
			continue
		}
		js, _ := simplejson.NewJson([]byte(reply))
		output, _ := js.Map()

		code, err := util.GetIntByInterface(output["bk_error_code"])
		if err != nil || 0 != code {
			continue
		}
		//add module process config
		addModule2Process(req, processNameStr, moduleName)
	}
}

//addModule2Process add process 2 module
func addModule2Process(req *restful.Request, processNameStr string, moduleName string) {
	appIDStr := strconv.Itoa(appID)
	processNameArr := strings.Split(processNameStr, ",")
	for _, processName := range processNameArr {
		processID, ok := procName2ID[processName]
		if false == ok {
			continue
		}
		processIDStr := strconv.Itoa(processID)
		url := procAPI + "/process/v1/module" + "/" + common.BKDefaultOwnerID + "/" + appIDStr + "/" + processIDStr + "/" + moduleName
		blog.Info("migrate add module process config url :%s", url)
		reply, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, nil)
		if err != nil {
			blog.Error("migrate add module process config %v", err)
			continue
		}
		blog.Info("migrate add module process return :%s", string(reply))
		js, err := simplejson.NewJson([]byte(reply))
		if err != nil {
			blog.Error("migrate add module process config json err %v", err)
			continue
		}
		output, err := js.Map()
		if err != nil {
			blog.Error("migrate add module process config data not map err %v", err)
			continue
		}

		code, err := util.GetIntByInterface(output["bk_error_code"])
		if err != nil || 0 != code {
			blog.Error("migrate add module process config code err %v", err)
			continue
		}
	}
}

//BKAppIsExist is bk app exist
func BKAppIsExist(req *restful.Request) (bool, error) {

	params := make(map[string]interface{})
	conditon := make(map[string]interface{})
	conditon[common.BKAppNameField] = common.BKAppName
	conditon[common.BKOwnerIDField] = ownerID
	params["condition"] = conditon
	params["fields"] = []string{common.BKAppIDField}
	params["start"] = 0
	params["limit"] = 20

	byteParams, _ := json.Marshal(params)
	url := topoAPI + "/topo/v1/app/search/" + ownerID
	blog.Info("Get bk app url :%s", url)
	blog.Info("Get bk app content :%s", string(byteParams))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, byteParams)
	if err != nil {
		blog.Error("Get bk app error :%v", err)
		return false, err
	}
	blog.Info("Get bk app return :%s", string(reply))
	js, err := simplejson.NewJson([]byte(reply))
	if nil != err {
		blog.Error("Get bk app data not json error :%v", err)
	}
	output, err := js.Map()
	if nil != err {
		blog.Error("Get bk app data not map error :%v", err)
	}
	code, err := util.GetIntByInterface(output["bk_error_code"])
	if err != nil {
		blog.Error("Get bk app data not map error :%v", err)
		return false, errors.New(reply)
	}
	if 0 != code {
		blog.Error("Get bk app data not map error :%v", err)
		return false, errors.New(output["message"].(string))
	}
	cnt, err := js.Get("data").Get("count").Int()
	if err != nil {
		blog.Error("Get bk app data not count error :%v", err)
		return false, errors.New(reply)
	}
	if 0 == cnt {
		return false, nil
	}
	return true, nil
}
