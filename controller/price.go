package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bCoder778/qitmeer-sync/params"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	getIpInfo = "http://ip-api.com/json/"
	//getPrices = "https://api.hotbit.io/api/v1/market.status?market=PMEER/USDT&period=1800"
	getPrices = "https://apiv4.upex.io/exchange-open-api/open/api/get_ticker?symbol=pmeerusdt"
)

type Price struct {
	NowPrice string `json:"nowPrice"`
	UpDown   string `json:"updown"`
}

type PriceApiResult struct {
	Id     int    `json:"id"`
	Result *Price `json:"result"`
	Error  string `json:"error"`
}

type UpexPriceApiResult struct {
	Code string     `json:"code"`
	Msg  string     `json:"msg"`
	Data *UpexPrice `json:"data"`
}

type UpexPrice struct {
	High float64 `json:"high"`
	Vol  float64 `json:"vol"`
	Last float64 `json:"last"`
	Low  float64 `json:"low"`
	Buy  string  `json:"buy"`
	Sell string  `json:"sell"`
	Rose float64 `json:"rose"`
	Time uint64  `json:"time"`
}

func (c *Controller) GetPrice() (*Price, error) {
	body, err := getBody("GET", getPrices, nil)
	prices := &UpexPriceApiResult{}
	err = json.Unmarshal(body, prices)
	if err != nil {
		return nil, err
	}
	if prices.Data != nil {
		return &Price{
			NowPrice: fmt.Sprintf("%.4f", prices.Data.Last),
			UpDown:   fmt.Sprintf("%.4f", prices.Data.Rose*100),
		}, nil
	}
	return nil, errors.New("no PMEER")
}

func (c *Controller) GetCirculating() string {
	value, err := c.cache.Value("GetCirculating", "GetCirculating")
	if err != nil {
		circulating := c.getCirculating()
		c.cache.Add("GetCirculating", "GetCirculating", 10*60*time.Second, circulating)
		return circulating
	}
	return value.(string)
}

func (c *Controller) getCirculating() string {
	count, err := c.storage.GetValidBlockCount()
	height, err := c.storage.GetLastHeight()
	if err != nil {
		return "0"
	}
	reward := params.Qitmeer10Params.BlockReward
	genesisUTXO := params.Qitmeer10Params.GenesisUTXO["MEER"]
	unlock := uint64(math.Floor(float64(height)/2880)) * 2640287564088

	reward = reward + unlock
	circulating := (uint64(count)-1)*reward + genesisUTXO

	sCirculating := strconv.FormatUint(circulating, 10)
	return sCirculating
}

func (c *Controller) GetCirculatingFloat() string {
	pMeer := c.GetCirculating()
	if len(pMeer) > 8 {
		return pMeer[0:len(pMeer)-8] + "." + pMeer[len(pMeer)-8:]
	} else {
		return "0"
	}
}

func (c *Controller) GetMaxPMeer() string {
	return "20028781000000000"
}

func (c *Controller) GetMaxFloatPMeer() string {
	return "200287810"
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
