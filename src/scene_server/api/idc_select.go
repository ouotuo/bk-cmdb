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
 
package api

import (
	"configcenter/src/common"
	"fmt"
)

func (cli *Client) ReForwardSelectMetaIdcModel(callfunc func(url, method string) (string, error), ownerid, clsid, objid string) func() (string, error) {

	return func() (string, error) {
		return callfunc(fmt.Sprintf("%s/idc/v1/model/%s/%s/%s", cli.address, ownerid, clsid, objid), common.HTTPSelectGet)
	}
}

func (cli *Client) ReForwardSelectMetaIdcInstChild(callfunc func(url, method string) (string, error), ownerid, objid, appid, instid string) func() (string, error) {

	return func() (string, error) {
		return callfunc(fmt.Sprintf("%s/idc/v1/inst/child/%s/%s/%s/%s", cli.address, ownerid, objid, appid, instid), common.HTTPSelectGet)
	}
}

func (cli *Client) ReForwardSelectMetaIdcInst(callfunc func(url, method string) (string, error), ownerid, appid string) func() (string, error) {

	return func() (string, error) {
		return callfunc(fmt.Sprintf("%s/idc/v1/inst/%s/%s", cli.address, ownerid, appid), common.HTTPSelectGet)
	}
}

func (cli *Client) ReForwardSelectMetaIdc(callfunc func(url, method string) (string, error), ownerid string) func() (string, error) {

	return func() (string, error) {
		return callfunc(fmt.Sprintf("%s/idc/v1/model/%s", cli.address, ownerid), common.HTTPSelectGet)
	}
}
func (cli *Client) ReForwardSelectMetaIdc1(callfunc func(url, method string) (string, error), ownerid, appid string) func() (string, error) {

	return func() (string, error) {
		return callfunc(fmt.Sprintf("%s/idc/v1/idc/search/%s/%s", cli.address, ownerid, appid), common.HTTPSelectPost)
	}
}
func (cli *Client) ReForwardSelectMetaIdcTopoInst(callfunc func(url, method string) (string, error), ownerid, appid string) func() (string, error) {

	return func() (string, error) {
		return callfunc(fmt.Sprintf("%s/idc/v1/inst/%s/%s", cli.address, ownerid, appid), common.HTTPSelectGet)
	}
}

func (cli *Client) ReForwardSelectMetaIdcTopoByClsID(callfunc func(url, method string) (string, error), ownerid, clsid, objid string) func() (string, error) {

	return func() (string, error) {
		return callfunc(fmt.Sprintf("%s/idc/v1/model/%s/%s/%s", cli.address, ownerid, clsid, objid), common.HTTPSelectGet)
	}
}
func (cli *Client) ReForwardSelectMetaIdcTopoInstChild(callfunc func(url, method string) (string, error), ownerid, objid, appid, instid string) func() (string, error) {

	return func() (string, error) {
		return callfunc(fmt.Sprintf("%s/idc/v1/inst/child/%s/%s/%s/%s", cli.address, ownerid, objid, appid, instid), common.HTTPSelectGet)
	}
}