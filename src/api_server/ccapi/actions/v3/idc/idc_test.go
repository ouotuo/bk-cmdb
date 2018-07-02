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
 
package idc_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"configcenter/src/common/http/httpclient"
	"time"
)

const Address = "http://192.168.148.138:8080/api/v3"

var hclient = httpclient.NewHttpClient()

func init() {
	hclient.SetHeader("Bk_user","runTest")
	hclient.SetHeader("Http_blueking_supplier_id","0")
}

func TestTopoInst(t *testing.T) {

	rsp, rspErr := hclient.GET(Address+"/idc/inst/0/2", nil, nil)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Printf("rsp:%s", rsp)
}

func TestCreateIdc(t *testing.T) {

	param := map[string]interface{}{
		"bk_idc_name":          "app_test",
		"bk_idc_operator":          "123",
	}

	paramsJs, _ := json.Marshal(param)
	rsp, rspErr := hclient.POST(Address+"/idc/0", nil, paramsJs)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Printf("rsp:%s", rsp)
}
func TestCreateRack(t *testing.T) {

	param := map[string]interface{}{
		"bk_rack_name":          "rack_test",
		"bk_supplier_account":          "0",
		"bk_rack_open_time":time.Now().Format("2006-01-02 15:04:05"),
		"bk_rack_use_time":time.Now().Format("2006-01-02 15:04:05"),
	}

	paramsJs, _ := json.Marshal(param)
	rsp, rspErr := hclient.POST(Address+"/rack/1", nil, paramsJs)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Printf("rsp:%s", rsp)
}
func TestCreatePos(t *testing.T) {

	param := map[string]interface{}{
		"bk_pos_name":          "pos_test",
		"bk_supplier_account":          "0",
		"bk_pos_open_time":time.Now().Format("2006-01-02 15:04:05"),
		"bk_pos_use_time":time.Now().Format("2006-01-02 15:04:05"),
	}

	paramsJs, _ := json.Marshal(param)
	rsp, rspErr := hclient.POST(Address+"/pos/1/1", nil, paramsJs)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Printf("rsp:%s", rsp)
}


func TestGetRack(t *testing.T) {

	param := map[string]interface{}{
		"condition":    map[string]interface{}{
			"bk_rack_id":"9",
		},
	}

	paramsJs, _ := json.Marshal(param)
	rsp, rspErr := hclient.POST(Address+"v3/rack/search/0/8", nil, paramsJs)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Printf("rsp:%s", rsp)
}