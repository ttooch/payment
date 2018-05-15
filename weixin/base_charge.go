package weixin

import (
	"bytes"
	"encoding/xml"
	"github.com/scholar-ink/payment/helper"
	"io/ioutil"
	"strings"
)

const (
	ReqUrl  = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	SUCCESS = "SUCCESS"
)

type Error struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	PrepayId   string `xml:"prepay_id"`
	CodeUrl    string `xml:"code_url"`
}

type Return struct {
	AppId      string `xml:"appid"`
	SubAppId   string `xml:"sub_appid"`
	MchId      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
}

type BaseCharge struct {
	*BaseConfig
}

type BaseConfig struct {
	AppId     string `xml:"appid" json:"appid"`
	MchId     string `xml:"mch_id" json:"mch_id"`
	SubAppId  string `xml:"sub_appid,omitempty" json:"sub_appid,omitempty"`
	SubMchId  string `xml:"sub_mch_id,omitempty" json:"sub_mch_id,omitempty"`
	TimeStart string `xml:"time_start" json:"time_start"`
	Md5Key    string `xml:"-" json:"-"`
	SignType  string `xml:"sign_type" json:"sign_type"`
	Sign      string `xml:"sign" json:"sign"`
	NonceStr  string `xml:"nonce_str" json:"nonce_str"`
}

func (base *BaseCharge) SendReq(pay interface{}) (b []byte) {

	buffer := bytes.NewBuffer(b)

	err := xml.NewEncoder(buffer).Encode(pay)

	if err != nil {
		panic(err)
	}

	client := helper.DefaultHttpClient

	httpResp, err := client.Post(ReqUrl, "text/xml; charset=utf-8", buffer)

	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	b, err = ioutil.ReadAll(httpResp.Body)

	return
}

func (base *BaseCharge) SetSign(pay interface{}) {

	mapData := helper.Struct2Map(pay)

	signStr := helper.CreateLinkString(&mapData)

	base.Sign = base.makeSign(signStr)

}

func (base *BaseCharge) makeSign(sign string) string {

	switch base.SignType {

	case "MD5":
		sign += "&key=" + base.Md5Key

		sign = helper.Md5(sign)
	}

	return strings.ToUpper(sign)
}

func (base *BaseCharge) initBaseConfig(config *BaseConfig) {
	config.NonceStr = helper.NonceStr()
	base.BaseConfig = config
}
