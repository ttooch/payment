/**
 * @author dengmeiyu
 * @since 20180711
 */
package alipay

import (
	"testing"
	"fmt"
)

func TestAliTrade(t *testing.T) {

	app := new(AliTrade)

	app.InitBaseConfig(&BaseAliConfig{
		AppId:    "2017122701286469",
		Method:   "alipay.trade.pay",
		AppAuthToken:"201801BBe10810a6828f40328bf24cabdc748X87",
		Charset:        "utf-8",
	})


	ret, err := app.Handle(map[string]interface{}{
		"body":           "123",
		"out_trade_no":   "2018061510564487166",
		"trans_currency": "CNY",
		"total_amount":   "0.01",
		"auth_code":      "289317329214168639",
		"scene":          "bar_code",
		"buyer_id":       "1",
		"subject":        "123",

	})

	fmt.Printf("%+v", ret)

	fmt.Println(err)
}

