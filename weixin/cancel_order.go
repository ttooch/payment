/**
 * @author dengmeiyu
 * @since 20180618
 */
package weixin

import (
	"encoding/xml"
	"encoding/json"
)

type OrderReverse struct {
	XMLName xml.Name `xml:"xml" json:"-"`
	BaseCharge
	*OrderOutTradeNo
}

type ReverseResult struct {
	Error
	Return
}

func (rev *OrderReverse) Reverse(data map[string]interface{}) (ReverseResult, error) {
	err := rev.BuildData(data)

	if err != nil {
		return ReverseResult{}, err
	}
	rev.SetSign(rev)
	ret := rev.SendReq(ReverseReqUrl, rev)
	return rev.RetData(ret)

}

/**
*
* 撤销订单，如果失败会重复调用10次
* @param string $out_trade_no
* @param 调用深度 $depth
*/
func (rev *OrderReverse) Cancel(data map[string]interface{}, depth int) bool {
	if depth > 10 {
		return false
	}

	result, _ := rev.Reverse(data)

	//接口调用失败
	if (result.ReturnCode != "SUCCESS") {
		return false
	}

	//如果结果为success且不需要重新调用撤销，则表示撤销成功
	if (result.ReturnCode != "SUCCESS" && result.Recall == "N") {
		return true
	} else if result.Recall == "Y" {
		depth ++
		return rev.Cancel(data, depth)
	}
	return false
}

func (que *OrderReverse) BuildData(conf map[string]interface{}) error {

	var OrderOutTradeNo OrderOutTradeNo
	b, _ := json.Marshal(conf)
	json.Unmarshal(b, &OrderOutTradeNo)

	que.OrderOutTradeNo = &OrderOutTradeNo

	return nil
}

func (que *OrderReverse) RetData(ret []byte) (ReverseResult, error) {

	result := ReverseResult{}
	xml.Unmarshal(ret, &result)
	return result, nil

}
