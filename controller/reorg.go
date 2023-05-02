package controller

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	"github.com/bCoder778/qitmeer-sync/utils"
	"sort"
	"strings"
	"time"
)

type ReorgResp struct {
	Order         uint64 `json:"order"`
	OldHash       string `json:"oldHash"`
	NewHash       string `json:"newHash"`
	OldMiner      string `json:"oldMiner"`
	NewMiner      string `json:"newMiner"`
	Confirmations uint64 `json:"confirmations"`
	EVMHeight     uint64 `json:"evmHeight"`
}

func (c *Controller) QueryReorg(page, size int) *types.ListResp {
	key := fmt.Sprintf("QueryReorg-%d-%d", page, size)
	value, err := c.cache.Value("QueryReorg", key)
	if err != nil {
		reorgs := c.queryReorg(page, size)
		c.cache.Add("QueryReorg", key, 60*3*time.Second, reorgs)
		return reorgs
	}
	return value.(*types.ListResp)
}

func (c *Controller) queryReorg(page, size int) *types.ListResp {
	reorgs, err := c.storage.QueryReorg(page, size)
	if err != nil {
		return &types.ListResp{}
	}
	count := c.storage.GetReorgCount()
	reorgResps := []*ReorgResp{}
	for _, reorg := range reorgs {
		reorgInfo := c.storage.GetReorgInfo(reorg.Hash)
		oldMiner, _ := utils.PkAddressToAddress(reorg.Address)
		_, err := c.rpcClient.GetBlockByHash(reorg.Hash)
		if err != nil && (strings.Contains(err.Error(), "no node") || strings.Contains(err.Error(), "no block")) {
			if reorgInfo.Hash == "" {
				//reorgResps = append(reorgResps, &ReorgResp{
				//	Order:     0,
				//	OldHash:   reorg.Hash,
				//	NewHash:   "unknown",
				//	OldMiner:  oldMiner,
				//	NewMiner:  "unknown",
				//	EVMHeight: 0,
				//})
			} else {
				newBlock, _ := c.storage.GetBlockByOrder(reorgInfo.OldOrder)

				newMiner, _ := utils.PkAddressToAddress(newBlock.Address)
				reorgResps = append(reorgResps, &ReorgResp{
					Order:         newBlock.Order,
					OldHash:       reorg.Hash,
					NewHash:       newBlock.Hash,
					OldMiner:      oldMiner,
					NewMiner:      newMiner,
					Confirmations: reorg.Confirmations,
					EVMHeight:     newBlock.EvmHeight,
				})

			}
		}

	}

	sort.Slice(reorgResps, func(i, j int) bool {
		return reorgResps[i].Order > reorgResps[j].Order
	})

	return &types.ListResp{
		Page:  page,
		Size:  size,
		Count: count,
		List:  reorgResps,
	}
}

type OrderReorgResp struct {
	Order            uint64 `json:"order"`
	Hash             string `json:"hash"`
	Miner            string `json:"miner"`
	EVMHeight        uint64 `json:"evmHeight"`
	ReorgTime        string `json:"reorgTime"`
	Old              string `json:"old"`
	OldConfirmations uint64 `json:"oldConfirmations"`
	OldMiner         string `json:"oldMiner"`
}

func (c *Controller) QueryOrderReorg(page, size int) *types.ListResp {
	reorgs, err := c.storage.QueryOrderReorg(page, size)
	if err != nil {
		return &types.ListResp{}
	}
	count := c.storage.GetOrderReorgCount()
	reorgResps := []*OrderReorgResp{}
	for _, reorg := range reorgs {
		reorgResps = append(reorgResps, &OrderReorgResp{
			Order:            0,
			Hash:             "",
			Miner:            "",
			EVMHeight:        0,
			ReorgTime:        "",
			Old:              "",
			OldConfirmations: 0,
			OldMiner:         "",
		})

	}

	return &types.ListResp{
		Page:  page,
		Size:  size,
		Count: count,
		List:  reorgResps,
	}
}
