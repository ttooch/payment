package weixin

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/scholar-ink/payment/helper"
	"time"
)

type PubCharge struct {
	XMLName xml.Name `xml:"xml" json:"-"`
	*PubConf
	BaseCharge
}

type PubConf struct {
	Openid         string `xml:"openid,omitempty" json:"openid,omitempty"`
	SubOpenid      string `xml:"sub_openid,omitempty" json:"sub_openid,omitempty"`
	Body           string `xml:"body" json:"body"`
	Detail         string `xml:"detail,omitempty" json:"detail,omitempty"`
	Attach         string `xml:"attach,omitempty" json:"attach,omitempty"`
	OutTradeNo     string `xml:"out_trade_no" json:"out_trade_no"`
	FeeType        string `xml:"fee_type,omitempty" json:"fee_type,omitempty"`
	TotalFee       int64  `xml:"total_fee" json:"total_fee"`
	SpbillCreateIp string `xml:"spbill_create_ip" json:"spbill_create_ip"`
	TimeStart      string `xml:"time_start,omitempty" json:"time_start,omitempty"`
	TimeExpire     string `xml:"time_expire,omitempty" json:"time_expire,omitempty"`
	GoodsTag       string `xml:"goods_tag,omitempty" json:"goods_tag,omitempty" `
	NotifyUrl      string `xml:"notify_url" json:"notify_url"`
	TradeType      string `xml:"trade_type" json:"trade_type"`
}

type PubReturn struct {
	AppId     string `json:"appId"`
	TimeStamp int64  `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

func (app *PubCharge) Handle(conf map[string]interface{}) (interface{}, error) {
	app.BuildData(conf)
	app.SetSign(app)
	ret := app.SendReq(app)
	return app.RetData(ret)
}

func (app *PubCharge) RetData(ret []byte) (pubReturn PubReturn, err error) {
	var result struct {
		Error
		Return
	}

	xml.Unmarshal(ret, &result)

	if result.ReturnCode == SUCCESS && result.ResultCode == SUCCESS {

		if result.SubAppId != "" {
			pubReturn.AppId = result.SubAppId
		} else {
			pubReturn.AppId = result.AppId
		}

		pubReturn.TimeStamp = time.Now().Unix()

		pubReturn.NonceStr = helper.NonceStr()

		pubReturn.Package = "prepay_id=" + result.PrepayId

		pubReturn.SignType = app.SignType

		app.SetSign(pubReturn)

		pubReturn.PaySign = app.Sign

	} else {

		return pubReturn, errors.New(result.ReturnMsg + result.ErrCodeDes)
	}

	return pubReturn, nil

}

func (app *PubCharge) BuildData(conf map[string]interface{}) {

	b, _ := json.Marshal(conf)

	var pubConf PubConf

	json.Unmarshal(b, &pubConf)

	app.PubConf = &pubConf

	app.PubConf.TradeType = "JSAPI"
}
