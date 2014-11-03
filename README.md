gotest
======
utils to upop srv req (商户接入中国银联UPOP系统), SDK rev 0.0.1.

==== 基本要求 ====


==== 使用说明 ====

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



