package helper

import (
	"net/http"
	"time"
	"crypto/tls"
	"net"
)

func NewHttpClient(timeout ...time.Duration) *http.Client {
	client := *http.DefaultClient

	if len(timeout) >0  {
		client.Timeout = timeout[0]
	}else {
		client.Timeout = time.Second * 5
	}
	return &client
}

func NewTLSHttpClient(certFile, keyFile string) (httpClient *http.Client) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   3 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       tlsConfig,
		},
		Timeout: 5 * time.Second,
	}
}

func NewTLSBlockHttpClient(certPEMBlock, keyPEMBlock  []byte) (httpClient *http.Client) {
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   3 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       tlsConfig,
		},
		Timeout: 5 * time.Second,
	}
}