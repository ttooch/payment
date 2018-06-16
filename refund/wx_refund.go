package refund

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/ttooch/payment/helper"
	"github.com/ttooch/payment/weixin"
	"io/ioutil"
	"strings"
	"time"
)

const (
	WXREFUNDURL = "https://api.mch.weixin.qq.com/secapi/pay/refund"
	SUCCESS     = "SUCCESS"
)

type Error struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	PrepayId   string `xml:"prepay_id"`
	CodeUrl    string `xml:"code_url"`
	Recall     string `xml:"recall"`
}

type WxReturn struct {
	AppId     string `json:"appId"`
	TimeStamp int64  `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

type Return struct {
	AppId      string `xml:"appid"`
	SubAppId   string `xml:"sub_appid"`
	MchId      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	TradeState string `xml:"trade_state"`
}

type WxRefund struct {
	XMLName xml.Name `xml:"xml" json:"-"`
	*weixin.BaseConfig
	*WxConf
}

type WxConf struct {
	TransactionId string `xml:"transaction_id" json:"transaction_id"`
	OutRefundNo   string `xml:"out_refund_no" json:"out_refund_no"`
	TotalFee      int64  `xml:"total_fee" json:"total_fee"`
	RefundFee     int64  `xml:"refund_fee" json:"refund_fee"`
	RefundFeeType string `xml:"refund_fee_type,omitempty" json:"refund_fee_type,omitempty"`
	RefundDesc    string `xml:"refund_desc,omitempty" json:"refund_desc,omitempty"`
}

func (wx *WxRefund) Handle(conf map[string]interface{}) (interface{}, error) {
	err := wx.BuildData(conf)
	if err != nil {
		return nil, err
	}
	wx.SetSign(wx)
	ret,err := wx.sendReq(WXREFUNDURL)
	return wx.RetData(ret)
}

func (wx *WxRefund) RetData(ret []byte) (wxReturn WxReturn, err error) {

	var result struct {
		Error
		Return
	}

	xml.Unmarshal(ret, &result)

	fmt.Printf("%+v", result)

	if result.ReturnCode == SUCCESS && result.ResultCode == SUCCESS {

		if result.SubAppId != "" {
			wxReturn.AppId = result.SubAppId
		} else {
			wxReturn.AppId = result.AppId
		}

		wxReturn.TimeStamp = time.Now().Unix()

		wxReturn.NonceStr = helper.NonceStr()

		wxReturn.Package = "prepay_id=" + result.PrepayId

		wx.SetSign(wxReturn)

		wxReturn.PaySign = wx.Sign

	} else {

		return wxReturn, errors.New(result.ErrCodeDes)
	}

	return wxReturn, nil

}

func (wx *WxRefund) sendReq(reqUrl string) (b []byte,err error) {

	buffer := bytes.NewBuffer(b)

	err = xml.NewEncoder(buffer).Encode(wx)

	if err != nil {
		return
	}

	client := helper.NewTLSHttpClient("/Users/zhouchao/go/src/github.com/ttooch/payment/refund/pem/apiclient_cert.pem", "/Users/zhouchao/go/src/github.com/ttooch/payment/refund/pem/apiclient_key.pem")

	fmt.Println(buffer.String())

	httpResp, err := client.Post(reqUrl, "text/xml; charset=utf-8", buffer)

	fmt.Println(err)

	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	b, err = ioutil.ReadAll(httpResp.Body)

	return
}

func (wx *WxRefund) SetSign(data interface{}) {

	mapData := helper.Struct2Map(data)

	signStr := helper.CreateLinkString(&mapData)

	wx.Sign = wx.makeSign(signStr)
}

func (wx *WxRefund) makeSign(sign string) string {

	sign += "&key=" + wx.Md5Key

	sign = helper.Md5(sign)

	return strings.ToUpper(sign)
}

func (wx *WxRefund) BuildData(conf map[string]interface{}) error {

	b, _ := json.Marshal(conf)

	var wxConf WxConf

	json.Unmarshal(b, &wxConf)

	if wxConf.TransactionId == "" {
		return errors.New("微信订单号不能为空")
	}

	if wxConf.OutRefundNo == "" {
		return errors.New("商户退款单号不能为空")
	}

	if wxConf.RefundFee == 0 {
		return errors.New("退款金额不能为0")
	}

	if wxConf.TotalFee == 0 {
		return errors.New("订单金额不能为0")
	}

	if wxConf.RefundFee > wxConf.TotalFee {
		return errors.New("退款金额不能大于订单金额")
	}

	wx.WxConf = &wxConf

	return nil
}

func (wx *WxRefund) InitBaseConfig(config *weixin.BaseConfig) {
	config.NonceStr = helper.NonceStr()

	wx.BaseConfig = config
}
