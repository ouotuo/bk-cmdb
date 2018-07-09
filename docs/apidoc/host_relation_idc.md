

###  主机转移到业务内模块
* API: POST /api/{version}/hosts/idc
* API名称： transfer_host_idc
* 功能说明：
	* 中文：主机设置机房
	* English ：transfer host to idc
* input body：
```
{
    "bk_idc_id":151,
    "bk_host_id":[
        10,
        9
    ],
    "bk_pos_id":[
        170
    ]
}
```
* input字段说明:

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_idc_id| int| 是|无|机房ID |  idc ID|
| bk_host_id| int数组| 是| 无|主机 ID|host ID|
| bk_pos_id| int数组| 是| 无|机架 id| POS ID |



* output：
```
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "",
    "data": null
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null | 请求返回的数据 |the data response|
