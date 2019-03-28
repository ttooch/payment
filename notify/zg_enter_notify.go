package notify

import (
	"errors"
	"github.com/ttooch/payment/gpay/enter"
	"github.com/ttooch/payment/helper"
	"encoding/json"
	"fmt"
)

type ZgEnterNotifyData struct {
	AccountNo    string `json:"accountNo"`
	MerchantName string `json:"merchantName"`
	SmMd5Key     string `json:"smMd5Key"`
	SmMerchantNo string `json:"smMerchantNo"`
	Status       string `json:"status"`
}

type ZgEnterNotify struct {
	*enter.BaseConfig
	notifyData *ZgEnterNotifyData
}

func (zg *ZgEnterNotify) InitBaseConfig(config *enter.BaseConfig) {
	zg.BaseConfig = config
}

func (zg *ZgEnterNotify) getNotifyData(ret string) error {

	if ret == "" {
		return errors.New("获取通知数据失败")
	}

	err := json.Unmarshal([]byte(ret),zg)

	if err != nil {
		return errors.New("解析返回数据失败:" + err.Error())
	}

	err = zg.checkNotify()

	if err != nil {
		return err
	}

	retData, err := helper.Rsa1Decrypt(zg.PfxData, zg.EncryptData, zg.CertPassWord)

	if err != nil {
		return err
	}

	notify := new(ZgEnterNotifyData)

	err = json.Unmarshal(retData,notify)

	if err != nil {
		return errors.New("解析通知数据失败:" + err.Error())
	}

	zg.notifyData = notify

	return nil
}

func (zg *ZgEnterNotify) checkNotify() error {
	if zg.ResponseCode != "0000" {
		return errors.New("中钢返回错误" + zg.ResponseMsg)
	}
	return zg.verifySign()
}

func (zg *ZgEnterNotify) verifySign() error {
	return nil
}

func (zg *ZgEnterNotify) replyNotify(err error) bool {

	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return true
	}
}

type CallBack func(data *ZgEnterNotifyData) error

func (zg *ZgEnterNotify) Handle(ret string, callBack CallBack) bool {
	err := zg.getNotifyData(ret)

	if err != nil {
		return zg.replyNotify(err)
	}

	err = callBack(zg.notifyData)

	if err != nil {
		return zg.replyNotify(err)
	}

	return zg.replyNotify(nil)
}
