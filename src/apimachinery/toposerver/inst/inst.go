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
	"context"
	"fmt"

	"configcenter/src/apimachinery/util"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/paraparse"
	"configcenter/src/scene_server/topo_server/topo_service/actions/inst"
)

// TODO: config this body data struct.
func (t *instanceClient) CreateInst(ctx context.Context, objID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/inst/%s/%s", h.OwnerID, objID)

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) DeleteInst(ctx context.Context, objID string, instID string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/inst/%s/%s/%s", h.OwnerID, objID, instID)

	err = t.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) UpdateInst(ctx context.Context, objID string, instID string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/inst/%s/%s/%s", h.OwnerID, objID, instID)

	err = t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) SelectInsts(ctx context.Context, objID string, h util.Headers, s *params.SearchParams) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/inst/search/%s/%s", h.OwnerID, objID)

	err = t.client.Post().
		WithContext(ctx).
		Body(s).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) SelectInstsAndAsstDetail(ctx context.Context, objID string, h util.Headers, s *params.SearchParams) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/inst/search/owner/%s/object/%s/detail", h.OwnerID, objID)

	err = t.client.Post().
		WithContext(ctx).
		Body(s).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) InstSearch(ctx context.Context, objID string, h util.Headers, s *params.SearchParams) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/inst/search/owner/%s/object/%s", h.OwnerID, objID)

	err = t.client.Post().
		WithContext(ctx).
		Body(s).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) SelectInstsByAssociation(ctx context.Context, objID string, h util.Headers, p *inst.AssociationParams) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/inst/association/search/owner/%s/object/%s", h.OwnerID, objID)

	err = t.client.Post().
		WithContext(ctx).
		Body(p).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) SelectInst(ctx context.Context, objID string, instID string, h util.Headers, p *params.SearchParams) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/inst/search/owner/%s/%s/%s", h.OwnerID, objID, instID)

	err = t.client.Post().
		WithContext(ctx).
		Body(p).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) SelectTopo(ctx context.Context, objID string, instID string, h util.Headers, p *params.SearchParams) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/inst/search/topo/owner/%s/object/%s/inst/%s", h.OwnerID, objID, instID)

	err = t.client.Post().
		WithContext(ctx).
		Body(p).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) SelectAssociationTopo(ctx context.Context, objID string, instID string, h util.Headers, p *params.SearchParams) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/inst/association/topo/search/owner/%sobject/%s/inst/%s", h.OwnerID, objID, instID)

	err = t.client.Post().
		WithContext(ctx).
		Body(p).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}
