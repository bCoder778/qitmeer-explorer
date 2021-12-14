package external

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	getIpInfo = "http://ip-api.com/json/"
	getPrices = "https://api.jbex.com/openapi/quote/v1/ticker/24hr?symbol=MEERUSDT"
)

func GetIpInfo(ip string) (*IpInfo, error) {
	body, err := getBody("POST", getIpInfo+ip, nil)
	ipInfo := &IpInfo{}
	err = json.Unmarshal(body, ipInfo)
	if err != nil {
		return nil, err
	}
	return ipInfo, err
}

func GetPrice() (*JbexResp, error) {
	body, err := getBody("GET", getPrices, nil)
	if err != nil {
		return nil, err
	}

	resp := &JbexResp{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getBody(method string, url string, param map[string]interface{}) ([]byte, error) {
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(paramBytes)
	req, err := http.NewRequest(method, url, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 100}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
