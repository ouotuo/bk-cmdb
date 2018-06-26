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

package controllers

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/core/cc/wactions"
	"fmt"

	//"configcenter/src/web_server/application/middleware/user"
	"github.com/gin-gonic/contrib/cache"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/ldap.v2"
	"strings"
	"time"
)

var UserCache *cache.InMemoryStore

const BkAccountUrl = "site.bk_account_url"

func init() {
	UserCache = cache.NewInMemoryStore(time.Hour)
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectGet, "/user/list", nil, GetUserList})
	wactions.RegisterNewAction(wactions.Action{common.HTTPUpdate, "/user/language/:language", nil, UpdateUserLanguage})

}

type userResult struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Code    string      `json:"code"`
	Result  bool        `json:"result"`
}

func GetUserByLDAP(ladpIp string, baseDN string, ldapPort string, bindDN string, bindPasswd string) []map[string]interface{} {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", ladpIp, ldapPort))
	if err != nil {
		return []map[string]interface{}{}
	}
	res := []map[string]interface{}{}
	err = l.Bind(bindDN, bindPasswd)
	if err != nil {
		return []map[string]interface{}{}
	}
	searchRequest := ldap.NewSearchRequest(
		baseDN, // The base dn to search
		ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=inetOrgPerson))",                                            // The filter to apply
		[]string{"dn", "cn", "sn", "mail", "mobile", "uniqueMember", "objectClass"}, // A list attributes to retrieve
		nil,
	)
	sr, err := l.Search(searchRequest)

	blog.Info("len Entries %s", len(sr.Entries))
	for _, entry := range sr.Entries {
		blog.Info(entry.DN)
		cellData := make(map[string]interface{})
		username := entry.GetAttributeValue("cn") + "." + entry.GetAttributeValue("sn")
		nikiname := strings.Split(entry.GetAttributeValue("mail"), "@")[0]
		cellData["chinese_name"] = username
		cellData["english_name"] = nikiname
		res = append(res, cellData)
	}
	defer l.Close()
	return res

}

//GetUserList get user list
// func GetUserList(c *gin.Context) {
// 	session := sessions.Default(c)
// 	skiplogin := session.Get("skiplogin")
// 	blog.Info("skiplogin %s", skiplogin)
// 	skiplogins, ok := skiplogin.(string)
// 	if ok && "1" == skiplogins {
// 		admindata := make([]interface{}, 0)
// 		admincell := make(map[string]interface{})
// 		admincell["chinese_name"] = "admin"
// 		admincell["english_name"] = "admin"
// 		admindata = append(admindata, admincell)
// 		c.JSON(200, gin.H{
// 			"result":        true,
// 			"bk_error_msg":  "get user list ok",
// 			"bk_error_code": "00",
// 			"data":          admindata,
// 		})
// 		return
// 	}
// }

type userDataResult struct {
	Message string      `json:"bk_error_msg"`
	Data    interface{} `json:"data"`
	Code    string      `json:"bk_error_code"`
	Result  bool        `json:"result"`
}

var getUserFailData = userDataResult{
	Result:  false,
	Message: "get user list false",
	Code:    "",
	Data:    nil,
}

// GetUserList get user list
func GetUserList(c *gin.Context) {
	a := api.NewAPIResource()
	config, _ := a.ParseConfig()
	ldapIp := config["ldap.ldap_ip"]
	ldapPort := config["ldap.ldap_port"]
	baseDN := config["ldap.ldap_baseDN"]
	bindDN := config["ldap.ldap_bindDN"]
	bindPasswd := config["ldap.ldap_bind_passwd"]

	res := []map[string]interface{}{}
	err := UserCache.Get("userList", &res)
	if err != nil {
		res = GetUserByLDAP(ldapIp, baseDN, ldapPort, bindDN, bindPasswd)
		UserCache.Set("userList", res, 3600*time.Second)
	}
	c.JSON(200, gin.H{
		"result":        true,
		"bk_error_msg":  "get user list ok",
		"bk_error_code": "00",
		"data":          res,
	})
	// accountURL := config[BkAccountUrl]
	// user := user.NewUser()
	// code, data := user.GetUserList(c, accountURL)
	// c.JSON(code, data)
	return
}

func UpdateUserLanguage(c *gin.Context) {
	session := sessions.Default(c)
	language := c.Param("language")

	session.Set("language", language)
	err := session.Save()

	if nil != err {
		blog.Errorf("user update language error:%s", err.Error())
		c.JSON(200, userDataResult{
			Result:  false,
			Message: "user update language error",
			Code:    fmt.Sprintf("%d", common.CCErrCommHTTPDoRequestFailed),
			Data:    nil,
		})
		return
	}

	c.SetCookie("blueking_language", language, 0, "/", "", false, true)

	c.JSON(200, userDataResult{
		Result:  true,
		Message: "",
		Code:    "00",
		Data:    nil,
	})
	return
}
