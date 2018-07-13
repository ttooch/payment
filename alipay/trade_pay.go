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
	"time"
	"net/url"
	"sort"
	"fmt"
	"encoding/base64"
	"github.com/smartwalle/alipay/encoding"
	"crypto"
	"strings"
)
const (
	ALITRADE = "https://openapi.alipay.com/gateway.do"
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

type AliPayTradePayResponse struct {
	AliPayTradePay struct {
		Code                string           `json:"code"`
		Msg                 string           `json:"msg"`
		SubCode             string           `json:"sub_code"`
		SubMsg              string           `json:"sub_msg"`
		BuyerLogonId        string           `json:"buyer_logon_id"`           // 买家支付宝账号
		BuyerPayAmount      string           `json:"buyer_pay_amount"`         // 买家实付金额，单位为元，两位小数。
		BuyerUserId         string           `json:"buyer_user_id"`            // 买家在支付宝的用户id
		CardBalance         string           `json:"card_balance"`             // 支付宝卡余额
		DiscountGoodsDetail string           `json:"discount_goods_detail"`    // 本次交易支付所使用的单品券优惠的商品优惠信息
		GmtPayment          string           `json:"gmt_payment"`
		InvoiceAmount       string           `json:"invoice_amount"`                // 交易中用户支付的可开具发票的金额，单位为元，两位小数。
		OutTradeNo          string           `json:"out_trade_no"`                  // 创建交易传入的商户订单号
		TradeNo             string           `json:"trade_no"`                      // 支付宝交易号
		PointAmount         string           `json:"point_amount"`                  // 积分支付的金额，单位为元，两位小数。
		ReceiptAmount       string           `json:"receipt_amount"`                // 实收金额，单位为元，两位小数
		StoreName           string           `json:"store_name"`                    // 发生支付交易的商户门店名称
		TotalAmount         string           `json:"total_amount"`                  // 发该笔退款所对应的交易的订单金额
	} `json:"alipay_trade_pay_response"`
	Sign string `json:"sign"`
}

type AliConf struct {
	OutTradeNo   string `xml:"out_trade_no" json:"out_trade_no"`
	Scene        string `xml:"scene" json:"scene"`
	AuthCode     string `xml:"auth_code" json:"auth_code"`
	ProductCode  string `xml:"product_code,omitempty" json:"product_code,omitempty"`
	Subject      string `xml:"subject" json:"subject"`
	BuyerId      string `xml:"buyer_id" json:"buyer_id"`
	SellerId     string `xml:"seller_id" json:"seller_id"`
	TotalAmount  string `xml:"total_amount" json:"total_amount"`
	TransCurrency string `xml:"trans_currency,omitempty" json:"trans_currency,omitempty"`
	SettleCurrency string `xml:"settle_currency,omitempty" json:"settle_currency,omitempty"`
	DiscountableAmount string `xml:"discountable_amount,omitempty" json:"discountable_amount,omitempty"`
	Body string `xml:"body" json:"body"`

}

type AliTrade struct {
	XMLName xml.Name `xml:"xml" json:"-"`
	*BaseAliConfig
	*AliConf
}


type TradeReturn struct {
	TradeNo string `xml:"trade_no" json:"trade_no"`
	OutTradeNo string `xml:"out_trade_no" json:"out_trade_no"`
	BuyerLogonId string `xml:"buyer_logon_id" json:"buyer_logon_id"`
	TotalAmount string `xml:"total_amount" json:"total_amount"`
	TransCurrency string `xml:"trans_currency" json:"trans_currency"`
}



func (tra *AliTrade) Handle(conf map[string]interface{}) (interface{}, error) {
	err := tra.BuildData(conf)
	if err != nil {
		return nil, err
	}
	ret, err := tra.sendReq(ALITRADE,tra)
	fmt.Println("===========")
	fmt.Println(string(ret))
	return tra.RetData(ret)
}

func (tra *AliTrade) RetData(ret []byte) (re *AliPayTradePayResponse, err error) {

	result := new(AliPayTradePayResponse)
	xml.Unmarshal(ret, result)
	return result, nil

}

func sign(m url.Values) string {
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
	b,_ := helper.RsaEncrypt([]byte(src))

	fmt.Println("加密：",b)
	fmt.Println("base加密：",base64.StdEncoding.EncodeToString(b))
	return base64.StdEncoding.EncodeToString(b)
}

func (tra *AliTrade) sendReq(reqUrl string, pay interface{}) (b []byte, err error) {

	client := helper.NewHttpClient()
	var data = url.Values{}

	data.Add("app_id", tra.AppId)
	data.Add("method", tra.Method)
	data.Add("charset", tra.Charset)
	data.Add("sign_type", "RSA2")
	data.Add("timestamp", tra.TimeStamp)
	data.Add("version", "1.0")
	data.Add("biz_content",tra.BizContent)
	data.Add("sign", sign(data))

	httpResp, err := client.PostForm(reqUrl,data)

	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	b, err = ioutil.ReadAll(httpResp.Body)

	return

}

func signWithPKCS1v15(param url.Values, privateKey []byte, hash crypto.Hash) (s string, err error) {
	if param == nil {
		param = make(url.Values, 0)
	}

	var pList = make([]string, 0, 0)
	for key := range param {
		var value = strings.TrimSpace(param.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	sort.Strings(pList)
	var src = strings.Join(pList, "&")
	fmt.Println(string(src))
	sig, err := encoding.SignPKCS1v15([]byte(src), privateKey, hash)
	if err != nil {
		return "", err
	}
	s = base64.StdEncoding.EncodeToString(sig)
	return s, nil
}

func (tra *AliTrade) BuildData(conf map[string]interface{}) error {

	b, _ := json.Marshal(conf)

	countInfo := string(b)

	tra.BaseAliConfig.BizContent = countInfo

	return nil
}

func (tra *AliTrade) InitBaseConfig(config *BaseAliConfig) {

	config.TimeStamp = time.Now().Format("2006-01-02 15:04:05")

	tra.BaseAliConfig = config
}


