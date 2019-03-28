package enter

import (
	"encoding/json"
	"github.com/ttooch/payment/helper"
)

type ModifyEnter struct {
	*ModifyConf
	BaseCharge
}

type ModifyConf struct {
	AccountNo    string     `json:"accountNo"`    //商户登陆账号
	FeeInfo      []*FeeInfo `json:"feeInfo"`      //费率信息
	FeeStartTime string     `json:"feeStartTime"` //费率生效开始时间
	FeeEndTime   string     `json:"feeEndTime"`   //费率生效结束时间
}

type ModifyReturn struct {
	ResultCode  string `json:"resultCode"`
	ErrorCode   string `json:"errorCode"`
	ErrCodeDesc string `json:"errCodeDesc"`
	FeeInfo     string `json:"feeInfo"` //费率信息
	Status      string `json:"status"`
	Ext         string `json:"ext"`
}

func (en *ModifyEnter) Handle(conf *ModifyConf) (interface{}, error) {
	err := en.BuildData(conf)
	if err != nil {
		return nil, err
	}
	en.SetSign()
	ret, err := en.SendReq(EnterUrl, en)
	if err != nil {
		return nil, err
	}
	return en.RetData(ret)
}

func (en *ModifyEnter) RetData(ret []byte) (*ModifyReturn, error) {

	ret, err := en.BaseCharge.RetData(ret)

	if err != nil {
		return nil, err
	}

	modifyReturn := new(ModifyReturn)

	err = json.Unmarshal(ret, &modifyReturn)

	if err != nil {
		return nil, err
	}

	return modifyReturn, nil
}

func (en *ModifyEnter) BuildData(conf *ModifyConf) error {

	b, err := json.Marshal(conf)

	if err != nil {
		return err
	}

	en.ModifyConf = conf

	en.ServiceName = "merchant.enter.modify"

	encryptData, err := helper.Rsa1Encrypt(en.PfxData, b, en.CertPassWord)

	if err != nil {
		return err
	}

	en.EncryptData = encryptData

	return nil
}
