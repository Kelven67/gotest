gotest
======

### u1) upop utils

utils golang sdk to upop srv req (商户接入中国银联UPOP系统), SDK rev 0.0.1.

==== 基本要求 ====

golang 1.2+ 

==== 使用说明 ====
```go
/*
	var req = upop.NewPayReq().And(upop.UpopPkgParams{ 
			"orderTime":    ot, 
			"orderTimeout": oto, 
			"orderNumber":  orderNumber, 
			"orderAmount":  orderAmount, 
		})
*/

func (this *ServiceController) Pay() {
	beego.Debug("---------2.20 Pay: ")
	now := time.Now()
	result := PayTnResult{Timestamp: now.Unix(), Code: CODE_FAIL, Msg: MSG_TOKEN_IS_INVALID}

	//check token
	token := this.GetReqString("token", "")
	uid, _, _, _, _, _ := this.validationToken(token)
	if uid > 0 {
		orderNumber := this.GetReqString("orderNumber", "")
		orderAmount := this.GetReqString("orderAmount", "")
		//check orderNumber
		if len(orderNumber) <= 0 {
			result.Msg = "ReqOrderNumberIsnull"
			goto ret
		}
		//check orderAmount
		if len(orderAmount) <= 0 {
			result.Msg = "ReqOrderAmountIsnull"
			goto ret
		}

		d, _ := time.ParseDuration("30m")
		ot := now.Format("20060102150405") //format -> 2006-01-02 15:04:05
		oto := now.Add(d).Format("20060102150405")

		var req = upop.UpopPkgParams{
			"version":          upop.Version,
			"charset":          upop.Charset,
			"transType":        upop.TransType_01, //交易类型
			"merId":            upop.MerCode,      //商户代码
			"backEndUrl":       upop.MerBackEndUrl,
			"frontEndUrl":      "",
			"acqCode":          "",
			"orderTime":        ot,  //交易开始时间
			"orderTimeout":     oto, //订单超时时间
			"orderNumber":      orderNumber,
			"orderAmount":      orderAmount,
			"orderCurrency":    upop.Currency_CN, //交易币种
			"orderDescription": "",
			"merReserved":      "",
			"reqReserved":      "",
			"sysReserved":      "",
		}
		//beego.Debug("req(1):", req)
		req.Encode(upop.SignType) //signstr
		//beego.Debug("req(2):", req)

		//req UPMP srv
		var strurl string
		var err error
		b := httplib.Post(upop.GateWay)
		//b.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		//upop.SortFor(req, func(idx int32, k, v string) {
		//	beego.Debug("k:", k, ",v:", v)
		//	b.Param(k, v)
		//})
		for k, v := range req {
			b.Param(k, v)
		}
		if false {
			strurl, err = b.SetTimeout(100*time.Second, 30*time.Second).String() //strurl
			if err != nil {
				beego.Debug("err:", err)
				result.Msg = "ReqUpopIsException"
				goto ret
			}
		} else {
			strurl = "charset=UTF-8&merReserved=&reqReserved=&respCode=00&respMsg=ok&sysReserved=&tn=A0000000000001&transType=01&version=1.0.0&signMethod=MD5&signature=e6d1877889422a8b21d9d1a2cd91dfb8"
		}

		beego.Debug("strurl:", strurl)

		var rep = upop.UpopPkgParams{}
		err = rep.Decode(strurl)
		cs := rep.CheckSecurity(upop.SecurityKey)

		//check reps sign
		if cs == false {
			result.Msg = "UpopRepCheckSecurityError"
			goto ret
		} else if cs == true && rep.RespCode() == "00" && err == nil {
			//do the business logic
			beego.Debug("tn", rep.Tn())
			result.Tn = rep.Tn()
			result.Code = CODE_SUCC
			result.Msg = MSG_SUCCESS
			goto ret
		}

		result.Msg = "Error_" + rep.RespCode() + "_" + rep.RespMsg()
		goto ret

	}
ret:
	this.Data["json"] = &result
	this.ServeJson()
}
```


* * *

<h2 id="user">二、服务相关接口</h2>
<h3 id="user1">2.1 服务列表</h3>
**接口说明：**  `服务首页，每个功能项的列表。`

**接口地址：**  `api/service/list?serviceId=[value]&token=[token]`

**请求方式：**  `GET`

**请求参数：**
>
| 参数名称       | 类型      	|必填    	| 说明              |
| ----------- 	|--------	|------		| -----------------|
| token      	| 字符串 	|		是	|          |
| serviceId      |字符串	|   	        是    |          |


**返回字段：**
>
| 名称       | 类型      | 说明              |
| ----------- 	|--------	| -----------------|
| serviceList      	| 数组 	|	集采类数组 |


**错误代码表：**
>
| 错误代码       |  说明              |
| ----------- 	|-----------------|
| TokenIsInvalid      	| 令牌无效  |	

**返回数据格式：**	`text/json`
	
**返回数据示例：**

```
{
    "serviceList": [
        {
            "title": "关于2014年5月份通勤车月票预定和领取的通知_0",
            "id": "dffc4a498ce94956b31011b79ad8c4fe",
            "targetUrl": "http://10.10.873.174:8080/api/web/service.htm",
            "dateTime": "03-19 09:15"
        },
        {
            "title": "关于2014年5月份通勤车月票预定和领取的通知_1",
            "id": "7e00b769596848559943802d836ba362",
            "targetUrl": "http://10.10.873.174:8080/api/web/service.htm",
            "dateTime": "03-19 09:15"
        },
        {
            "title": "关于2014年5月份通勤车月票预定和领取的通知_2",
            "id": "4c421feb992c4528b57e1747799f32b5",
            "targetUrl": "http://10.10.873.174:8080/api/web/service.htm",
            "dateTime": "03-19 09:15"
        },
        {
            "title": "关于2014年5月份通勤车月票预定和领取的通知_3",
            "id": "2218f96b0f854317b5c6304e9c5cd589",
            "targetUrl": "http://10.10.873.174:8080/api/web/service.htm",
            "dateTime": "03-19 09:15"
        },
        {
            "title": "关于2014年5月份通勤车月票预定和领取的通知_4",
            "id": "5f9812976b9b4358a4ab652fe7813609",
            "targetUrl": "http://10.10.873.174:8080/api/web/service.htm",
"dateTime": "03-19 09:15"
        }
    ],
 "total": 5,
"currentTimestamp ": 1399881385,    "code": 0,
"msg": ""
}

```

**返回代码说明：**

**备注：**
	

* * *

