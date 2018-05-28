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
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/core/cc/wactions"
	"fmt"
	"configcenter/src/common/types"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/ldap.v2"
	"errors"
	"crypto/md5"
	"encoding/hex"
	"time"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/common/blog"
	"github.com/bitly/go-simplejson"
	"strings"
	"encoding/json"
)

//LogOutUser log out user
func LogOutUser(c *gin.Context) {
	fmt.Sprintf("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
	a := api.NewAPIResource()
	config, _ := a.ParseConfig()
	site := config["site.domain_url"]
	loginURL := config["site.bk_login_url"]
	appCode := config["site.app_code"]
	loginPage := fmt.Sprintf(loginURL, appCode, site)
	session := sessions.Default(c)
	session.Clear()
	c.Redirect(302, loginPage)
}
func UserAuthentication(username string,password string,c *gin.Context) error {


	cc := api.NewAPIResource()

	objID := c.Param(common.BKObjIDField)

	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)

	a := api.NewAPIResource()
	config, _ := a.ParseConfig()
	ldapIp := config["ldap.ldap_ip"]
	ldapPort := config["ldap.ldap_port"]
	baseDN := config["ldap.ldap_baseDN"]
	bindDN := config["ldap.ldap_bindDN"]
	bindPasswd := config["ldap.ldap_bind_passwd"]
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", ldapIp, ldapPort))
	if err != nil {
		return err
	}
	defer l.Close()

	err = l.Bind(bindDN, bindPasswd)
	if err != nil {
		fmt.Sprintf("(%s=%s)",bindDN, bindPasswd)
		return errors.New("User does not exist ")
	}
	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)","cn", username),
		[]string{"dn", "cn", "sn", "mail", "mobile", "userPassword"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		return errors.New(err.Error()+fmt.Sprintf("(%s=%s)",bindDN, bindPasswd))
	}
	if len(sr.Entries) == 0 {
		return errors.New("User does not exist ")
	}
	for _, entriy := range sr.Entries {
		userdn := entriy.DN
		// Bind as the user to verify their password
		err = l.Bind(userdn, password)
		if err == nil {
			h := md5.New()
			h.Write([]byte(username+bindPasswd+time.Now().String())) // 需要加密的字符串为 123456
			cipherStr := h.Sum(nil)
			bk_token:= hex.EncodeToString(cipherStr)
			session := sessions.Default(c)
			nikiname := strings.Split(entriy.GetAttributeValue("mail"), "@")[0]
			session.Set("userName", nikiname)
			session.Set("chName", username)
			session.Set("phone", entriy.GetAttributeValue("mobile"))
			session.Set("email", entriy.GetAttributeValue("mail"))
			/// /api/{version}/topo/privilege/group/{bk_supplier_account}/search
			c.Request.Header.Set("Content-Type", "application/json;charset=UTF-8")
			c.Request.Header.Set("BK_USER", "ads")
			c.Request.Header.Set("HTTP_BLUEKING_SUPPLIER_ID", "0")
			url := fmt.Sprintf("%s/api/%s/topo/privilege/group/%s/search/", apiSite, webCommon.API_VERSION, "0")
			conds := common.KvMap{"user_list":nikiname}
			if config["session.admin_user"]==nikiname {
				session.Set("role", "1")
			}else {
				result, err := httpRequest(url, conds,  c.Request.Header)
				if nil != err {
					blog.Info("httpRequest error  %s  url:%s ", err.Error(), url)
					session.Set("role", "0")
				}
				blog.Info("get %s fields group  url:%s", objID, url)
				blog.Info("get %s fields group return:%s", objID, result)
				js, err := simplejson.NewJson([]byte(result))
				if nil != err {
					blog.Info("get %s fields group  url:%s return:%s", objID, url, result)
					session.Set("role", "0")
				}
				fields, _ := js.Get("data").Array()
				role:=""
				for _, field := range fields {
					ss,_:=json.Marshal(field)
					s, _ := simplejson.NewJson(ss)
					groupId,err:=s.Get("group_id").String()
					if err==nil{
						role=role+groupId+","
					}
					groupName,err:=s.Get("group_name").String()
					if err==nil && groupName=="admin"{
						role="1"
						break
					}

				}
				blog.Info("get %s fields group return:%s", fields)
				session.Set("role", role)
			}

			session.Set("bk_token", bk_token)
			session.Set("owner_uin", common.BKDefaultOwnerID)
			session.Set("skiplogin", "0")
			session.Save()
			blog.Info("11111111valid user login session token %s", bk_token)
			c.SetCookie("bk_token",bk_token,3600,"","",false,false)
			return nil
		}else{
			err1 := l.Bind(bindDN, bindPasswd)
			if err1 != nil {
				fmt.Sprintf("(%s=%s)(%s)",bindDN, bindPasswd,baseDN)
				return errors.New("User does not exist ")
			}
		}
	}
	return  errors.New(err.Error()+fmt.Sprintf("(%s=%s)(%s)",username, password))
}


func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	err:=UserAuthentication(username,password,c)
	if  err!=nil{
		c.JSON(401, gin.H{
			"status": err.Error(),
		})
	}else{
		c.JSON(200, gin.H{
			"status": "success",
		})
	}
}
func init() {
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectGet, "/logout", nil, LogOutUser})
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectGet, "/login", nil, Login})
}
