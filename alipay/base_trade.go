/**
 * @author dengmeiyu
 * @since 20180713
 */
package alipay

import (
	"net/url"
	"strings"
	"sort"
	"fmt"
	"github.com/ttooch/payment/helper"
	"encoding/base64"
)

const (
	ALITRADE = "https://openapi.alipay.com/gateway.do?charset=utf-8"
	SUCCESS     = "SUCCESS"
	FAIL     = "FAIL"
)

type BaseAliConfig struct {
	AppId          string        `xml:"app_id" json:"app_id"`
	Method         string        `xml:"method" json:"method"`
	SignType       string        `xml:"sign_type" json:"sign_type"`
	Sign           string        `xml:"sign" json:"sign"`
	TimeStamp      string        `xml:"time_stamp" json:"time_stamp"`
	Version        string        `xml:"version" json:"version"`
	BizContent     string        `xml:"-" json:"-"`
	Charset         string        `xml:"-" json:"-"`
	NotifyUrl      string        `xml:"-" json:"-"`
	AppAuthToken   string        `xml:"-" json:"-"`
}

func (base *BaseAliConfig) RSASign(m url.Values,privateKey string) string {
	//对url.values进行排序
	if m == nil {
		m = make(url.Values, 0)
	}

	var pList = make([]string, 0, 0)
	for key := range m {
		var value = strings.TrimSpace(m.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	sort.Strings(pList)
	var src = strings.Join(pList, "&")
	fmt.Println(string(src))
	//对排序后的数据进行rsa2加密，获得sign
	b,_ := helper.RsaEncrypt([]byte(src),privateKey)

	fmt.Println("base加密：",base64.StdEncoding.EncodeToString(b))
	return base64.StdEncoding.EncodeToString(b)
}
