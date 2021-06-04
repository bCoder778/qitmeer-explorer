package controller

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	"time"
)

func (c *Controller) QueryTokenTxs(page, size int, coinId, stat string) (*types.ListResp, error) {
	key := fmt.Sprintf("%s_%d_%d_%s", coinId, page, size, stat)
	value, err := c.cache.Value("TokenTxs", key)
	if err != nil {
		list, err := c.queryTokenTxs(page, size, coinId, stat)
		if err != nil {
			return nil, err
		}
		c.cache.Add("TokenTxs", key, 2*60*time.Second, list)
		return list, nil
	}
	return value.(*types.ListResp), nil
}

func (c *Controller) queryTokenTxs(page, size int, coinId, stat string) (*types.ListResp, error) {
	vs, err := c.storage.QueryTokenTransaction(page, size, coinId, stat)
	if err != nil {
		return nil, err
	}
	count, err := c.storage.GetTokenTransactionCount(coinId, stat)
	if err != nil {
		return nil, err
	}
	height, err := c.storage.GetLastHeight()
	if err != nil {
		return nil, err
	}
	return &types.ListResp{
		Page:  page,
		Size:  size,
		List:  types.ToVoutListResp(vs, height),
		Count: count,
	}, nil
}
