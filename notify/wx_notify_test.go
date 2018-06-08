package notify

import (
	"errors"
	"fmt"
	"github.com/ttooch/payment/weixin"
	"testing"
)

func TestWxNotify(t *testing.T) {

	notify := new(WxNotify)

	notify.InitBaseConfig(&weixin.BaseConfig{
		AppId:    "wxf06ac118ca3d9533",
		MchId:    "1495589652",
		SubAppId: "wxa33cba2b69f869f3",
		SubMchId: "1495746312",
		Md5Key:   "057177a8459352933f755c535b0ab0ef",
		SignType: "MD5",
	})

	xml := "<xml><appid><![CDATA[wxf06ac118ca3d9533]]></appid><bank_type><![CDATA[CFT]]></bank_type><cash_fee><![CDATA[1]]></cash_fee><fee_type><![CDATA[CNY]]></fee_type><is_subscribe><![CDATA[N]]></is_subscribe><mch_id><![CDATA[1495589652]]></mch_id><nonce_str><![CDATA[a5ca845515b55b1164cfb3fe0095e943]]></nonce_str><openid><![CDATA[oJb7cwYJW-YC6ynUtvVLFm9UXfgs]]></openid><out_trade_no><![CDATA[2018060615095472882]]></out_trade_no><result_code><![CDATA[SUCCESS]]></result_code><return_code><![CDATA[SUCCESS]]></return_code><sign><![CDATA[A71D51B60BCAAEF4285977F4E4C629DA]]></sign><sub_appid><![CDATA[wxa33cba2b69f869f3]]></sub_appid><sub_is_subscribe><![CDATA[N]]></sub_is_subscribe><sub_mch_id><![CDATA[1495746312]]></sub_mch_id><sub_openid><![CDATA[oyA310LEnY_JW_-BDHVJguSpFyKQ]]></sub_openid><time_end><![CDATA[20180608165743]]></time_end><total_fee>1</total_fee><trade_type><![CDATA[JSAPI]]></trade_type><transaction_id><![CDATA[4200000123201806080945667886]]></transaction_id></xml>"

	result := notify.handle(xml, func(data *WxNotifyData) error {
		fmt.Println(data)
		return nil
		return errors.New("回调错误")
	})

	fmt.Println(result)
}
