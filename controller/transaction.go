package controller

import (
	types2 "github.com/Qitmeer/qitmeer/core/types"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	"github.com/bCoder778/qitmeer-sync/verify/stat"
)

func (c *Controller) LastTransactions(page, size int) (*types.ListResp, error) {
	txs, err := c.db.LastTransactions(page, size)
	if err != nil {
		return nil, err
	}
	count, err := c.db.GetTransactionCount()
	if err != nil {
		return nil, err
	}
	return &types.ListResp{
		Page:  page,
		Size:  size,
		List:  types.DBTransactionsToTransactions(txs),
		Count: count,
	}, nil
}

func (c *Controller) LastAddressTransactions(page, size int, address string) (*types.ListResp, error) {
	txIds, err := c.db.LastAddressTxId(page, size, address)
	if err != nil {
		return nil, err
	}
	count, err := c.db.GetAddressTransactionCount(address)
	if err != nil {
		return nil, err
	}
	txsDetail := []*types.TransactionDetail{}
	for _, txId := range txIds {
		tx, err := c.TransactionDetail(txId, address)
		if err != nil {
			return nil, err
		}
		txsDetail = append(txsDetail, tx)
	}

	return &types.ListResp{
		Page:  page,
		Size:  size,
		List:  txsDetail,
		Count: count,
	}, nil
}

func (c *Controller) TransactionDetail(txId string, address string) (*types.TransactionDetail, error) {
	var header *types.Transaction
	vin := []*types.Vinout{}
	vout := []*types.Vinout{}
	var totalVin, totalVout uint64
	txs, err := c.db.GetTransactionByTxId(txId)
	if err != nil {
		return nil, err
	}

	if len(txs) != 0 {
		header = types.DBTransactionToTransaction(txs[0])
	}
	if len(txs) > 1 {
		for _, tx := range txs {
			if tx.Stat == stat.TX_Unconfirmed || tx.Stat == stat.TX_Confirmed {
				header = types.DBTransactionToTransaction(tx)
				break
			}
		}
	}
	dbVin, err := c.db.QueryTransactionVin(txId)
	if err != nil {
		return nil, err
	}
	for _, in := range dbVin {
		vin = append(vin, types.DBVinoutToVinout(in))
		if in.Address == address {
			totalVin += in.Amount
		}
	}
	dbVout, err := c.db.QueryTransactionVout(txId)
	if err != nil {
		return nil, err
	}
	for _, out := range dbVout {
		vout = append(vout, types.DBVinoutToVinout(out))
		if out.Address == address {
			totalVin += out.Amount
		}
	}
	header.AddressChange = types2.Amount(totalVin).ToCoin() - types2.Amount(totalVout).ToCoin()
	return &types.TransactionDetail{Header: header, Vout: vout, Vin: vin}, nil
}
