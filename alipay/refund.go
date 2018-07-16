/**
 * @author dengmeiyu
 * @since 20180713
 */
package alipay

import (
	"time"
	"encoding/xml"
	"fmt"
	"github.com/ttooch/payment/helper"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"errors"
)

type AliRefund struct {
	XMLName xml.Name `xml:"xml" json:"-"`
	*BaseAliConfig
	*AliRefundConf
}

type AliRefundConf struct {
	OutTradeNo string `xml:"out_trade_no" json:"out_trade_no"`
	TradeNo  string `xml:"trade_no" json:"trade_no"`
	RefundAmount string `xml:"refund_amount" json:"refund_amount"`
	RefundCurrency string `xml:"refund_currency" json:"refund_currency"`
	RefundReason string `xml:"refund_reason,omitempty" json:"refund_reason,omitempty"`

}

func (tra *AliRefund) Handle(conf map[string]interface{},privateKey string, aliPublicKey string) (string, error) {

	err := tra.BuildData(conf)
	if err != nil {
		return "", err
	}
	ret, err := tra.sendReq(ALITRADE,tra,privateKey)
	if err != nil {
		return "", err
	}
	fmt.Println(string(ret))
	return tra.RetData(ret, aliPublicKey)
}

func (tra *AliRefund) RetData(ret []byte,aliPublicKey string) (re string, err error) {

	result := new(AliPayTradePayResponse)
	json.Unmarshal(ret, result)

	b,_:=json.Marshal(result.AliPayTradePay)

	//TODO 调用之后验签 验签方法还要改一下
	err = helper.RSAVerify([]byte(b),result.Sign, aliPublicKey)
	if err != nil {
		fmt.Println(err.Error())
		//return nil ,err
	}

	if result.AliPayTradePay.Code != "10000" && result.AliPayTradePay.Msg != "Success"{
		return "", errors.New("支付宝退款失败："+result.AliPayTradePay.SubMsg)
	}
	return result.AliPayTradePay.TradeNo, nil

}

func (tra *AliRefund) sendReq(reqUrl string, pay interface{},privateKey string) (b []byte, err error) {

	client := helper.NewHttpClient()
	var data = url.Values{}

	data.Add("app_id", tra.AppId)
	data.Add("method", tra.Method)
	data.Add("charset", tra.Charset)
	data.Add("sign_type", tra.SignType)
	data.Add("timestamp", tra.TimeStamp)
	data.Add("version", tra.Version)
	data.Add("biz_content",tra.BizContent)
	data.Add("sign", tra.RSASign(data,privateKey))

	httpResp, err := client.PostForm(reqUrl,data)

	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	b, err = ioutil.ReadAll(httpResp.Body)

	return

}


func (tra *AliRefund) BuildData(conf map[string]interface{}) error {

	b, _ := json.Marshal(conf)

	countInfo := string(b)

	tra.BaseAliConfig.BizContent = countInfo

	return nil
}


func (tra *AliRefund) InitBaseConfig(config *BaseAliConfig) {

	config.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
	config.Method = "alipay.trade.refund"
	config.Charset = "utf-8"
	config.Version = "1.0"
	config.SignType = "RSA2"
	tra.BaseAliConfig = config
}

