/**
 * @author dengmeiyu
 * @since 20180608
 */
package weixin

import (
	"encoding/xml"
	"github.com/ttooch/payment/helper"
	"time"
	"fmt"
	"log"
	"encoding/json"
)

type MicroCharge struct {
	XMLName xml.Name `xml:"xml" json:"-"`
	*MicroConf
	BaseCharge
}

const SYSTEMERROR = "SYSTEMERROR"
const BANKERROR = "BANKERROR"
const USERPAYING = "USERPAYING"

type MicroConf struct {
	//Openid         string    `xml:"openid,omitempty" json:"openid,omitempty"`
	//SubOpenid      string    `xml:"sub_openid,omitempty" json:"sub_openid,omitempty"`
	Body           string    `xml:"body" json:"body"`
	Detail         string    `xml:"detail,omitempty" json:"detail,omitempty"`
	Attach         string    `xml:"attach,omitempty" json:"attach,omitempty"`
	OutTradeNo     string    `xml:"out_trade_no" json:"out_trade_no"`
	FeeType        string    `xml:"fee_type,omitempty" json:"fee_type,omitempty"`
	TotalFee       int64     `xml:"total_fee" json:"total_fee"`
	SpbillCreateIp string    `xml:"spbill_create_ip" json:"spbill_create_ip"`
	TimeStart      string    `xml:"time_start,omitempty" json:"time_start,omitempty"`
	TimeExpire     string    `xml:"time_expire,omitempty" json:"time_expire,omitempty"`
	GoodsTag       string    `xml:"goods_tag,omitempty" json:"goods_tag,omitempty" `
	TradeType      string    `xml:"trade_type" json:"trade_type"`
	DeviceInfo     string    `xml:"device_info" json:"device_info"` //店铺编号
	AuthCode       string    `xml:"auth_code" json:"auth_code"`     //授权码
	SceneInfo      SceneInfo `xml:"scene_info,omitempty" json:"scene_info,omitempty"`
}

type SceneInfo struct {
	Id       string `xml:"id" json:"id"`
	Name     string `xml:"name" json:"name"`
	AreaCode string `xml:"areaCode" json:"areaCode"`
	Address  string `xml:"address" json:"address"`
}

type MicroReturn struct {
	AppId     string `json:"appId"`
	SubAppid  string `json:"sub_appid"`
	TimeStamp int64  `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	//Package   string `json:"package"`
	SignType string `json:"signType"`
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

	var MicroConf MicroConf

	json.Unmarshal(b, &MicroConf)

	if MicroConf.SpbillCreateIp == "" {
		MicroConf.SpbillCreateIp = "127.0.0.1"
	}

	if MicroConf.FeeType == "" {
		MicroConf.FeeType = "CNY"
	}

	app.MicroConf = &MicroConf


	return nil
}

func (app *MicroCharge) Handle(conf map[string]interface{}) (Result, error) {

	err := app.BuildData(conf)

	if err != nil {
		return Result{}, err
	}
	app.SetSign(app)
	ret := app.SendReq(MicropayReqUrl,app)
	fmt.Println(string(ret))
	return app.RetData(ret)
}

/*
* 提交刷卡支付，针对软POS，尽可能做成功
*
* 内置重试机制，最多60s
*/
func (app *MicroCharge) MicroPayWithPosConnectTimeout(data map[string]interface{}, connectTimeoutMs int) (Result, error) {

	remainingTimeMs := 60 * 1000
	var err error
	var result Result
	for {
		startTimestampMs := helper.CurrentTimeStampMS()
		readTimeoutMs := remainingTimeMs - connectTimeoutMs
		if readTimeoutMs > 1000 {
			result, err := app.Handle(data)
			if err != nil {
				goto ERROR
			}
			if result.ReturnCode == SUCCESS {
				errCode := result.ErrCode
				if resultCode := result.ResultCode; resultCode == SUCCESS {
					break
				}
				// 查看错误码，若支付结果未知，则重试提交刷卡支付
				if errCode == SYSTEMERROR || errCode == BANKERROR || errCode == USERPAYING {
					remainingTimeMs = remainingTimeMs - (int)(helper.CurrentTimeStampMS()-startTimestampMs)
					if remainingTimeMs <= 100 {
						break
					}
					log.Println("microPayWithPos: try micropay again")
					if remainingTimeMs > 5*1000 {
						time.Sleep(5 * time.Second)
					} else {
						time.Sleep(time.Second)
					}
					continue

				} else {
					break
				}

			} else {
				break
			}

		} else {
			break
		}
	}
	//return MicroReturn
	return result, nil
ERROR:
	result = Result{}
	return result, err

}


/**
*
* 提交刷卡支付，并且确认结果，接口比较慢
* @param WxPayMicroPay $microPayInput
* @throws WxpayException
* @return 返回查询接口的结果
*/
func (app *MicroCharge) pay(data map[string]interface{}) (interface{}, error) {
	//①、提交被扫支付
	result, err := app.Handle(data)
	if err != nil {
		return MicroReturn{}, err
	}
	//如果返回成功
	if result.ReturnCode == "" || result.OutTradeNo == "" || result.ResultCode == "" {
		fmt.Println("接口调用失败,请确认是否输入是否有误！")
	}

	//签名验证
	//	outTradeNo := result.OutTradeNo

	//②、接口调用成功，明确返回调用失败
	if result.ReturnCode == "SUCCESS" && result.ResultCode == "FAIL" &&
		result.ErrCode != "USERPAYING" && result.ErrCode != "SYSTEMERROR" {
		return false, nil
	}

	//③、确认支付是否成功
	queryTimes := 10

	for {
		queryTimes --
		succResult := 0
		queryResult, succResult := app.query(data,succResult)
		if (succResult == 2) {
			time.Sleep(2 * time.Second)
			continue
		} else if (succResult == 1) { //查询成功
			return queryResult, nil
		} else { //订单交易失败
			return false, nil
		}
	}

	//④、次确认失败，则撤销订单
	if !app.cancel(data,0) {
		fmt.Println("撤销单失败！")
	}

	return false,nil
}

/**
*
* 撤销订单，如果失败会重复调用10次
* @param string $out_trade_no
* @param 调用深度 $depth
*/
func (app *MicroCharge) cancel(data map[string]interface{}, depth int) bool {
	if depth > 10 {
		return false
	}

	result,_ := app.reverse(data)

	//接口调用失败
	if (result.ReturnCode != "SUCCESS") {
		return false
	}

	//如果结果为success且不需要重新调用撤销，则表示撤销成功
	if (result.ReturnCode != "SUCCESS" && result.Recall == "N") {
		return true
	} else if result.Recall == "Y" {
		depth ++
		return app.cancel(data, depth)
	}
	return false
}

/**
	 *
	 * 查询订单情况
	 * @param string $out_trade_no  商户订单号
	 * @param int $succCode         查询订单结果
	 * @return 0 订单不成功，1表示订单成功，2表示继续等待
	 */
func (app *MicroCharge) query(data map[string]interface{}, succCode int) (interface{}, int) {

	result, _ := app.orderQuery(data)

	if result.ReturnCode == "SUCCESS" && result.ResultCode == "SUCCESS" {
		//支付成功
		if result.TradeState == "SUCCESS" {
			succCode = 1
			return result, succCode
		} else if result.TradeState == "USERPAYING" { //用户支付中

			succCode = 2
			return false, succCode
		}
	}

	//如果返回错误码为“此交易订单号不存在”则直接认定失败
	if result.ErrCode == "ORDERNOTEXIST" {
		succCode = 0
	} else {
		//如果是系统错误，则后续继续
		succCode = 2
	}
	return false, succCode
}

func (app *MicroCharge) reverse(data map[string]interface{}) (Result, error) {
	err := app.BuildData(data)

	if err != nil {
		return Result{}, err
	}
	app.SetSign(app)
	ret := app.SendReq(ReverseReqUrl,app)
	return app.RetData(ret)

}
func (app *MicroCharge) orderQuery(data map[string]interface{}) (Result, error) {
	err := app.BuildData(data)

	if err != nil {
		return Result{}, err
	}
	app.SetSign(app)
	ret := app.SendReq(OrderqueryReqUrl,app)
	return app.RetData(ret)
}
