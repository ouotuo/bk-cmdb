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

package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"net/http"
)

// Property object fields
type Property struct {
	ID            string
	Name          string
	PropertyType  string
	Option        interface{}
	IsPre         bool
	IsRequire     bool
	Group         string
	Index         int
	ExcelColIndex int
	NotObjPropery bool //Not an attribute of the object, indicating that the field to be exported is needed for export,
}

// PropertyGroup property group
type PropertyGroup struct {
	Name  string
	Index int
}

// GetObjFieldIDs get object fields
func GetObjFieldIDs(objID, url string, header http.Header) (map[string]Property, error) {

	fields, err := getObjFieldIDs(objID, url, header)
	if nil != err {
		return nil, err
	}
	groups, err := getObjectGroup(objID, url, header)
	if nil != err {
		return nil, err
	}

	ret := make(map[string]Property)
	index := 0

	for _, group := range groups {
		for _, field := range fields {
			if field.Group == group.Name {
				field.ExcelColIndex = index
				ret[field.ID] = field
				index++
			}
		}
	}

	return ret, nil
}

func getObjectGroup(objID, url string, header http.Header) ([]PropertyGroup, error) {
	///api/v3/objectatt/group/property/owner/0/object/host
	url = fmt.Sprintf("%s/api/%s/objectatt/group/property/owner/%s/object/%s", url, webCommon.API_VERSION, util.GetActionOnwerIDByHTTPHeader(header), objID)
	conds := common.KvMap{common.BKObjIDField: objID, common.BKOwnerIDField: common.BKDefaultOwnerID, "page": common.KvMap{"start": 0, "limit": common.BKNoLimit, "sort": common.BKPropertyGroupIndexField}}
	result, err := httpRequest(url, conds, header)
	if nil != err {
		return nil, err
	}
	blog.Info("get %s fields group  url:%s", objID, url)
	blog.Info("get %s fields group return:%s", objID, result)
	js, err := simplejson.NewJson([]byte(result))
	if nil != err {
		blog.Info("get %s fields group  url:%s return:%s", objID, url, result)
		return nil, err
	}
	fields, _ := js.Get("data").Array()
	ret := []PropertyGroup{}
	for _, field := range fields {
		mapField, _ := field.(map[string]interface{})
		propertyGroup := PropertyGroup{}
		propertyGroup.Index, _ = util.GetIntByInterface(mapField[common.BKPropertyGroupIndexField])
		propertyGroup.Name, _ = mapField[common.BKPropertyGroupNameField].(string)
		ret = append(ret, propertyGroup)
	}
	return ret, nil

}

func getObjFieldIDs(objID, url string, header http.Header) ([]Property, error) {
	url = fmt.Sprintf("%s/api/%s/object/attr/search", url, webCommon.API_VERSION)
	conds := common.KvMap{
		common.BKObjIDField:   objID,
		common.BKOwnerIDField: common.BKDefaultOwnerID,
		"page": common.KvMap{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  fmt.Sprintf("-%s,bk_property_index", common.BKIsRequiredField),
		},
	}
	result, err := httpRequest(url, conds, header)
	if nil != err {
		return nil, err
	}
	blog.Info("get %s fields  url:%s", objID, url)
	blog.Info("get %s fields return:%s", objID, result)
	js, err := simplejson.NewJson([]byte(result))
	if nil != err {
		blog.Info("get %s fields  url:%s return:%s", objID, url, result)
		return nil, err
	}
	fields, _ := js.Get("data").Array()
	ret := []Property{}

	for _, field := range fields {
		mapField, _ := field.(map[string]interface{})

		fieldName, _ := mapField[common.BKPropertyNameField].(string)
		fieldID, _ := mapField[common.BKPropertyIDField].(string)
		fieldType, _ := mapField[common.BKPropertyTypeField].(string)
		fieldIsRequire, _ := mapField[common.BKIsRequiredField].(bool)
		fieldIsOption, _ := mapField[common.BKOptionField]
		fieldIsPre, _ := mapField[common.BKIsPre].(bool)
		fieldGroup, _ := mapField["bk_property_group_name"].(string)
		fieldIndex, _ := util.GetIntByInterface(mapField["bk_property_index"])

		ret = append(ret, Property{
			ID:           fieldID,
			Name:         fieldName,
			PropertyType: fieldType,
			IsRequire:    fieldIsRequire,
			IsPre:        fieldIsPre,
			Option:       fieldIsOption,
			Group:        fieldGroup,
			Index:        fieldIndex,
		})
	}

	return ret, nil
}

// getPropertyTypeAliasName  return propertyType name, whether to export,
func getPropertyTypeAliasName(propertyType string, defLang lang.DefaultCCLanguageIf) (string, bool) {
	var skip bool
	name := defLang.Language("field_type_" + propertyType)
	switch propertyType {
	case common.FieldTypeSingleChar:
	case common.FieldTypeLongChar:
	case common.FieldTypeInt:
	case common.FieldTypeEnum:
	case common.FieldTypeDate:
	case common.FieldTypeTime:
	case common.FieldTypeUser:
	case common.FieldTypeSingleAsst:
	case common.FieldTypeMultiAsst:
	case common.FieldTypeBool:
	case common.FieldTypeTimeZone:

	}
	if "" == name {
		name = propertyType
	}
	return name, skip
}
