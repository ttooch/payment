/**
 * @author dengmeiyu
 * @since 20180618
 */
package weixin

import (
	"encoding/xml"
	"encoding/json"
)

type OrderQuery struct {
	XMLName xml.Name `xml:"xml" json:"-"`
	BaseCharge
	*OrderOutTradeNo
}

type QueryReturn struct {
	//DeviceInfo string `xml:"device_info" json:"device_info"`
	OpenId         string `xml:"open_id" json:"open_id"`
	SubIsSubScribe string `xml:"sub_is_sub_scribe" json:"sub_is_sub_scribe"`
	TradeType      string `xml:"trade_type" json:"trade_type"`
	BankType       string `xml:"bank_type" json:"bank_type"`
	TotalFee       int    `xml:"total_fee" json:"total_fee"`
	CashFee        int    `xml:"cash_fee" json:"cash_fee"`
	OutTradeNo     string `xml:"out_trade_no" json:"out_trade_no"`
	TransactionId  string `xml:"transaction_id" json:"transaction_id"`
	TimeEnd        string `xml:"time_end" json:"time_end"`
	TradeStateDesc string `xml:"trade_state_desc" json:"trade_state_desc"`
}

type QueryResult struct {
	Error
	Return
	QueryReturn
}

/**
 *
 * 查询订单情况
 * @param string $out_trade_no  商户订单号
 * @param int $succCode         查询订单结果
 * @return 0 订单不成功，1表示订单成功，2表示继续等待
 */
func (que *OrderQuery) Query(data map[string]interface{}, succCode int) (interface{}, int) {

	queryResult, _ := que.OrderQuery(data)

	if queryResult.ReturnCode == "SUCCESS" && queryResult.ResultCode == "SUCCESS" {
		//支付成功
		if queryResult.TradeState == "SUCCESS" {
			succCode = 1
			return queryResult, succCode
		} else if queryResult.TradeState == "USERPAYING" { //用户支付中

			succCode = 2
			return false, succCode
		}
	}

	//如果返回错误码为“此交易订单号不存在”则直接认定失败
	if queryResult.ErrCode == "ORDERNOTEXIST" {
		succCode = 0
	} else {
		//如果是系统错误，则后续继续
		succCode = 2
	}
	return false, succCode
}

func (que *OrderQuery) OrderQuery(data map[string]interface{}) (QueryResult, error) {
	err := que.BuildData(data)
	if err != nil {
		return QueryResult{}, err
	}
	que.SetSign(que)
	ret := que.SendReq(OrderQueryReqUrl, que)
	return que.RetData(ret)
}

func (que *OrderQuery) RetData(ret []byte) (QueryResult, error) {

	result := QueryResult{}
	xml.Unmarshal(ret, &result)
	return result, nil

}

func (que *OrderQuery) BuildData(conf map[string]interface{}) error {

	var OrderOutTradeNo OrderOutTradeNo
	b, _ := json.Marshal(conf)
	json.Unmarshal(b, &OrderOutTradeNo)

	que.OrderOutTradeNo = &OrderOutTradeNo

	return nil
}
