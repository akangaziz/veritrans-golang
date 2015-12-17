package veritrans

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	SANDBOX_BASE_URL    = "https://api.sandbox.veritrans.co.id/v2"
	PRODUCTION_BASE_URL = "https://api.veritrans.co.id/v2"
)

type veritrans struct {
	ServerKey    string
	isProduction bool
	curlOptions  http.Header
	certificate  string
}

var vt *veritrans = nil

func New(key string, prod bool, pemFile string) *veritrans {
	if vt == nil {
		vt = new(veritrans)
		vt.ServerKey = key
		vt.isProduction = prod
		vt.certificate = pemFile
	}

	return vt
}

func (v *veritrans) RemoteCall(method string, url string, key string, payload interface{}) (*http.Response, error) {
	pem, err := ioutil.ReadAll(bytes.NewBufferString(v.certificate))
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(pem)

	tlsConfig := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}

	tlsConfig.BuildNameToCertificate()

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	client := &http.Client{
		Transport: transport,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBufferString(string(data)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s:", base64.StdEncoding.EncodeToString([]byte(key))))

	if len(vt.curlOptions) > 0 {
		for k, v := range vt.curlOptions {
			req.Header[k] = v
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func (v *veritrans) GetBaseUrl() (string, error) {
	if v == nil {
		return "", errors.New("Error, Have you initialize?")
	}

	if v.isProduction {
		return PRODUCTION_BASE_URL, nil
	}

	return SANDBOX_BASE_URL, nil
}

func (v *veritrans) Get(url, key string, payload interface{}) (*http.Response, error) {
	resp, err := v.RemoteCall("GET", url, key, payload)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (v *veritrans) Post(url, key string, payload interface{}) (*http.Response, error) {
	resp, err := v.RemoteCall("POST", url, key, payload)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (v *veritrans) VtWebCharge(payload interface{}) (*http.Response, error) {
	url, err := v.GetBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := v.Post(fmt.Sprintf("%s/charge", url), v.ServerKey, payload)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (v *veritrans) VtDirectCharge(payload interface{}) (*http.Response, error) {
	url, err := v.GetBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := v.Post(fmt.Sprintf("%s/charge", url), v.ServerKey, payload)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (v *veritrans) Status(id string) (*http.Response, error) {
	url, err := v.GetBaseUrl()
	if err != nil {
		return nil, err
	}

	resp, err := v.Get(fmt.Sprintf("%s/%s/status", url, id), v.ServerKey, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (v *veritrans) Approve(id string) (int, error) {
	url, err := v.GetBaseUrl()
	if err != nil {
		return http.StatusBadGateway, err
	}

	resp, err := v.Post(fmt.Sprintf("%s/%s/approve", url, id), v.ServerKey, nil)
	if err != nil {
		return http.StatusBadGateway, err
	}

	return resp.StatusCode, nil
}

func (v *veritrans) Cancel(id string) (int, error) {
	url, err := v.GetBaseUrl()
	if err != nil {
		return http.StatusBadGateway, err
	}

	resp, err := v.Post(fmt.Sprintf("%s/%s/cancel", url, id), v.ServerKey, nil)
	if err != nil {
		return http.StatusBadGateway, err
	}

	return resp.StatusCode, nil
}

func (v *veritrans) Expire(id string) (int, error) {
	url, err := v.GetBaseUrl()
	if err != nil {
		return http.StatusBadGateway, err
	}

	resp, err := v.Post(fmt.Sprintf("%s/%s/expire", url, id), v.ServerKey, nil)
	if err != nil {
		return http.StatusBadGateway, err
	}

	return resp.StatusCode, nil
}
