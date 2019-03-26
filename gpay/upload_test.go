package gpay

import (
	"testing"
	"io/ioutil"
	"fmt"
	"github.com/ttooch/payment/helper"
)

func TestUpload_Handle(t *testing.T) {
	up := new(Upload)

	up.InitBaseConfig(&BaseConfig{
		AgentNo:"zl@cs.sh.cn",
		Key:"ZD7D0ZM7SC8CVV6Q",
	})

	b,_ := ioutil.ReadFile("bank,jpg")

	ret, err := up.Handle(&UploadConf{
		ImageNo:helper.CreateSn(),
		ImageData:b,
	})

	fmt.Println(ret)
	fmt.Println(err)
}