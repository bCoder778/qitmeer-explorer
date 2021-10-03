package controller

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	types2 "github.com/bCoder778/qitmeer-sync/storage/types"
	"strconv"
	"time"
)

func (c *Controller) LastBlocks(page, size int) (*types.ListResp, error) {
	key := fmt.Sprintf("%d-%d", page, size)
	value, err := c.cache.Value("LastBlocks", key)
	if err != nil {
		blockList, err := c.lastBlocks(page, size)
		if err != nil {
			return nil, err
		}
		c.cache.Add("LastBlocks", key, 2*time.Second, blockList)
		return blockList, nil
	} else {
		return value.(*types.ListResp), nil
	}
}

func (c *Controller) lastBlocks(page, size int) (*types.ListResp, error) {
	blocks, err := c.storage.LastBlocks(page, size)
	if err != nil {
		return nil, err
	}
	count, err := c.storage.GetBlockCount("0,1,3,4")
	if err != nil {
		return nil, err
	}
	return &types.ListResp{
		Page:  page,
		Size:  size,
		List:  types.ToBlockRespList(blocks),
		Count: count,
	}, nil
}

func (c *Controller) BlockDetail(condition string) (*types.BlockDetailResp, error) {
	value, err := c.cache.Value("BlockDetail", condition)
	if err != nil {
		detail, err := c.blockDetail(condition)
		if err != nil {
			return nil, err
		}
		c.cache.Add("BlockDetail", condition, 10*time.Second, detail)
		return detail, nil
	} else {
		return value.(*types.BlockDetailResp), nil
	}
}

func (c *Controller) blockDetail(condition string) (*types.BlockDetailResp, error) {
	order, err := strconv.ParseUint(condition, 10, 64)
	var blockHeader *types2.Block
	if err == nil {
		blockHeader, err = c.storage.GetBlockByOrder(order)
		if err != nil {
			return nil, err
		}
	} else {
		blockHeader, err = c.storage.GetBlock(condition)
		if err != nil {
			return nil, err
		}
	}

	var txDetails []*types.TransactionDetailResp
	if blockHeader.Transactions <= 10 {
		txs, _ := c.storage.QueryTransactionsByBlockHash(blockHeader.Hash, 100, 1)
		for _, tx := range txs {
			tx, err := c.TransactionDetail(tx.TxId, blockHeader.Hash, "no address")
			if err != nil {
				return nil, err
			}
			txDetails = append(txDetails, tx)
		}
	}

	return &types.BlockDetailResp{Header: types.ToBlockResp(blockHeader), Transactions: txDetails}, nil
}

func (c *Controller) BlockTransactions(hash string, size, p int) (*types.ListResp, error) {

	var txs []*types.TransactionDetailResp
	rs, _ := c.storage.QueryTransactionsByBlockHash(hash, size, p)
	for _, tx := range rs {
		tx, err := c.TransactionDetail(tx.TxId, hash, "no address")
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return &types.ListResp{List: txs, Size: size, Page: p, Count: 0}, nil
}

func (c *Controller) QueryBlockByStatus(page, size int, stat string) (*types.ListResp, error) {
	key := fmt.Sprintf("%d-%d-%s", page, size, stat)
	value, err := c.cache.Value("LastBlocks", key)
	if err != nil {
		blockList, err := c.queryBlockByStatus(page, size, stat)
		if err != nil {
			return nil, err
		}
		c.cache.Add("LastBlocks", key, 30*time.Second, blockList)
		return blockList, nil
	}
	return value.(*types.ListResp), nil

}

func (c *Controller) queryBlockByStatus(page, size int, stat string) (*types.ListResp, error) {
	blocks, err := c.storage.QueryBlock(page, size, stat)
	if err != nil {
		return nil, err
	}
	count, err := c.storage.GetBlockCount(stat)
	if err != nil {
		return nil, err
	}
	return &types.ListResp{
		Page:  page,
		Size:  size,
		List:  types.ToBlockRespList(blocks),
		Count: count,
	}, nil
}
