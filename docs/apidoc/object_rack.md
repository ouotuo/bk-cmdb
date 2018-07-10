
### 创建机架

- API: POST  /api/{version}/rack/{bk_idc_id}   
- API 名称：create_rack
- 功能说明：
	- 中文： 新建机架
	- English：create rack

- input body
``` json
{

    "bk_rack_name":"",
    "bk_parent_id":0,
    "bk_supplier_account":"",
    "bk_idc_id":1,
    "default":0
}
```

**注:以上 JSON 数据中各字段为必填字段或内置字段，它们在示例中的值仅为示例数据。**

- input 参数说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_rack_id|int|是|无|机架ID|the rack id|
|bk_idc_id|int|是|无|机房ID|business ID|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_rack_name|string|是|无|机架名字 |rack name|
|bk_rack_desc|string|否|无|机架描述|the rack description|

 **注: 用户自定义的字段也可以作为参数传入。**


### 删除机架

- API: DELETE  /api/{version}/rack/{bk_idc_id}/{bk_rack_id}   
- API 名称：delete_rack
- 功能说明：
	- 中文： 删除机架
	- English：delete rack

- input body

    无

- input参数说明：

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_idc_id|int|是|无|机房ID|business ID|
|bk_rack_id|int|是|无|机架ID|the rack id|


- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": "success"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|操作结果|the result|

### 更新机架
- API： PUT /api/{version}/rack/{bk_idc_id}/{bk_rack_id}   
- API 名称：update_rack
- 功能说明：
	- 中文： 更新机架
	- English：update rack

- input body

``` json
{

    "bk_rack_name":"",
    "default":0
}
```

**注:以上 JSON 数据中各字段为必填字段或内置字段，它们在示例中的值仅为示例数据。**

- input 参数说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_rack_id|int|是|无|机架ID|the rack id|
|bk_idc_id|int|是|无|机房ID|business ID|
|bk_rack_name|string|否|无|机架名字 |rack name|
|bk_capacity|int|否|无|设计容量|the design the capacity|
|description|string|否|无|备注|the remark|
|bk_service_status|enum|否|开放|服务状态:开发/关闭|the service status:open/close|
|bk_rack_env|enum|否|正式|环境类型：测试/体验/正式|environment type:test/experience/formal|
|bk_rack_desc|string|否|无|机架描述|the rack description|


 **注: 用户在使用的时候可以为每个Set增加字段的数量，这些自定义的字段也可以作为参数传入。**

- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": "success"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|操作结果|the result|

### 查询机架

- API： POST /api/{version}/rack/search/{bk_supplier_account}/{bk_idc_id}   
- API 名称：search_rack
- 功能说明：
	- 中文： 查询机架
	- English：search rack

-  input body:
``` json
{
    "fields":[
        "bk_rack_name"
    ],
    "page":{
        "start":0,
        "limit":100,
        "sort":"bk_rack_name"
    },
    "condition":{
        "bk_rack_name":"rack_new"
    }
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input参数说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
| bk_supplier_account| string| 是| 无|开发商账号|supplier account code|
| bk_idc_id| int| 是|无|机房ID |  business ID|
| page| object| 是|无|分页参数 |page parameter|
| fields| array | 是| 无|查询字段|search fields|
| condition|  object| 是| 无|查询条件|search condition|

page 参数说明：

|名称|类型|必填| 默认值 | 说明 | Description|
|---|---| --- |---  | --- | ---|
| start|int|是|无|记录开始位置 |start record|
| limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
| sort| string| 否| 无|排序字段|the field for sort|

fields参数说明：

|名称|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_parent_id|int|否|无|父节点的ID|the parent inst identifier|
|bk_rack_id|int|是|无|机架ID|the rack id|
|bk_rack_name|string|否|无|机架名字 |rack name|
|bk_capacity|int|否|无|设计容量|the design the capacity|
|description|string|否|无|备注|the remark|
|bk_service_status|enum|否|开放|服务状态:开发/关闭|the service status:open/close|
|bk_rack_env|enum|否|正式|环境类型：测试/体验/正式|environment type:test/experience/formal|
|bk_rack_desc|string|否|无|机架描述|the rack description|

**注:所有字段均为Set定义的字段，这些字段包括预置字段，也包括用户自定义字段。**

condition 参数说明：

|名称|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_parent_id|int|否|无|父节点的ID|the parent inst identifier|
|bk_rack_id|int|是|无|机架ID|the rack id|
|bk_rack_name|string|否|无|机架名字 |rack name|
|bk_capacity|int|否|无|设计容量|the design the capacity|
|description|string|否|无|备注|the remark|
|bk_service_status|enum|否|开放|服务状态:开发/关闭|the service status:open/close|
|bk_rack_env|enum|否|正式|环境类型：测试/体验/正式|environment type:test/experience/formal|
|bk_rack_desc|string|否|无|机架描述|the rack description|

**注:所有字段均为Set定义的字段，这些字段包括预置字段，也包括用户自定义字段。**

- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": {
        "count": 1,
        "info": [
            {
                "bk_rack_name": "内置模块集"
            }
        ]
    }
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|操作结果|the result|

data 说明

| 字段|类型|说明|Description|
|---|---|---|---|
|count|int|数据条数|the data item count|
|info|array|数据集合|the data array|

info 说明

| 字段|类型|说明|Description|
|---|---|---|---|
|bk_parent_id|int|父节点的ID|the parent inst identifier|
|bk_rack_id|int|机架ID|the rack id|
|bk_rack_name|string|机架名字 |rack name|
|bk_capacity|int|设计容量|the design the capacity|
|description|string|备注|the remark|
|bk_service_status|enum|服务状态:开发/关闭|the service status:open/close|
|bk_rack_env|enum|环境类型：测试/体验/正式|environment type:test/experience/formal|
|bk_rack_desc|string|机架描述|the rack description|

**注：此处按照fields所指定的字段进行配置，所有字段均为Set定义的字段，这些字段包括预置字段，也包括用户自定义字段。**

