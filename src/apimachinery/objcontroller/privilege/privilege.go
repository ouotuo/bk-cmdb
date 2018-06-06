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

package privilege

import (
	"context"
	"fmt"

	"configcenter/src/apimachinery/util"
	"configcenter/src/common/core/cc/api"
)

func (t *privilege) CreateUserGroup(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/privilege/group/%s", h.OwnerID)

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *privilege) UpdateUserGroup(ctx context.Context, groupID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/privilege/group/%s/%s", h.OwnerID, groupID)

	err = t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *privilege) DeleteUserGroup(ctx context.Context, groupID string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/privilege/group/%s/%s", h.OwnerID, groupID)

	err = t.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *privilege) SearchUserGroup(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/privilege/group/%s/search", h.OwnerID)

	err = t.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *privilege) CreateUserGroupPrivi(ctx context.Context, groupID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/privilege/group/detail/%s/%s", h.OwnerID, groupID)

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *privilege) UpdateUserGroupPrivi(ctx context.Context, groupID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/privilege/group/detail/%s/%s", h.OwnerID, groupID)

	err = t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *privilege) GetUserGroupPrivi(ctx context.Context, groupID string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/privilege/group/detail/%s/%s", h.OwnerID, groupID)

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *privilege) CreateRolePri(ctx context.Context, objID string, propertyID string, h util.Headers, role []string) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/role/%s/%s/%s", h.OwnerID, objID, propertyID)

	err = t.client.Post().
		WithContext(ctx).
		Body(role).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *privilege) GetRolePri(ctx context.Context, objID string, propertyID string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/role/%s/%s/%s", h.OwnerID, objID, propertyID)

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *privilege) UpdateRolePri(ctx context.Context, objID string, propertyID string, h util.Headers, role []string) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/role/%s/%s/%s", h.OwnerID, objID, propertyID)

	err = t.client.Put().
		WithContext(ctx).
		Body(role).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *privilege) GetSystemFlag(ctx context.Context, flag string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/system/%s/%s", flag, h.OwnerID)

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}
