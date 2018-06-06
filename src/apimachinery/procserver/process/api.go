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

package process

import (
    "context"
    
    "configcenter/src/apimachinery/rest"
    "configcenter/src/apimachinery/util"
    "configcenter/src/common/core/cc/api"
)

type ProcessClientInterface interface {
    GetProcessDetailByID(ctx context.Context, appID string, procID string, h util.Headers) (resp *api.BKAPIRsp, err error)
    GetProcessBindModule(ctx context.Context, businessID string, procID string, h util.Headers) (resp *api.BKAPIRsp, err error)
    BindModuleProcess(ctx context.Context, businessID string, procID string, moduleName string, h util.Headers) (resp *api.BKAPIRsp, err error)
    DeleteModuleProcessBind(ctx context.Context, businessID string, procID string, moduleName string, h util.Headers) (resp *api.BKAPIRsp, err error)
    CreateProcess(ctx context.Context, businessID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    DeleteProcess(ctx context.Context, businessID string, procID string, h util.Headers) (resp *api.BKAPIRsp, err error)
    SearchProcess(ctx context.Context, businessID string, h util.Headers) (resp *api.BKAPIRsp, err error)
    UpdateProcess(ctx context.Context, businessID string, procID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)

}

func NewProcessClientInterface(client rest.ClientInterface) ProcessClientInterface {
    return &process{client:client}
}

type process struct {
    client rest.ClientInterface
}
