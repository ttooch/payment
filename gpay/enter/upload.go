package enter

import (
	"bytes"
	"code.aliyun.com/udian_pay/payment/utils"
	"encoding/json"
	"github.com/ttooch/payment/helper"
	"io"
	"io/ioutil"
	"mime/multipart"
)

type Upload struct {
	*UploadConf
	BaseCharge
}

type UploadConf struct {
	ImageData []byte `json:"image_data"`
	ImageNo   string `json:"image_no"`
}

type UploadReturn struct {
	ResultCode  string `json:"resultCode"`
	ErrorCode   string `json:"errorCode"`
	ErrCodeDesc string `json:"errCodeDesc"`
	ImageNo     string `json:"imageNo"`
	ImagePath   string `json:"imagePath"`
	Ext         string `json:"ext"`
}

func (up *Upload) Handle(conf *UploadConf) (interface{}, error) {
	err := up.BuildData(conf)
	if err != nil {
		return nil, err
	}
	up.SetSign()
	ret, err := up.SendReq(UploadUrl)
	if err != nil {
		return nil, err
	}
	return up.RetData(ret)
}

func (up *Upload) RetData(ret []byte) (*UploadReturn, error) {

	var baseReturn BaseConfig

	err := json.Unmarshal(ret, &baseReturn)

	if err != nil {
		return nil, err
	}

	uploadReturn := new(UploadReturn)

	err = json.Unmarshal([]byte(baseReturn.EncryptData), &uploadReturn)

	if err != nil {
		return nil, err
	}

	return uploadReturn, nil
}

func (up *Upload) BuildData(conf *UploadConf) error {

	_, err := json.Marshal(conf)

	if err != nil {
		return err
	}

	up.UploadConf = conf

	up.ServiceName = "merchant.enter.upload"

	return nil
}

func (up *Upload) SetSign() {
	up.SignData = utils.Md5(up.AgentNo + "3.0" + up.ImageNo + up.Key)
}

func (up *Upload) SendReq(reqUrl string) (b []byte, err error) {
	client := helper.NewHttpClient()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("agentNo", up.AgentNo)
	writer.WriteField("serviceName", up.ServiceName)
	writer.WriteField("version", up.Version)
	writer.WriteField("imageNo", up.ImageNo)
	writer.WriteField("signData", up.SignData)

	//关键的一步操作
	part, err := writer.CreateFormFile("imageData", "file.png")

	if err != nil {
		return nil, err
	}

	//iocopy
	_, err = io.Copy(part, bytes.NewReader(up.ImageData))

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	contentType := writer.FormDataContentType()

	rsp, err := client.Post(reqUrl, contentType, body)

	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	b, err = ioutil.ReadAll(rsp.Body)

	if err != nil {
		return nil, err
	}

	return b, nil

}
