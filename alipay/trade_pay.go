/**
 * @author dengmeiyu
 * @since 20180711
 */
package alipay

import (
	"encoding/xml"
	"github.com/ttooch/payment/helper"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"fmt"
	"errors"
	"time"
)

type AliPayTradePayResponse struct {
	AliPayTradePay struct {
		Code                string `json:"code"`
		Msg                 string `json:"msg"`
		SubCode             string `json:"sub_code"`
		SubMsg              string `json:"sub_msg"`
		BuyerLogonId        string `json:"buyer_logon_id"`        // 买家支付宝账号
		BuyerPayAmount      string `json:"buyer_pay_amount"`      // 买家实付金额，单位为元，两位小数。
		BuyerUserId         string `json:"buyer_user_id"`         // 买家在支付宝的用户id
		CardBalance         string `json:"card_balance"`          // 支付宝卡余额
		DiscountGoodsDetail string `json:"discount_goods_detail"` // 本次交易支付所使用的单品券优惠的商品优惠信息
		GmtPayment          string `json:"gmt_payment"`
		InvoiceAmount       string `json:"invoice_amount"` // 交易中用户支付的可开具发票的金额，单位为元，两位小数。
		OutTradeNo          string `json:"out_trade_no"`   // 创建交易传入的商户订单号
		TradeNo             string `json:"trade_no"`       // 支付宝交易号
		PointAmount         string `json:"point_amount"`   // 积分支付的金额，单位为元，两位小数。
		ReceiptAmount       string `json:"receipt_amount"` // 实收金额，单位为元，两位小数
		StoreName           string `json:"store_name"`     // 发生支付交易的商户门店名称
		TotalAmount         string `json:"total_amount"`   // 发该笔退款所对应的交易的订单金额
	} `json:"alipay_trade_pay_response"`
	Sign string `json:"sign"`
}

type AliConf struct {
	OutTradeNo         string `xml:"out_trade_no" json:"out_trade_no"`
	Scene              string `xml:"scene" json:"scene"`
	AuthCode           string `xml:"auth_code" json:"auth_code"`
	ProductCode        string `xml:"product_code,omitempty" json:"product_code,omitempty"`
	Subject            string `xml:"subject" json:"subject"`
	BuyerId            string `xml:"buyer_id" json:"buyer_id"`
	SellerId           string `xml:"seller_id" json:"seller_id"`
	TotalAmount        string `xml:"total_amount" json:"total_amount"`
	TransCurrency      string `xml:"trans_currency,omitempty" json:"trans_currency,omitempty"`
	SettleCurrency     string `xml:"settle_currency,omitempty" json:"settle_currency,omitempty"`
	DiscountableAmount string `xml:"discountable_amount,omitempty" json:"discountable_amount,omitempty"`
	Body               string `xml:"body" json:"body"`
}

type AliTrade struct {
	XMLName xml.Name `xml:"xml" json:"-"`
	*BaseAliConfig
	*AliConf
}

func (tra *AliTrade) Handle(conf map[string]interface{}, privateKey string, aliPublicKey string) (*AliPayTradePayResponse, error) {

	err := tra.BuildData(conf)
	if err != nil {
		return nil, err
	}
	ret, err := tra.sendReq(ALITRADE, tra, privateKey)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(ret))
	return tra.RetData(ret, privateKey, aliPublicKey)
}

func (tra *AliTrade) RetData(ret []byte, privateKey string,aliPublicKey string) (re *AliPayTradePayResponse, err error) {

	result := new(AliPayTradePayResponse)
	json.Unmarshal(ret, result)

	b, _ := json.Marshal(result.AliPayTradePay)

	//TODO 调用之后验签 验签方法还要改一下
	err = helper.RSAVerify([]byte(b), result.Sign, aliPublicKey)
	if err != nil {
		fmt.Println(err.Error())
		//return nil ,err
	}
	
	if result.AliPayTradePay.Code != "10000" && result.AliPayTradePay.Msg != "Success" {
		if result.AliPayTradePay.Code == "10003" {

			app := new(AliQuery)

			app.InitBaseConfig(&BaseAliConfig{
				AppId:    tra.AppId,
				AppAuthToken:tra.AppAuthToken,
			})

			ret, err := app.Handle(map[string]interface{}{
				"out_trade_no":   result.AliPayTradePay.OutTradeNo,

			},privateKey,aliPublicKey)

			if err != nil {
				return result, errors.New("调用支付宝查询接口失败：" + err.Error())
			}

			if ret.AliPayTradeQuery.TradeStatus =="TRADE_SUCCESS" {
				return result, nil
			}
		}
		return result, errors.New("支付宝条码支付失败：" + result.AliPayTradePay.SubMsg)
	}
	return result, nil

}

func (tra *AliTrade) sendReq(reqUrl string, pay interface{}, privateKey string) (b []byte, err error) {

	client := helper.NewHttpClient()
	var data = url.Values{}

	data.Add("app_id", tra.AppId)
	data.Add("method", tra.Method)
	data.Add("charset", tra.Charset)
	data.Add("sign_type", tra.SignType)
	data.Add("timestamp", tra.TimeStamp)
	data.Add("version", tra.Version)
	data.Add("biz_content", tra.BizContent)
	data.Add("sign", tra.RSASign(data, privateKey))

	httpResp, err := client.PostForm(reqUrl, data)

	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	b, err = ioutil.ReadAll(httpResp.Body)

	return

}

func (tra *AliTrade) BuildData(conf map[string]interface{}) error {

	b, _ := json.Marshal(conf)

	countInfo := string(b)

	tra.BaseAliConfig.BizContent = countInfo

	return nil
}

func (tra *AliTrade) InitBaseConfig(config *BaseAliConfig) {

	config.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
	config.Method = "alipay.trade.pay"
	config.Charset = "UTF-8"
	config.Version = "1.0"
	config.SignType = "RSA2"
	tra.BaseAliConfig = config
}
