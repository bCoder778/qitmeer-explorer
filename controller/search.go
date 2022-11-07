package controller

import (
	"fmt"
	"github.com/Qitmeer/qng-core/core/address"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
)

var InvalidAddr = fmt.Errorf("invalid address")

func (c *Controller) SearchV2(value string) (interface{}, error) {
	value = PkAddressToAddress(value)
	if len(value) == 35 {
		if !CheckAddress(value, c.conf.Qitmeer.Network) {
			return &types.SearchResult{
				Type:  "address",
				Value: value,
			}, InvalidAddr
		} else {
			return &types.SearchResult{
				Type:  "address",
				Value: value,
			}, nil
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

func PkAddressToAddress(addr string) string {
	if len(addr) == 53 {
		address, err := address.DecodeAddress(addr)
		if err != nil {
			return ""
		}
		return address.Encode()
	} else {
		return addr
	}
}
