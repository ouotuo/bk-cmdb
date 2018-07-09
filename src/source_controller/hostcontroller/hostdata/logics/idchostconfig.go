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
	"configcenter/src/common/util"
	eventtypes "configcenter/src/scene_server/event_server/types"
	metadataTable "configcenter/src/source_controller/api/metadata"
	"configcenter/src/source_controller/common/eventdata"
	"configcenter/src/source_controller/common/instdata"
	"errors"

	"gopkg.in/mgo.v2/bson"
)

var (
	posBaseTaleName = "cc_PosBase"
)

type posHostConfigParams struct {
	IdcId int   `json:"bk_idc_id"`
	HostID        int   `json:"bk_host_id"`
	PosID      []int `json:"bk_pos_id"`
}

//DelSingleHostPosRelation delete single host pos relation
func DelSingleHostPosRelation(ec *eventdata.EventContext, cc *api.APIResource, hostID, posID, appID int) (bool, error) {

	//get host info
	hostFieldArr := []string{common.BKHostInnerIPField}
	hostResult := make(map[string]interface{}, 0)
	errHost := instdata.GetObjectByID(common.BKInnerObjIDHost, hostFieldArr, hostID, &hostResult, common.BKHostIDField)
	blog.Infof("DelSingleHostPosRelation hostID:%d, hostinfo:%v", hostID, hostResult)
	if errHost != nil {
		blog.Error("delSingleHostPosRelation get host error:%s, host:%v", errHost.Error(), hostID)
		return false, errHost
	}

	posFieldArr := []string{common.BKPosNameField}
	var posResult interface{}
	errPos := instdata.GetObjectByID(common.BKInnerObjIDPos, posFieldArr, posID, &posResult, common.BKPosNameField)
	blog.Infof("DelSingleHostPosRelation pos:%d, pos info:%v", posID, posResult)
	if errPos != nil {
		blog.Error("delSingleHostPosRelation get pos posID:%d, error:%s,", posID, errPos.Error())
		return false, errPos
	}

	tableName := metadataTable.IdcHostConfig{}

	delCondition := make(map[string]interface{})
	delCondition[common.BKIdcIDField] = appID
	delCondition[common.BKHostIDField] = hostID
	delCondition[common.BKPosIDField] = posID
	num, numError := cc.InstCli.GetCntByCondition(tableName.TableName(), delCondition)
	blog.Infof("DelSingleHostPosRelation  get pos host relation condition:%v", delCondition)
	if numError != nil {
		blog.Error("delSingleHostPosRelation get pos host relation error:", numError.Error())
		return false, numError
	}
	//no config, return
	if num == 0 {
		return true, nil
	}

	// retrieve original datas
	origindatas := make([]map[string]interface{}, 0)
	getErr := cc.InstCli.GetMutilByCondition(tableName.TableName(), nil, delCondition, &origindatas, "", 0, 0)
	if getErr != nil {
		blog.Error("retrieve original datas error:%v", getErr)
		return false, getErr
	}

	delErr := cc.InstCli.DelByCondition(tableName.TableName(), delCondition)
	blog.Infof("DelSingleHostPosRelation delCondition:%v", delCondition)
	if delErr != nil {
		blog.Error("delSingleHostPosRelation del pos host relation error:", delErr.Error())
		return false, delErr
	}

	// send events
	for _, origindata := range origindatas {
		err := ec.InsertEvent(eventtypes.EventTypeRelation, "postransfer", eventtypes.EventActionDelete, nil, origindata)
		if err != nil {
			blog.Error("create event error:%v", err)
		}
	}

	return true, nil
}

//AddSingleHostPosRelation add single host pos relation
func AddSingleHostPosRelation(ec *eventdata.EventContext, cc *api.APIResource, hostID, posID, appID int) (bool, error) {
	//get host info
	hostFieldArr := []string{common.BKHostInnerIPField}
	hostResult := make(map[string]interface{})

	errHost := instdata.GetObjectByID(common.BKInnerObjIDHost, hostFieldArr, hostID, &hostResult, common.BKHostIDField)
	if errHost != nil {
		blog.Error("addSingleHostPosRelation get host error:%s", errHost.Error())
		return false, errHost
	}

	posFieldArr := []string{common.BKPosNameField, common.BKRackIDField}
	posResult := make(map[string]interface{})
	errPos := instdata.GetObjectByID(common.BKInnerObjIDPos, posFieldArr, posID, &posResult, common.BKPosIDField)
	if errPos != nil {
		blog.Error("addSingleHostPosRelation get pos posid:%d, error:%s", posID, errPos.Error())
		return false, errPos
	}
	blog.Error("error:%s", posResult)
	posName, _ := posResult[common.BKPosNameField].(string)
	rackID, _ := util.GetIntByInterface(posResult[common.BKRackIDField])

	if "" == posName || 0 == rackID {
		blog.Error("addSingleHostPosRelation get pos error:not find pos width PosID:%d", posID)
		return false, errors.New("未找到对应的模块")
	}

	tableName := metadataTable.IdcHostConfig{}
	posHostConfig := make(map[string]interface{})

	posHostConfig[common.BKIdcIDField] = appID
	posHostConfig[common.BKHostIDField] = hostID
	posHostConfig[common.BKPosIDField] = posID

	num, numError := cc.InstCli.GetCntByCondition(tableName.TableName(), posHostConfig)
	if numError != nil {
		blog.Error("addSingleHostPosRelation get pos host relation error:", numError.Error())
		return false, numError
	}
	//config exsit, return
	if num > 0 {
		return true, nil
	}

	posHostConfig[common.BKRackIDField] = rackID
	_, err := cc.InstCli.Insert(tableName.TableName(), posHostConfig)
	if err != nil {
		blog.Error("addSingleHostPosRelation add pos host relation error:", err.Error())
		return false, err
	}

	err = ec.InsertEvent(eventtypes.EventTypeRelation, "postransfer", eventtypes.EventActionCreate, posHostConfig, nil)
	if err != nil {
		blog.Error("create event error:%v", err)
	}

	return true, nil
}

//GetDefaultPosIDs get default pos ids
func GetDefaultPosIDs(cc *api.APIResource, appID int) ([]int, error) {
	defaultPosCond := make(map[string]interface{}, 2)
	defaultPosCond[common.BKDefaultField] = common.KvMap{common.BKDBIN: []int{common.DefaultFaultPosFlag, common.DefaultResPosFlag}}
	defaultPosCond[common.BKIdcIDField] = appID
	result := make([]interface{}, 0)
	var ret []int

	err := cc.InstCli.GetMutilByCondition(posBaseTaleName, []string{common.BKPosIDField, common.BKDefaultField}, defaultPosCond, &result, "ID", 0, 100)
	blog.Infof("defaultPosCond:%v", defaultPosCond)
	if nil != err {
		blog.Errorf("getDefaultPosIds error:%s, params:%v, %v", err.Error(), defaultPosCond, result)
		return ret, errors.New("未找到模块")
	}

	for _, r := range result {
		item := r.(bson.M)
		ID, err := util.GetIntByInterface(item[common.BKPosIDField])
		if nil != err {
			return ret, errors.New("未找到模块")
		}
		ret = append(ret, ID)
	}
	if 0 == len(ret) {
		return ret, errors.New("未找到模块")
	}

	return ret, nil
}

//GetPosIDsByHostID get pos id by hostid
func GetPosIDsByHostID(cc *api.APIResource, posCond interface{}) ([]int, error) {
	result := make([]interface{}, 0)
	var ret []int

	tableName := metadataTable.IdcHostConfig{}
	err := cc.InstCli.GetMutilByCondition(tableName.TableName(), []string{common.BKPosIDField}, posCond, &result, "", 0, 100)
	blog.Infof("GetPosIDsByHostID condition:%v", posCond)
	blog.Infof("result:%v", result)
	if nil != err {
		blog.Error("getPosIDsByHostID error:%", err.Error())
		return ret, errors.New("未找到主机所属模块")
	}
	for _, r := range result {
		item := r.(bson.M)
		ID, getErr := util.GetIntByInterface(item[common.BKPosIDField])
		if nil != getErr || ID == 0 {
			return ret, errors.New("未找到模块")
		}
		ret = append(ret, ID)
	}
	return ret, err
}



//获取业务下的默认模块
func GetIDlePosID(cc *api.APIResource, appID int) (int, error) {
	defaultPosCond := make(map[string]interface{}, 2)
	defaultPosCond[common.BKIdcIDField] = appID
	result := make(map[string]interface{}, 0)
	err := cc.InstCli.GetOneByCondition(posBaseTaleName, []string{common.BKPosIDField}, defaultPosCond, &result)

	if nil != err {
		blog.Error("getDefaultPosIDs error:%s", err.Error())
		return 0, errors.New("未找到模块")
	}

	ID, ok := util.GetIntByInterface(result[common.BKPosIDField])
	if nil != ok {
		return ID, errors.New("未找到模块")
	}

	return ID, nil
}
