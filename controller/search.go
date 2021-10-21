package controller

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
)

func (c *Controller) SearchV2(value string) (interface{}, error) {

	if len(value) == 35 {
		return &types.SearchResult{
			Type:  "address",
			Value: value,
		}, nil
	}

	block, _ := c.storage.GetBlock(value)
	if block != nil {
		return &types.SearchResult{
			Type:  "block",
			Value: value,
		}, nil
	}

	tx, _ := c.storage.GetTransactionByTxId(value)
	if tx != nil {
		return &types.SearchResult{
			Type:  "transaction",
			Value: value,
		}, nil
	}

	return nil, fmt.Errorf("not found")
}
