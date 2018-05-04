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

package models

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/source_controller/api/metadata"
	dbStorage "configcenter/src/storage"
)

func AddObjClassificationData(tableName string, metaCli dbStorage.DI) error {
	blog.Errorf("add data for  %s table ", tableName)
	rows := getObjClassificationData()
	for _, row := range rows {
		selector :=
			map[string]interface{}{
				common.BKClassificationIDField: row.ClassificationID,
			}
		isExist, err := metaCli.GetCntByCondition(tableName, selector)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tableName, err)
			return err
		}
		if isExist > 0 {
			continue
		}
		id, err := metaCli.GetIncID(tableName)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tableName, err)
			return err
		}
		row.ID = int(id)
		_, err = metaCli.Insert(tableName, row)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tableName, err)
			return err
		}
	}

	blog.Errorf("add data for  %s table  ", tableName)
	return nil
}

func getObjClassificationData() []*metadata.ObjClassification {

	dataRows := []*metadata.ObjClassification{
		&metadata.ObjClassification{ClassificationID: "bk_host_manage", ClassificationName: "主机管理", ClassificationType: "inner", ClassificationIcon: "icon-cc-business"},
		&metadata.ObjClassification{ClassificationID: "bk_biz_topo", ClassificationName: "业务拓扑", ClassificationType: "inner", ClassificationIcon: "icon-cc-square"},
		&metadata.ObjClassification{ClassificationID: "bk_organization", ClassificationName: "组织架构", ClassificationType: "inner", ClassificationIcon: "icon-cc-free-pool"},
		&metadata.ObjClassification{ClassificationID: "bk_network", ClassificationName: "网络", ClassificationType: "inner", ClassificationIcon: "icon-cc-networks"},
		&metadata.ObjClassification{ClassificationID: "bk_middleware", ClassificationName: "中间件", ClassificationType: "inner", ClassificationIcon: "icon-cc-record"},
		&metadata.ObjClassification{ClassificationID: "bk_idc", ClassificationName: "IDC", ClassificationType: "inner", ClassificationIcon: "icon-cc-idc"},
	}

	return dataRows

}
