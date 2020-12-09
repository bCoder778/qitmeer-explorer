package controller

import (
	"github.com/bCoder778/qitmeer-explorer/controller/types"
)

func (c *Controller) LastBlocks(page, size int) (*types.ListResp, error) {
	blocks, err := c.storage.LastBlocks(page, size)
	if err != nil {
		return nil, err
	}
	count, err := c.storage.GetBlockCount()
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

func (c *Controller) BlockDetail(hash string) (*types.BlockDetailResp, error) {
	blockHeader, err := c.storage.GetBlock(hash)
	if err != nil {
		return nil, err
	}
	txDetails := []*types.TransactionDetailResp{}
	txs, err := c.storage.QueryTransactionsByBlockHash(hash)
	for _, tx := range txs {
		tx, err := c.TransactionDetail(tx.TxId, "no address")
		if err != nil {
			return nil, err
		}
		txDetails = append(txDetails, tx)
	}
	return &types.BlockDetailResp{Header: types.ToBlockResp(blockHeader), Transactions: txDetails}, nil
}
