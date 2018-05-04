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

package data

import (
	"configcenter/src/common"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/validator"
	"configcenter/src/source_controller/api/metadata"
)

// default group
var (
	groupBaseInfo = mCommon.BaseInfo
)

// Distribution init revision
var Distribution = "community" // could be community or enterprise

/*
	&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "", PropertyName: "", IsRequired: , IsOnly: , PropertyGroup: , PropertyType: , Option: ""},
*/

// AppRow app structure
func AppRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDApp

	groupAppRole := mCommon.AppRole

	lifeCycleOption := []validator.EnumVal{{ID: "1", Name: "测试中", Type: "text"}, {ID: "2", Name: "已上线", Type: "text", IsDefault: true}, {ID: "3", Name: "停运", Type: "text"}}
	languageOption := []validator.EnumVal{{ID: "1", Name: "中文", Type: "text"}, {ID: "2", Name: "English", Type: "text"}}
	dataRows := []*metadata.ObjectAttDes{
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_biz_name", PropertyName: "业务名", IsRequired: true, IsOnly: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "life_cycle", PropertyName: "生命周期", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: lifeCycleOption},

		//role
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKMaintainersField, PropertyName: "运维人员", IsRequired: true, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FieldTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKProductPMField, PropertyName: "产品人员", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FieldTypeUser, Option: ""},

		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKTesterField, PropertyName: "测试人员", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FieldTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_biz_developer", PropertyName: "开发人员", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FieldTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKOperatorField, PropertyName: "操作人员", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FieldTypeUser, Option: ""},
	}

	if Distribution == common.RevisionEnterprise {
		dataRows = append(dataRows,
			&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "time_zone", PropertyName: "时区", IsRequired: true, IsOnly: false, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeTimeZone, Option: "", IsReadOnly: true},
			&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "language", PropertyName: "语言", IsRequired: true, IsOnly: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: languageOption, IsReadOnly: true},
		)
	}

	return dataRows

}

// SetRow set structure
func SetRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDSet
	serviceStatusOption := []validator.EnumVal{{ID: "1", Name: "开放", Type: "text", IsDefault: true}, {ID: "2", Name: "关闭", Type: "text"}}

	dataRows := []*metadata.ObjectAttDes{
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKAppIDField, PropertyName: "业务ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: common.KvMap{}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_set_name", PropertyName: "集群名字", IsRequired: true, IsOnly: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_set_desc", PropertyName: "集群描述", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_set_env", PropertyName: "环境类型", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "测试", Type: "text"}, {ID: "2", Name: "体验", Type: "text"}, {ID: "3", Name: "正式", Type: "text", IsDefault: true}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_service_status", PropertyName: "服务状态", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: serviceStatusOption},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "description", PropertyName: "备注", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_capacity", PropertyName: "设计容量", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: common.KvMap{"min": "1", "max": "999999999"}},

		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKChildStr, PropertyName: "", IsRequired: false, IsOnly: false, PropertyType: "", Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKInstParentStr, PropertyName: "", IsSystem: true, IsRequired: true, IsOnly: true, PropertyType: common.FieldTypeInt, Option: ""},
	}
	return dataRows
}

// ModuleRow module structure
func ModuleRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDModule
	moduleTypeOption := []validator.EnumVal{{ID: "1", Name: "普通", Type: "text", IsDefault: true}, {ID: "2", Name: "数据库", Type: "text"}}

	dataRows := []*metadata.ObjectAttDes{
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKAppIDField, PropertyName: "业务ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: common.KvMap{}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKSetIDField, PropertyName: "集群ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: common.KvMap{}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKModuleNameField, PropertyName: "模块名", IsRequired: true, IsOnly: true, Editable: true, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKChildStr, PropertyName: "", IsRequired: false, IsOnly: false, PropertyType: "", Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_module_type", PropertyName: "模块类型", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: moduleTypeOption},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "operator", PropertyName: "主要维护人", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_bak_operator", PropertyName: "备份维护人", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeUser, Option: ""},
	}
	return dataRows
}

// PlatRow plat structure
func PlatRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDPlat
	dataRows := []*metadata.ObjectAttDes{
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKCloudNameField, PropertyName: "云区域", IsRequired: true, IsOnly: true, IsPre: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKOwnerIDField, PropertyName: "供应商", IsRequired: true, IsOnly: true, IsPre: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
	}
	return dataRows
}

// HostRow host structure
func HostRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDHost
	groupAgent := mCommon.HostAutoFields
	dataRows := []*metadata.ObjectAttDes{
		//基本信息分组
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKHostInnerIPField, PropertyName: "内网IP", IsRequired: true, IsOnly: true, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: common.PatternMultipleIP},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKHostOuterIPField, PropertyName: "外网IP", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: common.PatternMultipleIP},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "operator", PropertyName: "主要维护人", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_bak_operator", PropertyName: "备份维护人", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_asset_id", PropertyName: "固资编号", IsRequired: false, IsOnly: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_sn", PropertyName: "设备SN", IsRequired: false, IsOnly: true, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_comment", PropertyName: "备注", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_uuid", PropertyName: "UUID", IsRequired: false, IsOnly: true, Editable: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_host_type", PropertyName: "主机类型", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "物理机", Type: "text"}, {ID: "2", Name: "虚拟机", Type: "text"}, {ID: "3", Name: "Docker", Type: "text"}, {ID: "4", Name: "路由器", Type: "text"}, {ID: "5", Name: "交换机", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_level", PropertyName: "重要级别", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "重要", Type: "text"}, {ID: "2", Name: "一般", Type: "text"}, {ID: "3", Name: "非常重要", Type: "text"}, {ID: "4", Name: "不重要", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_status", PropertyName: "运行状态", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "在线", Type: "text"}, {ID: "2", Name: "离线", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_current_status", PropertyName: "当前状态", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "运营中", Type: "text"}, {ID: "2", Name: "待运营", Type: "text"}, {ID: "3", Name: "故障中", Type: "text"}, {ID: "4", Name: "重装中", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_sla", PropertyName: "SLA级别", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "L1", Type: "text"}, {ID: "2", Name: "L2", Type: "text"}, {ID: "3", Name: "L3", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKCloudIDField, PropertyName: "云区域", IsRequired: false, IsOnly: true, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleAsst, Option: common.KvMap{"value": "plat", "label": "云区域"}}, //common.FieldTypeInt, Option: "{}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_service_term", PropertyName: "质保年限", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: common.KvMap{"min": "1", "max": "10"}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "buytime", PropertyName: "采购时间", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeTime, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "expiretime", PropertyName: "过保时间", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeTime, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "inputtime", PropertyName: "上架时间", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeTime, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "outputtime", PropertyName: "下架时间", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeTime, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "firstusetime", PropertyName: "首次使用时间", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeTime, Option: ""},

		//自动发现分组
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_host_name", PropertyName: "主机名称", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKOSTypeField, PropertyName: "操作系统类型", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "Linux", Type: "text"}, {ID: "2", Name: "Windows", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKOSNameField, PropertyName: "操作系统名称", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_os_version", PropertyName: "操作系统版本", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_os_bit", PropertyName: "操作系统位数", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_cpu", PropertyName: "CPU逻辑核心数", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeInt, Option: common.KvMap{"min": "1", "max": "1000000"}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_cpu_mhz", PropertyName: "CPU频率", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeInt, Unit: "Hz", Option: common.KvMap{"min": "1", "max": "100000000"}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_cpu_module", PropertyName: "CPU型号", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_mem", PropertyName: "内存容量", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeInt, Unit: "MB", Option: common.KvMap{"min": "1", "max": "100000000"}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_disk", PropertyName: "磁盘容量", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeInt, Unit: "GB", Option: common.KvMap{"min": "1", "max": "100000000"}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_mac", PropertyName: "内网MAC", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_lan_mask", PropertyName: "内网掩码", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_lan_gateway", PropertyName: "内网网关", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: common.PatternIP},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_outer_mac", PropertyName: "外网MAC", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_outer_mask", PropertyName: "外网掩码", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_outer_gateway", PropertyName: "外网网关", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: common.PatternIP},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_host_manageip", PropertyName: "带外IP", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: common.PatternIP},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_manage_mask", PropertyName: "带外掩码", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_manage_gateway", PropertyName: "带外网关", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_manufacturer", PropertyName: "厂商", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_productName", PropertyName: "型号", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},

		//agent 没有分组
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.CreateTimeField, PropertyName: "录入时间", IsRequired: false, IsOnly: false, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeTime, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "import_from", PropertyName: "录入方式", IsRequired: false, IsOnly: false, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "excel", Type: "text"}, {ID: "2", Name: "agent", Type: "text"}, {ID: "3", Name: "api", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_agent_version", PropertyName: "Agent版本", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_agent_status", PropertyName: "Agent状态", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "正常", Type: "text"}, {ID: "2", Name: "异常", Type: "text"}, {ID: "3", Name: "未安装", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_agent_update_time", PropertyName: "Agent更新时间", IsRequired: false, IsOnly: false, Editable: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeTime, Option: ""},
	}

	return dataRows
}

// ProcRow proc structure
func ProcRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDProc
	groupPort := mCommon.ProcPort
	// groupGsekit := mCommon.Proc_gsekit_base_info
	// groupGsekitManage := mCommon.Proc_gsekit_manage_info
	dataRows := []*metadata.ObjectAttDes{
		//base info
		//&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_process_id", PropertyName: "进程ID", IsSystem: true, IsRequired: true, IsOnly: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: "{}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKAppIDField, PropertyName: "业务ID", IsAPI: true, IsRequired: true, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: "{}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKProcessNameField, PropertyName: "进程名称", IsRequired: true, IsOnly: true, IsPre: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "description", PropertyName: "进程描述", IsRequired: false, IsOnly: false, IsPre: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},

		//监听端口分组
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bind_ip", PropertyName: "绑定IP", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupPort, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "127.0.0.1", Type: "text"}, {ID: "2", Name: "0.0.0.0", Type: "text"}, {ID: "3", Name: "第一内网IP", Type: "text"}, {ID: "4", Name: "第一外网IP", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "port", PropertyName: "端口", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupPort, PropertyType: common.FieldTypeSingleChar, Option: `^((\d+-\d+)|(\d+))(,((\d+)|(\d+-\d+)))*$`, Placeholder: `单个端口：8080 </br>多个连续端口：8080-8089 </br>多个不连续端口：8080-8089,8199`},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "protocol", PropertyName: "协议", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupPort, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "TCP", Type: "text"}, {ID: "2", Name: "UDP", Type: "text"}}},

		//gsekit 基础信息
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_func_id", PropertyName: "功能ID", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_func_name", PropertyName: "功能名称", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "work_path", PropertyName: "工作路径", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "user", PropertyName: "启动用户", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "proc_num", PropertyName: "启动数量", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeInt, Option: common.KvMap{"min": "1", "max": "1000000"}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "priority", PropertyName: "启动优先级", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeInt, Option: common.KvMap{"min": "1", "max": "100"}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "timeout", PropertyName: "操作超时时长", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeInt, Option: common.KvMap{"min": "1", "max": "1000000"}},

		//gsekit 进程信息
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "start_cmd", PropertyName: "启动命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "stop_cmd", PropertyName: "停止命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "restart_cmd", PropertyName: "重启命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "face_stop_cmd", PropertyName: "强制停止命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "reload_cmd", PropertyName: "进程重载命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "pid_file", PropertyName: "PID文件路径", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "auto_start", PropertyName: "是否自动拉起", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeBool, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "auto_time_gap", PropertyName: "拉起间隔", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeInt, Option: common.KvMap{"min": "1", "max": "1000000"}},
	}
	return dataRows
}

// IdcRow app structure
func IdcRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDIdc
	dataRows := []*metadata.ObjectAttDes{
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_idc_name", PropertyName: "机房", IsRequired: true, IsOnly: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_idc_city_name", PropertyName: "城市", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_idc_short_name", PropertyName: "机房简称", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_idc_code", PropertyName: "机房代码", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_iop_name", PropertyName: "运营商", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_idc_operator", PropertyName: "负责人", IsRequired: true, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_idc_operator_tel", PropertyName: "负责人电话", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_idc_memo", PropertyName: "备注", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
	}
	return dataRows

}

// RackRow set structure
func RackRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDRack

	dataRows := []*metadata.ObjectAttDes{
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_idc_id", PropertyName: "机房ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: common.KvMap{}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_name", PropertyName: "机架", IsRequired: true, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_tor_name", PropertyName: "TOR分组", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_high", PropertyName: "机架高度", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_sn", PropertyName: "SP编号", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_power_type", PropertyName: "电路类型", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_electricity_type", PropertyName: "电流类型", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_rated_power", PropertyName: "可用功率", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_city_power", PropertyName: "市电", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_ratedi", PropertyName: "额定电流", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_maxi", PropertyName: "峰值电流", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_max_power", PropertyName: "峰值功率", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_usedi", PropertyName: "已用电流", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_used", PropertyName: "是否启用", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "0", Name: "否", Type: "text"}, {ID: "1", Name: "是", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_power_open", PropertyName: "是否开电", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "0", Name: "否", Type: "text"}, {ID: "1", Name: "是", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_use_time", PropertyName: "启用时间", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeTime, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_open_time", PropertyName: "开电时间", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeTime, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_memo", PropertyName: "备注", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
	}
	return dataRows
}

// PosRow module structure
func PosRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDPos

	dataRows := []*metadata.ObjectAttDes{
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_idc_id", PropertyName: "机房ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: common.KvMap{}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_rack_id", PropertyName: "机架ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: common.KvMap{}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_pos_name", PropertyName: "机位", IsRequired: true, IsOnly: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_pos_high", PropertyName: "机位高度", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_pos_used", PropertyName: "是否启用", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "0", Name: "否", Type: "text"}, {ID: "1", Name: "是", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_pos_power_open", PropertyName: "是否开电", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "0", Name: "否", Type: "text"}, {ID: "1", Name: "是", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_pos_use_time", PropertyName: "启用时间", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeTime, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_pos_open_time", PropertyName: "开电时间", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeTime, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_pos_memo", PropertyName: "备注", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
	}
	return dataRows
}
