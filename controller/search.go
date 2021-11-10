package controller

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
)

func (c *Controller) SearchV2(value string) (interface{}, error) {

	if len(value) == 35 {
		if !CheckAddress(value, c.conf.Qitmeer.Network){
			return &types.SearchResult{
				Type:  "address",
				Value: value,
			}, fmt.Errorf("invalid address")
		}
	}

	block, _ := c.storage.GetBlock(value)
	if len(block.Hash) > 0 {
		return &types.SearchResult{
			Type:  "block",
			Value: value,
		}, nil
	}

	tx, _ := c.storage.GetTransactionByTxId(value)
	if len(tx) > 0 {
		return &types.SearchResult{
			Type:  "transaction",
			Value: value,
		}, nil
	}

	return nil, fmt.Errorf("SEARCH_NOT_FOUND")
}
