package refund

import (
	"fmt"
	"github.com/ttooch/payment/weixin"
	"testing"
)

func TestWxRefund_Handle(t *testing.T) {

	wxRefund := new(WxRefund)

	wxRefund.InitBaseConfig(&weixin.BaseConfig{
		AppId:     "wxf06ac118ca3d9533",
		MchId:     "1495589652",
		SubAppId:  "wxa33cba2b69f869f3",
		SubMchId:  "1495746312",
		Md5Key:    "057177a8459352933f755c535b0ab0ef",
		SignType:  "MD5",
		NotifyUrl: "http://api.store.udian.me/v1/payment/CallBackWxPay",
	})

	ret, err := wxRefund.Handle(map[string]interface{}{
		"transaction_id":  "4200000114201806162872791672",
		"out_refund_no":   "2018061218342609197",
		"total_fee":       2,
		"refund_fee":      1,
		"refund_desc":     "测试退款",
	})

	fmt.Println(ret)
	fmt.Println(err)
}
