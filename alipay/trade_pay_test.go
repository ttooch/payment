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
		AppId:    "2088922549974876",
		Method:   "alipay.trade.pay",
		AppAuthToken:"201801BB635d7c6ab47a459ba579f876d3736X68",
		Charset:        "utf-8",
	})


	ret, err := app.Handle(map[string]interface{}{
		"body":           "123",
		"out_trade_no":   "2018061510564487165",
		"trans_currency": "CNY",
		"total_amount":   "0.01",
		"auth_code":      "283986513883525558",
		"scene":          "bar_code",
		"buyer_id":       "1",
		"subject":        "123",

	})

	fmt.Println("33333333333333333333333333")
	fmt.Printf("%+v", ret)

	fmt.Println(err)
}

