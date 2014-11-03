// upop utils
// @description upop utils is an open-source for the Go programming language.
// @authors     yif

package upop

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"
)

const (
	Version              = "1.0.0"
	Charset              = "UTF-8"
	MerCode              = "105550149170027"
	MerFrontEndUrl       = ""
	MerBackEndUrl        = "http://www.yourdomain.com/your_path/yourBackEndUrl"
	SignType             = "MD5"
	SignType_SHA1withRSA = "SHA1withRSA"
	SecurityKey          = "88888888"

	TransType_01 = "01"
	Currency_CN  = "156"

	// 基础网址（请按相应环境切换）
	/* 前台交易测试环境 */
	UPOP_BASE_URL = "http://58.246.226.99/UpopWeb/api/"

	/* 前台交易PM环境（准生产环境） */
	//UPOP_BASE_URL = "https://www.epay.lxdns.com/UpopWeb/api/"

	/* 前台交易生产环境 */
	//UPOP_BASE_URL = "https://unionpaysecure.com/api/"

	/* 后台交易测试环境 */
	UPOP_BSPAY_BASE_URL = "http://58.246.226.99/UpopWeb/api/"

	/* 后台交易PM环境（准生产环境） */
	//UPOP_BSPAY_BASE_URL = "https://www.epay.lxdns.com/UpopWeb/api/"

	/* 后台交易生产环境 */
	//UPOP_BSPAY_BASE_URL = "https://besvr.unionpaysecure.com/api/"

	/* 查询交易测试环境 */
	UPOP_QUERY_BASE_URL = "http://58.246.226.99/UpopWeb/api/"

	/* 查询交易PM环境（准生产环境） */
	//UPOP_QUERY_BASE_URL = "https://www.epay.lxdns.com/UpopWeb/api/"

	/* 查询交易生产环境 */
	//UPOP_QUERY_BASE_URL = "https://query.unionpaysecure.com/api/";

	// 支付网址
	GateWay = UPOP_BASE_URL + "Pay.action"
)

func Md5Sign(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//func SortFor(v UpopPkgParams, ff func(int32, string, string)) {
//	if v == nil {
//		return
//	}
//	keys := make([]string, 0, len(v))
//	for k := range v {
//		keys = append(keys, k)
//	}
//	sort.Strings(keys)
//	var idx int32 = 0
//	for _, k := range keys {
//		idx++
//		vs := v[k]
//		ff(idx, k, vs)
//	}
//}

func SignUrlEncode(v UpopPkgParams) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := url.QueryEscape(k) + "="
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(prefix)
		//buf.WriteString(url.QueryEscape(v))
		buf.WriteString(vs)
	}
	return buf.String()
}

func ParseQuery(m UpopPkgParams, query string) (err error) {
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		m[key] = value
	}
	return err
}

func SignMap(req UpopPkgParams, signMethod string, securityKey string) string {
	/*
		idx := 1
			v := url.Values{}
			for key, item := range req {
				println(key + ",  " + item)
				if idx == 1 {
					v.Set(key, item)
				} else {
					v.Add(key, item)
				}
				idx = idx + 1
			}
	*/

	if SignType == signMethod {
		strBeforeMd5 := SignUrlEncode(req) + "&" + Md5Sign(securityKey)
		//println(strBeforeMd5)
		//println(Md5Sign(strBeforeMd5))
		return Md5Sign(strBeforeMd5)
	} else if SignType_SHA1withRSA == signMethod {

	}
	return ""

}

type UpopPkgParams map[string]string

func (c UpopPkgParams) SignMap(signMethod string) string {
	return SignMap(c, signMethod, SecurityKey)
}

func (c UpopPkgParams) clone() UpopPkgParams {
	/*
		newMap := UpopPkgParams{}
		for k, v := range c {
			newMap[k] = v
		}
		return newMap
	*/
	return c
}

func (c UpopPkgParams) And(cond UpopPkgParams) UpopPkgParams {
	c = c.clone()
	if cond != nil {
		for k, v := range cond {
			c[k] = v
		}
	}
	return c
}

func (c UpopPkgParams) OrNot(cond UpopPkgParams) UpopPkgParams {
	c = c.clone()
	if cond != nil {
		for k, _ := range cond {
			delete(c, k)
		}
	}
	return c
}

func (c UpopPkgParams) Encode(signMethod string) string {
	c.OrNot(UpopPkgParams{
		"signature":  "",
		"signMethod": ""})
	var signstr string = SignMap(c, signMethod, SecurityKey)
	c.And(UpopPkgParams{
		"signature":  signstr,
		"signMethod": signMethod})
	return signstr
}

func (c UpopPkgParams) Decode(restr string) (rr error) {
	rr = ParseQuery(c, restr)
	return
}

func (c UpopPkgParams) CheckSecurity(key string) bool {
	st := c["signature"]
	sm := c["signMethod"]

	//println("st:", st)
	//println("sm:", sm)

	if SignType == sm {
		c.OrNot(UpopPkgParams{
			"signature":  "",
			"signMethod": ""})
		signstr := SignMap(c, sm, key)
		c.And(UpopPkgParams{
			"signature":  st,
			"signMethod": sm})
		//println("ss:", signstr)
		return signstr == st
	} else if SignType_SHA1withRSA == sm {

	}
	return true
}

func (c UpopPkgParams) Get(key string) string {
	if v, ok := c[key]; ok {
		return v
	}
	return "<nil>"
}

func (c UpopPkgParams) RespCode() string {
	return c.Get("respCode")
}

func (c UpopPkgParams) RespMsg() string {
	return c.Get("respMsg")
}

func (c UpopPkgParams) Tn() string {
	return c.Get("tn")
}

func NewPayReq() UpopPkgParams {
	return UpopPkgParams{
		"version":          Version,
		"charset":          Charset,
		"transType":        TransType_01, //交易类型
		"merId":            MerCode,      //商户代码
		"backEndUrl":       MerBackEndUrl,
		"frontEndUrl":      "",
		"acqCode":          "",
		"orderTime":        "",          //交易开始时间
		"orderTimeout":     "",          //订单超时时间
		"orderNumber":      "",          //商户订单号
		"orderAmount":      "",          //交易金额
		"orderCurrency":    Currency_CN, //交易币种
		"orderDescription": "",
		"merReserved":      "",
		"reqReserved":      "",
		"sysReserved":      "",
	}
}

func NewRayRep() UpopPkgParams {
	return UpopPkgParams{
		"version": Version,
		"charset": Charset,
		//"signMethod": "",
		//"signature": "",
		"transType":   "",
		"tn":          "",
		"respCode":    "",
		"respMsg":     "",
		"merReserved": "",
		"reqReserved": "",
		"sysReserved": "",
	}
}

/*
now := time.Now()
		d, _ := time.ParseDuration("30m")
		ot := now.Format("2006-01-02 15:04:05")
		oto := now.Add(d).Format("2006-01-02 15:04:05")

		var req = upop.UpopPkgParams{
			"version":          upop.Version,
			"charset":          upop.Charset,
			"transType":        "01",         //交易类型
			"merId":            upop.MerCode, //商户代码
			"backEndUrl":       upop.MerBackEndUrl,
			"frontEndUrl":      "",
			"acqCode":          "",
			"orderTime":        ot,
			"orderTimeout":     oto, //交易超时时间
			"orderNumber":      "",
			"orderAmount":      orderAmount,
			"orderCurrency":    "156", //交易币种
			"orderDescription": "",
			"merReserved":      "",
			"reqReserved":      "",
			"sysReserved":      "",
		}
		beego.Debug("old:", req)
		signstr := req.Encode(upop.SignType)

		beego.Debug("signstr:", signstr)
		beego.Debug("new:", req)

		var tgt = upop.UpopPkgParams{
			"version": upop.Version,
			"charset": upop.Charset,
			//"signMethod": "",
			//"signature": "",
			"transType":   "01",
			"tn":          "A0000000000001",
			"respCode":    "00",
			"respMsg":     "ok",
			"merReserved": "",
			"reqReserved": "",
			"sysReserved": "",
		}
		strurl := upop.SignUrlEncode(tgt)
		strurlbef := strurl + "&" + upop.Md5Sign("asdf1")
		strurl = strurl + "&signMethod=MD5&signature=" + upop.Md5Sign(strurlbef)

		var rep = upop.UpopPkgParams{}
		err := rep.Decode(strurl)
		cs := rep.CheckSecurity("asdf")

		beego.Debug("rep:", rep)
		beego.Debug("err:", err)
		beego.Debug("cs:", cs)
*/
