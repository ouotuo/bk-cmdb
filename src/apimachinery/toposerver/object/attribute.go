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

package object

import (
	"context"
	"fmt"

	"configcenter/src/apimachinery/util"
	"configcenter/src/common/core/cc/api"
	obj "configcenter/src/source_controller/api/object"
)

func (t *object) CreateObjectAtt(ctx context.Context, h util.Headers, obj *obj.ObjAttDes) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/objectattr"

	err = t.client.Post().
		WithContext(ctx).
		Body(obj).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *object) SelectObjectAttWithParams(ctx context.Context, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/objectattr/search"

	err = t.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *object) UpdateObjectAtt(ctx context.Context, objID string, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/objectattr/%s", objID)

	err = t.client.Put().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *object) DeleteObjectAtt(ctx context.Context, objID string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/objectattr/%s", objID)

	err = t.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}
