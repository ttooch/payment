package refund

import (
	"fmt"
	"github.com/ttooch/payment/weixin"
	"testing"
	"github.com/ttooch/payment/helper"
	"git.coding.net/ttouch_/util"
)

func TestWxRefund_Handle(t *testing.T) {

	wxRefund := new(WxRefund)

	wxRefund.InitBaseConfig(&weixin.BaseConfig{
		AppId:     "wxf06ac118ca3d9533",
		MchId:     "1495589652",
		SubAppId:  "wxa33cba2b69f869f3",
		SubMchId:  "1495746312",
		Md5Key:    "057177a8459352933f755c535b0ab0ef",
		SignType:  "MD5",
		//NotifyUrl:WXREFUNDURL,
		//Cert:"/Users/zhouchao/go/src/github.com/ttooch/payment/refund/pem/apiclient_cert.pem",
		//Key:"/Users/zhouchao/go/src/github.com/ttooch/payment/refund/pem/apiclient_key.pem",
		Cert:util.WX_CERT,
		Key:util.WX_KEY,
	})

	ret, err := wxRefund.Handle(map[string]interface{}{
		"transaction_id":  "4200000114201806162872791672",
		"out_refund_no":   "2018061218342609197",
		"total_fee":       2,
		"refund_fee":      1,
		"refund_desc":     "测试退款",
	})

	fmt.Println(ret)
	fmt.Println(err)
}

func TestWxRefund_BuildData(t *testing.T) {

	cert := `
-----BEGIN CERTIFICATE-----
	MIIEaTCCA9KgAwIBAgIEAZlrdzANBgkqhkiG9w0BAQUFADCBijELMAkGA1UEBhMC
	Q04xEjAQBgNVBAgTCUd1YW5nZG9uZzERMA8GA1UEBxMIU2hlbnpoZW4xEDAOBgNV
	BAoTB1RlbmNlbnQxDDAKBgNVBAsTA1dYRzETMBEGA1UEAxMKTW1wYXltY2hDQTEf
	MB0GCSqGSIb3DQEJARYQbW1wYXltY2hAdGVuY2VudDAeFw0xNzEyMjgwOTAzMjVa
	Fw0yNzEyMjYwOTAzMjVaMIGYMQswCQYDVQQGEwJDTjESMBAGA1UECBMJR3Vhbmdk
	b25nMREwDwYDVQQHEwhTaGVuemhlbjEQMA4GA1UEChMHVGVuY2VudDEOMAwGA1UE
	CxMFTU1QYXkxLTArBgNVBAMUJOS4iua1t+mAmua4lOS/oeaBr+enkeaKgOaciemZ
	kOWFrOWPuDERMA8GA1UEBBMIODU4MjU1NzgwggEiMA0GCSqGSIb3DQEBAQUAA4IB
	DwAwggEKAoIBAQC344FNL050uckEjHndUz+Sq2mquYJVPGYZLtbFzdQe3uUgzhRP
	LyQPyrNqM1joWBaxSDUmBynhgYIG6T/kl5xFb8DEO6V3LXPMaGB1qWPJBVdyKgVt
	v1Rj86hUqhzUeQwCnw/WN5s0zfSER/ah8DzVP7n1dy1nGvJCVI0lirhA7x72ujV5
	+RKPLWiRB2RWRHhiDAZ7GxRVsMbY6yVByDLEbBgptOQiJM9kskQuNFisIOCq8LQM
	cntX03VNPOEmYLI4uj9j05f774tMPw16LtJUoYLPKSky02iyyMRiOgxfMagRE3rS
	3u7xbBrbJ2Y2TMyRJT9LoGy6Y+innbOXXoRZAgMBAAGjggFGMIIBQjAJBgNVHRME
	AjAAMCwGCWCGSAGG+EIBDQQfFh0iQ0VTLUNBIEdlbmVyYXRlIENlcnRpZmljYXRl
	IjAdBgNVHQ4EFgQUdQtA/i/NQmCI86fJ+k7Q9UUsmUowgb8GA1UdIwSBtzCBtIAU
	PgUm9iJitBVbiM1kfrDUYqflhnShgZCkgY0wgYoxCzAJBgNVBAYTAkNOMRIwEAYD
	VQQIEwlHdWFuZ2RvbmcxETAPBgNVBAcTCFNoZW56aGVuMRAwDgYDVQQKEwdUZW5j
	ZW50MQwwCgYDVQQLEwNXWEcxEzARBgNVBAMTCk1tcGF5bWNoQ0ExHzAdBgkqhkiG
	9w0BCQEWEG1tcGF5bWNoQHRlbmNlbnSCCQC7VJcrvADoVzAOBgNVHQ8BAf8EBAMC
	BsAwFgYDVR0lAQH/BAwwCgYIKwYBBQUHAwIwDQYJKoZIhvcNAQEFBQADgYEACEO6
	mB1FeERYJR4sUE48O5IEtT1/2UifTqRtHSHA32ZFtv4EYD9FqfMpRmxdu6DvB65b
	8M34dZjMVNzF0ErEX8CwY2/ocsZRMoeRaVxLFxhczbw6+jk8ZJecIz1t+u3E6cGm
	cn9MzPyKsjluS2NJRHE1z3qQ1b3eTcQOBLcv0Kg=
-----END CERTIFICATE-----
`

	key := `
-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQC344FNL050uckE
jHndUz+Sq2mquYJVPGYZLtbFzdQe3uUgzhRPLyQPyrNqM1joWBaxSDUmBynhgYIG
6T/kl5xFb8DEO6V3LXPMaGB1qWPJBVdyKgVtv1Rj86hUqhzUeQwCnw/WN5s0zfSE
R/ah8DzVP7n1dy1nGvJCVI0lirhA7x72ujV5+RKPLWiRB2RWRHhiDAZ7GxRVsMbY
6yVByDLEbBgptOQiJM9kskQuNFisIOCq8LQMcntX03VNPOEmYLI4uj9j05f774tM
Pw16LtJUoYLPKSky02iyyMRiOgxfMagRE3rS3u7xbBrbJ2Y2TMyRJT9LoGy6Y+in
nbOXXoRZAgMBAAECggEAK8tK6A70KGh2r1alT0icQ1n/eGFl4NbZmfXFwKYcdp2v
A/kQbStU803bHNfNvMf934rK3ZunJcWqKeszBEX3ObI7voBWD+QLSAszrdIBqcZT
5SSM1m5Ssts8o60/4HpTmew34VSs6g6CxV4+JqMIWeErcIJJldydBZ8NXnuDFjN1
KNhbjNoiEkstqE7p7AND9nf+cIg/XZJ3Wto+q2tiixYvkozouvi3DSMNH0QKua+3
B7EV6voUecriNQzHrlW6Pf6Ts2pH1DGtUMlNLnKd6oLNLSvXQOeWQHBLU30AOb20
wNkQdFSU4N5aGKXVEVYYV3MyqIYBk8umZZ3LctmhgQKBgQDm2RIYPDMsQgudBzwK
WeKfBFBorRTxO6TXmkiV3cvuV3zYyHM/Xx6DIoAh6CphZgeVF8h196aG+DRW3KIc
idgrBm7B8bpoZSnx6ieZl5YfOtwxI7rb7veH1yiT7cRVmwA4lLFrEN8Al6tYq9qL
25kgvt33KPdHWzPkHxnN9iLdDQKBgQDL7J+DX0VTqWvNDj8L3XYwEhaZVbi42oI1
nqhwl+mrD39os1LHZ4tRBH2Lmd22wOA0faLWF673EMWhZK49zYTYOeRT0WS2Xxyz
dFpn0UeJjX2Mt3xmgv4CDgk7BXJfZeI020mAxsWm4WBEXcBcAhR66ya0uF2pvOue
bWqOGVepfQKBgQDIVVT4pzWZkO9XJLIcIYkjirDlRl9IQgR5/rBDLqoNeusjjkFF
m5x1HPmpcztkLxEhd/RFO9fOhZOWVf5xWca4/+oacVbdxf0yrYwUJJLpI0F3e5Rr
zOQPhM74IX0i9VKpx5B4Y7nDX9wZJXqRqM8otbsXBPwjdqh/reXM8W+r+QKBgQDB
hn7fGuPdALS4dmOoq/RErHXb8LkMXPojTr2FlPWBjISaZUwSxxY0vfEzMcNyc1qT
FgQZ74HxIG6dusGND2SaG16vmNFeLac8Oxis27RrOubCS0N4uam7Y3ypEYM4O6VQ
CieYWYsr00kbuGkcKDEtccpayXjB4/Mrd6Ue07gYtQKBgQCEJqXiuu0LgtwQbyol
SOranzp+hC8G8/FGuHIcEYZIFeQZhodDv3nN35+mnxF5nwk63VBKjNRCo5qF2HcB
JtEtiNqt2Dvfqv69Vkl1ojDbAziKcWYuQwh2slV7X9Jq90PqbygZPx4qNGyPZuEe
VXa6QBDids/MsLyun4vzdPQyJA==
-----END PRIVATE KEY-----
`

	helper.NewTLSBlockHttpClient([]byte(cert),[]byte(key))

}
