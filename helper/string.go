package helper

import (
	"crypto/md5"
	"encoding/hex"
	rand1 "github.com/chanxuehong/rand"
	"time"
	"fmt"
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/rand"
	"crypto"
	"github.com/smartwalle/alipay/encoding"
)

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	cipherStr := h.Sum(nil)

	return hex.EncodeToString(cipherStr)
}

func NonceStr() string {
	return string(rand1.NewHex())
}

//CurrentTimeStampMS get current time with millisecond
func CurrentTimeStampMS() int64 {
	return time.Now().UnixNano() / time.Millisecond.Nanoseconds()
}

//rsa2加密
func RsaEncrypt(origData []byte,privateKey string) ([]byte, error) {
	//私钥切片处理
	key :=encoding.ParsePrivateKey(privateKey)

	block, _ := pem.Decode([]byte(key))//PiravteKeyData为私钥文件的字节数组
	if block == nil {
		fmt.Println("block空")
		return nil,nil
	}
	//priv即私钥对象,block2.Bytes是私钥的字节流
	var priv *rsa.PrivateKey

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		fmt.Println("无法还原私钥")
		return nil, nil
	}
	h2 := sha256.New()
	h2.Write(origData)
	hashed := h2.Sum(nil)
	signature2, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed) //签名
	return signature2,err
}

