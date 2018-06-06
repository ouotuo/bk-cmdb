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

package hostserver

import (
    "context"
    "fmt"
    
    "configcenter/src/apimachinery/util"
    "configcenter/src/common/core/cc/api"
    "configcenter/src/common/paraparse"
    "configcenter/src/scene_server/host_server/host_service/actions/hosts"
    "configcenter/src/source_controller/common/commondata"
)

func (hs *hostServer) DeleteHostBatch(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/host/batch"

	err = hs.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetHostDetailByID(ctx context.Context, hostID string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/hosts/%s/%s", h.OwnerID, hostID)

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostSnapInfo(ctx context.Context, hostID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/host/snapshot/%s", hostID)

	err = hs.client.Get().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AddHost(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/hosts/addhost"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AddHostFromAgent(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/host/add/agent"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetHostFavourites(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "hosts/favorites/search"

	err = hs.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AddHostFavourite(ctx context.Context, h util.Headers, dat *hosts.FavouriteParms) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "hosts/favorites"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) UpdateHostFavouriteByID(ctx context.Context, id string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("hosts/favorites/%s", id)

	err = hs.client.Put().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) DeleteHostFavouriteByID(ctx context.Context, id string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("hosts/favorites/%s", id)

	err = hs.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) IncrHostFavouritesCount(ctx context.Context, id string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/hosts/favorites/%s/incr", id)

	err = hs.client.Put().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AddHistory(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/history"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetHistorys(ctx context.Context, start string, limit string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/history/%s/%s", start, limit)

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AddHostMutiltAppModuleRelation(ctx context.Context, h util.Headers, dat *hosts.CloudHostModuleParams) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/hosts/modules/biz/mutiple"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostModuleRelation(ctx context.Context, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/hosts/modules"

	err = hs.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) MoveHost2EmptyModule(ctx context.Context, h util.Headers, dat *hosts.DefaultModuleHostConfigParams) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/hosts/emptymodule"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) MoveHost2FaultModule(ctx context.Context, h util.Headers, dat *hosts.DefaultModuleHostConfigParams) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/hosts/faultmodule"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) MoveHostToResourcePool(ctx context.Context, h util.Headers, dat *hosts.DefaultModuleHostConfigParams) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/hosts/resource"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AssignHostToApp(ctx context.Context, h util.Headers, dat *hosts.DefaultModuleHostConfigParams) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/hosts/assgin"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AssignHostToAppModule(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/host/add/module"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) SaveUserCustom(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/usercustom"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetUserCustom(ctx context.Context, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/usercustom/user/search"

	err = hs.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetDefaultCustom(ctx context.Context, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/usercustom/default/search"

	err = hs.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetAgentStatus(ctx context.Context, appID string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("getAgentStatus/%s", appID)

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) UpdateHost(ctx context.Context, appID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/openapi/host/%s", appID)

	err = hs.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) UpdateHostByAppID(ctx context.Context, appID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/host/updateHostByAppID/%s", appID)

	err = hs.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetHostListByAppidAndField(ctx context.Context, appID string, field string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/host/getHostListByAppidAndField/%s/%s", appID, field)

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostSearchByIP(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/gethostlistbyip"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostSearchByModuleID(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/getmodulehostlist"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostSearchBySetID(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/getsethostlist"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostSearchByAppID(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/getapphostlist"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostSearchByProperty(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/gethostsbyproperty"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetIPAndProxyByCompany(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/getIPAndProxyByCompany"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) UpdateCustomProperty(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/openapi/updatecustomproperty"

	err = hs.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) CloneHostProperty(ctx context.Context, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "openapi/host/clonehostproperty"

	err = hs.client.Put().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetHostAppByCompanyId(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/openapi/host/getHostAppByCompanyId"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) DelHostInApp(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/openapi/host/delhostinapp"

	err = hs.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetGitServerIp(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/openapi/host/getGitServerIp"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetPlat(ctx context.Context, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/plat"

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) CreatePlat(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/plat"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) DelPlat(ctx context.Context, cloudID string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/plat/%s", cloudID)

	err = hs.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostSearch(ctx context.Context, h util.Headers, dat *params.HostCommonSearch) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/search"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostSearchWithAsstDetail(ctx context.Context, h util.Headers, dat *params.HostCommonSearch) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/search/asstdetail"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) UpdateHostBatch(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/host/batch"

	err = hs.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AddUserCustomQuery(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/userapi"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) UpdateUserCustomQuery(ctx context.Context, businessID string, id string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/userapi/%s/%s", businessID, id)

	err = hs.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) DeleteUserCustomQuery(ctx context.Context, businessID string, id string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/userapi/%s/%s", businessID, id)

	err = hs.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetUserCustomQuery(ctx context.Context, businessID string, h util.Headers, dat *commondata.ObjQueryInput) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/userapi/search/%s", businessID)

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetUserCustomQueryDetail(ctx context.Context, businessID string, id string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/userapi/detail/%s/%s", businessID, id)

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetUserCustomQueryResult(ctx context.Context, businessID, id, start, limit string, h util.Headers) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/userapi/data/%s/%s/%s/%s", businessID, id, start, limit)

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}
