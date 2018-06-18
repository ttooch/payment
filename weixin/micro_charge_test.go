/**
 * @author dengmeiyu
 * @since 20180615
 */

package weixin

import (
"fmt"
"testing"
	"time"
)

func TestMicroPayWithPosConnectTimeout(t *testing.T) {

	app := new(MicroCharge)

	app.InitBaseConfig(&BaseConfig{
		AppId:    "wxf06ac118ca3d9533",
		SubAppId: "wxa33cba2b69f869f3",
		MchId:    "1495589652",
		SubMchId: "1495746312",
		Md5Key:   "057177a8459352933f755c535b0ab0ef",
		SignType: "MD5",
		ExpireDuration: time.Second * 300,
	})

	ret, err := app.Handle(map[string]interface{}{
		"device_info":      "013467007045764",
		"body":             "image形象店-深圳腾大- QQ公仔",
		"out_trade_no":     "2018061511000063153",
		"fee_type":         "CNY",
		"total_fee":        1,
		"spbill_create_ip": "127.0.0.1",
		"auth_code":        "134628524482100179",
	})

	fmt.Printf("%+v", ret)

	fmt.Println(err)
}
