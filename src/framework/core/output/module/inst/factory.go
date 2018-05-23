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
 
package inst

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/model"
	mtypes "configcenter/src/framework/core/output/module/types"
	"configcenter/src/framework/core/types"
)

// CreateInst creat a new inst for the model
func CreateInst(target model.Model) (Inst, error) {

	switch target.GetID() {
	case mtypes.BKInnerObjIDBusiness:
		return createBusiness(target)
	case mtypes.BKInnerObjIDHost:
		return createHost(target)
	case mtypes.BKInnerObjIDModule:
		return createModule(target)
	case mtypes.BKInnerObjIDPlat:
		return createPlat(target)
	case mtypes.BKInnerObjIDProc:
		return createProc(target)
	case mtypes.BKInnerObjIDSet:
		return createSet(target)
	default:
		return &inst{target: target, datas: types.MapStr{}}, nil

	}

}

// FindInstsLikeName find all insts by inst name
func FindInstsLikeName(target model.Model, instName string) (Iterator, error) {

	switch target.GetID() {
	case mtypes.BKInnerObjIDBusiness:
		return findBusinessLikeName(target, instName)
	case mtypes.BKInnerObjIDHost:
		return findHostsLikeName(target, instName)
	case mtypes.BKInnerObjIDModule:
		return findModulesLikeName(target, instName)
	case mtypes.BKInnerObjIDPlat:
		return findPlatsLikeName(target, instName)
	case mtypes.BKInnerObjIDProc:
		return findProcsLikeName(target, instName)
	case mtypes.BKInnerObjIDSet:
		return findSetsLikeName(target, instName)
	default:
		cond := common.CreateCondition().Field(InstName).Like(instName)
		return newIteratorInst(target, cond)

	}

}

// FindInstsByCondition find all insts by condition
func FindInstsByCondition(target model.Model, cond common.Condition) (Iterator, error) {
	// TODO:按照条件读取所有实例
	switch target.GetID() {
	case mtypes.BKInnerObjIDBusiness:
		return findBusinessByCondition(target, cond)
	case mtypes.BKInnerObjIDHost:
		return findHostsByCondition(target, cond)
	case mtypes.BKInnerObjIDModule:
		return findModulesByCondition(target, cond)
	case mtypes.BKInnerObjIDPlat:
		return findPlatsByCondition(target, cond)
	case mtypes.BKInnerObjIDProc:
		return findProcsByCondition(target, cond)
	case mtypes.BKInnerObjIDSet:
		return findSetsByCondition(target, cond)
	default:
		return newIteratorInst(target, cond)

	}

}
