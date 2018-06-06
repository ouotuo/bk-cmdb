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

package history

import (
    "context"
    "fmt"
    
    "configcenter/src/apimachinery/util"
    "configcenter/src/common/core/cc/api"
)

func(hi *history) AddHistory(ctx context.Context, user string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/history/%s", user)

        err = hi.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func(hi *history) GetHistorys(ctx context.Context, user string, start string, limit string, h util.Headers) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/history/%s/%s/%s", user, start, limit)

        err = hi.client.Get().
        WithContext(ctx).
        Body(nil).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}