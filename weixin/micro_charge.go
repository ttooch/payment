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
	*PayConf
	BaseCharge
}

const SYSTEMERROR = "SYSTEMERROR"
const BANKERROR = "BANKERROR"
const USERPAYING = "USERPAYING"

type PayConf struct {
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
	ret := app.SendReq(MicropayReqUrl,app)
	fmt.Println(string(ret))
	return app.RetData(ret)
}

/*
* 提交刷卡支付，针对软POS，尽可能做成功
*
* 内置重试机制，最多60s
*/
func (app *MicroCharge) Handle(data map[string]interface{}, connectTimeoutMs int) (Result, error) {

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
