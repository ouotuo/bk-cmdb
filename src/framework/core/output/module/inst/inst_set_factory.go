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
	"configcenter/src/framework/core/types"
)

func createSet(target model.Model) (Inst, error) {
	return &set{target: target, datas: types.MapStr{}}, nil
}

// findSetsLikeName find all insts by inst name
func findSetsLikeName(target model.Model, setName string) (Iterator, error) {
	cond := common.CreateCondition().Field(SetName).Like(setName)
	return newIteratorInstSet(target, cond)
}

// findSetsByCondition find all insts by condition
func findSetsByCondition(target model.Model, cond common.Condition) (Iterator, error) {
	return newIteratorInstSet(target, cond)
}
