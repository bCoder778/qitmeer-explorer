package controller

import (
	"github.com/bCoder778/qitmeer-explorer/controller/types"
)

func (c *Controller) LastBlocks(page, size int) (*types.ListResp, error) {
	blocks, err := c.db.LastBlocks(page, size)
	if err != nil {
		return nil, err
	}
	count, err := c.db.GetBlockCount()
	if err != nil {
		return nil, err
	}
	return &types.ListResp{
		Page:  page,
		Size:  size,
		List:  blocks,
		Count: count,
	}, nil
}

func (c *Controller) BlockDetail(hash string) (*types.BlockDetail, error) {
	blockHeader, err := c.db.GetBlock(hash)
	if err != nil {
		return nil, err
	}
	txDetails := []*types.TransactionDetail{}
	txs, err := c.db.QueryTransactionsByBlockHash(hash)
	for _, tx := range txs {
		tx, err := c.TransactionDetail(tx.TxId)
		if err != nil {
			return nil, err
		}
		txDetails = append(txDetails, tx)
	}
	return &types.BlockDetail{Header: types.DBBlockToBlock(blockHeader), Transactions: txDetails}, nil
}
