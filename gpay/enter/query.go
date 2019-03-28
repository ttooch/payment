package enter

import (
	"encoding/json"
	"github.com/ttooch/payment/helper"
)

type QueryEnter struct {
	*QueryConf
	BaseCharge
}

type QueryConf struct {
	AccountNo string `json:"accountNo"` //商户登陆账号
}

type QueryReturn struct {
	AccountNo    string `json:"accountNo"`
	Status       string `json:"status"`
	AccountType  string `json:"accountType"`
	MerchantName string `json:"merchantName"`
	MerchantMail string `json:"merchantMail"`
	PhoneNo      string `json:"phoneNo"`
	SettleInfo   string `json:"settleInfo"`
	FeeInfo      string `json:"feeInfo"`
	FeeStartTime string `json:"feeStartTime"`
	FeeEndTime   string `json:"feeEndTime"`
}

func (en *QueryEnter) Handle(conf *QueryConf) (interface{}, error) {
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

func (en *QueryEnter) RetData(ret []byte) (*QueryReturn, error) {

	ret, err := en.BaseCharge.RetData(ret)

	if err != nil {
		return nil, err
	}

	modifyReturn := new(QueryReturn)

	err = json.Unmarshal(ret, &modifyReturn)

	if err != nil {
		return nil, err
	}

	return modifyReturn, nil
}

func (en *QueryEnter) BuildData(conf *QueryConf) error {

	b, err := json.Marshal(conf)

	if err != nil {
		return err
	}

	en.QueryConf = conf

	en.ServiceName = "merchant.enter.query"

	encryptData, err := helper.Rsa1Encrypt(en.PfxData, b, en.CertPassWord)

	if err != nil {
		return err
	}

	en.EncryptData = encryptData

	return nil
}
