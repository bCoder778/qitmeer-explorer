package controller

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	"time"
)

func (c *Controller) AlgorithmList() []*types.AlgorithmResp {
	value, err := c.cache.Value("AlgorithmList", "AlgorithmList")
	if err != nil {
		list := c.qitmeer.AlgorithmList()
		c.cache.Add("AlgorithmList", "AlgorithmList", 120*time.Second, list)
		return list
	}
	return value.([]*types.AlgorithmResp)
}

func (c *Controller) AlgorithmLine(algorithm string, sec int) *types.AlgorithmLineResp {
	key := fmt.Sprintf("%s-%s", algorithm, sec)
	value, err := c.cache.Value("AlgorithmLine", key)
	if err != nil {
		line := c.qitmeer.AlgorithmLine(algorithm, sec)
		c.cache.Add("AlgorithmLine", key, time.Second*time.Duration(sec/10), line)
		return line
	}
	return value.(*types.AlgorithmLineResp)
}

func (c *Controller) GetCoinIds() []string {
	key := "getCoinIds"
	value, err := c.cache.Value("getCoinIds", key)
	if err != nil {
		tokens := c.getCoinIds()
		c.cache.Add("getCoinIds", key, time.Second*60*10, tokens)
		return tokens
	}
	return value.([]string)
}

func (c *Controller) getCoinIds() []string {
	coinIds := []string{}
	tokens, _ := c.rpcClient.GetTokens()
	for _, token := range tokens {
		coinIds = append(coinIds, token.CoinName)
	}
	return coinIds
}
