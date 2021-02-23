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
	getPrices = "https://www.ubcoin.pro/api/ubcoin-recharge/runDown/getAllRunDown"
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

	client := &http.Client{Timeout: time.Second * 20}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
