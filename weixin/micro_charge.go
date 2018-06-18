/**
 * @author dengmeiyu
 * @since 20180608
 */
package weixin

import (
	"encoding/xml"
	"fmt"
	"encoding/json"
	"time"
)

type MicroCharge struct {
	XMLName xml.Name `xml:"xml" json:"-"`
	*PayConf
	BaseCharge
}

const SYSTEMERROR = "SYSTEMERROR"
const BANKERROR = "BANKERROR"
const USERPAYING = "USERPAYING"

type PayConf struct {
	Openid         string    `xml:"openid,omitempty" json:"openid,omitempty"`
	SubOpenid      string    `xml:"sub_openid,omitempty" json:"sub_openid,omitempty"`
	Body string `xml:"body" json:"body"`
	Detail         string    `xml:"detail,omitempty" json:"detail,omitempty"`
	Attach         string    `xml:"attach,omitempty" json:"attach,omitempty"`
	OutTradeNo     string `xml:"out_trade_no" json:"out_trade_no"`
	FeeType        string `xml:"fee_type,omitempty" json:"fee_type,omitempty"`
	TotalFee       int64  `xml:"total_fee" json:"total_fee"`
	SpbillCreateIp string `xml:"spbill_create_ip" json:"spbill_create_ip"`
	TimeStart      string    `xml:"time_start,omitempty" json:"time_start,omitempty"`
	TimeExpire     string    `xml:"time_expire,omitempty" json:"time_expire,omitempty"`
	GoodsTag       string    `xml:"goods_tag,omitempty" json:"goods_tag,omitempty" `
	TradeType      string    `xml:"trade_type" json:"trade_type"`
	DeviceInfo string `xml:"device_info,omitempty" json:"device_info,omitempty"` //店铺编号
	AuthCode   string `xml:"auth_code" json:"auth_code"`     //授权码
	SceneInfo      SceneInfo `xml:"scene_info,omitempty" json:"scene_info,omitempty"`
}

type SceneInfo struct {
	Id       string `xml:"id" json:"id"`
	Name     string `xml:"name" json:"name"`
	AreaCode string `xml:"areaCode" json:"areaCode"`
	Address  string `xml:"address" json:"address"`
}

type MicroReturn struct {
	SubAppid  string `json:"sub_appid"`
	TimeStamp int64  `json:"time_stamp"`
	NonceStr  string `json:"nonce_str"`
	//Package   string `json:"package"`
	SignType string `json:"sign_type"`
	PaySign  string `json:"paySign"`

	Openid             string `json:"openid"`
	SubOpenid          string `json:"sub_openid"`
	IsSubscribe        string `json:"is_subscribe"`
	SubIsSubscribe     string `json:"sub_is_subscribe"`
	TradeType          string `json:"trade_type"`
	BankType           string `json:"bank_type"`
	FeeType            string `json:"fee_type"`
	TotalFee           string `json:"total_fee"`
	CashFeeType        string `json:"cash_fee_type"`
	CashFee            string `json:"cash_fee"`
	SettlementTotalFee string `json:"settlement_total_fee"`
	CouponFee          string `json:"coupon_fee"`
	TransactionId      string `json:"transaction_id"`
	OutTradeNo         string `json:"out_trade_no"`
	TimeEnd            string `json:"time_end"`
}

type PayData map[string]interface{}

type Result struct {
	Error
	Return
	MicroReturn
}

func (app *MicroCharge) RetData(ret []byte) (Result, error) {

	result := Result{}
	xml.Unmarshal(ret, &result)
	return result, nil

}

func (app *MicroCharge) BuildData(conf map[string]interface{}) error {

	b, _ := json.Marshal(conf)

	var PayConf PayConf

	json.Unmarshal(b, &PayConf)

	if PayConf.SpbillCreateIp == "" {
		PayConf.SpbillCreateIp = "127.0.0.1"
	}

	if PayConf.FeeType == "" {
		PayConf.FeeType = "CNY"
	}

	app.PayConf = &PayConf

	return nil
}

func (app *MicroCharge) MicroPayWithPosConnectTimeout(conf map[string]interface{}) (Result, error) {

	err := app.BuildData(conf)

	if err != nil {
		return Result{}, err
	}
	app.SetSign(app)
	ret := app.SendReq(MicropayReqUrl, app)
	fmt.Println(string(ret))
	return app.RetData(ret)
}

/**
*
* 提交刷卡支付，并且确认结果，接口比较慢
* @param WxPayMicroPay $microPayInput
* @throws WxpayException
* @return 返回查询接口的结果
*/
func (app *MicroCharge) Handle(data map[string]interface{}) (interface{}, error) {
	//①、提交被扫支付
	result, err := app.MicroPayWithPosConnectTimeout(data)
	if err != nil {
		return MicroReturn{}, err
	}
	//如果返回成功
	if result.ReturnCode == "" || result.ResultCode == "" {
		fmt.Println("接口调用失败,请确认是否输入是否有误！")
	}

	//②、接口调用成功，明确返回调用失败
	if result.ReturnCode == "SUCCESS" && result.ResultCode == "FAIL" &&
		result.ErrCode != "USERPAYING" && result.ErrCode != "SYSTEMERROR" {
		return false, nil
	}

	//③、确认支付是否成功
	queryTimes := 10
	appid := app.AppId
	subAppId := app.SubAppId
	mchId := app.MchId
	subMchId := app.SubMchId
	signType := app.SignType
	md5Key := app.Md5Key

	for {
		queryTimes --
		succResult := 0
		que := new(OrderQuery)
		que.InitBaseConfig(&BaseConfig{
			AppId:    appid,
			SubAppId: subAppId,
			MchId:    mchId,
			SubMchId: subMchId,
			SignType: signType,
			Md5Key:   md5Key,
		})
		queryResult, succResult := que.Query(map[string]interface{}{"out_trade_no": app.OutTradeNo,
		}, succResult)
		if (succResult == 2) {
			time.Sleep(2 * time.Second)
			continue
		} else if (succResult == 1) { //查询成功
			return queryResult, nil
		} else { //订单交易失败
			return false, nil
		}
	}

	//④、确认失败，则撤销订单
	rev := new(OrderReverse)
	rev.InitBaseConfig(&BaseConfig{
		AppId:    appid,
		SubAppId: subAppId,
		MchId:    mchId,
		SubMchId: subMchId,
		SignType: signType,
		Md5Key:   md5Key,
	})

	if !rev.Cancel(map[string]interface{}{"out_trade_no": app.OutTradeNo,
	}, 0) {
		fmt.Println("撤销单失败！")
	}

	return false, nil
}
