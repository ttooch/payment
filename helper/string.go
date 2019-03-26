package helper

import (
	"bytes"
	"crypto"
	"crypto/md5"
	rand2 "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	rand1 "github.com/chanxuehong/rand"
	"golang.org/x/crypto/pkcs12"
	"math/rand"
	"strings"
	"time"
)

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	cipherStr := h.Sum(nil)

	return hex.EncodeToString(cipherStr)
}

//生产orderSn
func CreateSn() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strings.Replace(time.Now().Format("060102150405.000000"), ".", "", -1) + fmt.Sprintf("%06v", rnd.Int31n(1000000))
}

func NonceStr() string {
	return string(rand1.NewHex())
}

//CurrentTimeStampMS get current time with millisecond
func CurrentTimeStampMS() int64 {
	return time.Now().UnixNano() / time.Millisecond.Nanoseconds()
}

//rsa1公钥加密
func Rsa1Encrypt(pfxData, origData []byte, certPassWord string) (string, error) {

	fmt.Println(pfxData)
	fmt.Println(certPassWord)

	_, cert, err := pkcs12.Decode(pfxData, certPassWord)

	if err != nil {
		return "",err
	}

	pubKey := cert.PublicKey

	var pub = pubKey.(*rsa.PublicKey)

	partLen := pub.N.BitLen()/8 - 11

	chunks := split([]byte(origData), partLen)

	buffer := bytes.NewBufferString("")

	for _, chunk := range chunks {

		bytes, err := rsa.EncryptPKCS1v15(rand2.Reader, pub, chunk)

		if err != nil {

			return "", err

		}

		buffer.Write(bytes)

	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil

}

//rsa1私钥解密
func Rsa1Decrypt(pfxData []byte,origData string,certPassWord string) ([]byte, error) {

	private, _, err := pkcs12.Decode(pfxData, certPassWord)

	if err != nil {
		return nil,err
	}

	var pri = private.(*rsa.PrivateKey)

	//分段加密
	partLen := pri.N.BitLen() / 8

	raw, err := base64.StdEncoding.DecodeString(origData)

	chunks := split([]byte(raw), partLen)

	buffer := bytes.NewBufferString("")

	for _, chunk := range chunks {

		decrypted, err := rsa.DecryptPKCS1v15(rand2.Reader, pri, chunk)

		if err != nil {
			return nil, err
		}

		buffer.Write(decrypted)

	}

	return buffer.Bytes(), err
}

//rsa2加密
func RsaEncrypt(origData []byte, privateKey string) ([]byte, error) {
	//私钥切片处理
	key := ParsePrivateKey(privateKey)

	block, _ := pem.Decode([]byte(key)) //PiravteKeyData为私钥文件的字节数组
	if block == nil {
		fmt.Println("block空")
		return nil, nil
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
	signature2, err := rsa.SignPKCS1v15(rand2.Reader, priv, crypto.SHA256, hashed) //签名
	return signature2, err
}

func RSAVerify(src []byte, sign, publicKey string) error {

	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	//支付宝公钥做切片处理
	key := ParsePublicKey(publicKey)
	block, _ := pem.Decode([]byte(key)) //PublicKeyData为私钥文件的字节数组
	if block == nil {
		fmt.Println("Public block空")
		return nil
	}

	var pubI interface{}
	pubI, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	h2 := sha256.New()
	h2.Write(src)
	hashed := h2.Sum(nil)
	var pub = pubI.(*rsa.PublicKey)
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed, signBytes) //验签

	if err != nil {
		return err
	}
	return nil
}

func split(buf []byte, lim int) [][]byte {

	var chunk []byte

	chunks := make([][]byte, 0, len(buf)/lim+1)

	for len(buf) >= lim {

		chunk, buf = buf[:lim], buf[lim:]

		chunks = append(chunks, chunk)

	}
	if len(buf) > 0 {

		chunks = append(chunks, buf[:len(buf)])

	}

	return chunks
}
