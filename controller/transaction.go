package controller

import (
	"fmt"
	types2 "github.com/Qitmeer/qitmeer/core/types"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	"github.com/bCoder778/qitmeer-sync/verify/stat"
	"time"
)

func (c *Controller) LastTransactions(page, size int) (*types.ListResp, error) {
	key := fmt.Sprintf("%s-%s", page, size)
	value, err := c.cache.Value("LastTransactions", key)
	if err != nil {
		list, err := c.lastTransactions(page, size)
		if err != nil {
			return nil, err
		}
		c.cache.Add("LastTransactions", key, 2*60*time.Second, list)
		return list, nil
	}
	return value.(*types.ListResp), nil
}

func (c *Controller) lastTransactions(page, size int) (*types.ListResp, error) {
	txs, err := c.storage.LastTransactions(page, size)
	if err != nil {
		return nil, err
	}
	count, err := c.storage.GetTransactionCount()
	if err != nil {
		return nil, err
	}
	return &types.ListResp{
		Page:  page,
		Size:  size,
		List:  types.ToTransactionRespList(txs),
		Count: count,
	}, nil
}

func (c *Controller) LastAddressTransactions(page, size int, address string) (*types.ListResp, error) {
	key := fmt.Sprintf("%s-%s-%s", page, size, address)
	value, err := c.cache.Value("LastAddressTransactions", key)
	if err != nil {
		list, err := c.lastAddressTransactions(page, size, address)
		if err != nil {
			return nil, err
		}
		c.cache.Add("LastAddressTransactions", key, 60*time.Second, list)
		return list, nil
	}
	return value.(*types.ListResp), nil
}

func (c *Controller) lastAddressTransactions(page, size int, address string) (*types.ListResp, error) {
	txIds, err := c.storage.LastAddressTxId(page, size, address)
	if err != nil {
		return nil, err
	}
	count, err := c.storage.GetAddressTransactionCount(address)
	if err != nil {
		return nil, err
	}
	txsDetail := []*types.TransactionDetailResp{}
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

func (c *Controller) TransactionDetail(txId string, address string) (*types.TransactionDetailResp, error) {
	key := fmt.Sprintf("%s-%s", txId, address)
	value, err := c.cache.Value("TransactionDetail", key)
	if err != nil {
		detail, err := c.transactionDetail(txId, address)
		if err != nil {
			return nil, err
		}
		c.cache.Add("TransactionDetail", key, 2*60*time.Second, detail)
		return detail, nil
	}
	return value.(*types.TransactionDetailResp), nil
}

func (c *Controller) transactionDetail(txId string, address string) (*types.TransactionDetailResp, error) {
	var header *types.TransactionResp = &types.TransactionResp{}
	vin := []*types.VinResp{}
	vout := []*types.VoutResp{}
	changeMap := map[string]*types.TransferChange{}
	feesMap := map[string]*types.Fees{}
	changeList := []*types.TransferChange{}
	feesList := []*types.Fees{}
	txs, err := c.storage.GetTransactionByTxId(txId)
	if err != nil {
		return nil, err
	}

	if len(txs) != 0 {
		header = types.ToTransactionResp(txs[0])
	}
	if len(txs) > 1 {
		for _, tx := range txs {
			if tx.Stat == stat.TX_Unconfirmed || tx.Stat == stat.TX_Confirmed {
				header = types.ToTransactionResp(tx)
				break
			}
		}
	}
	dbVin, err := c.storage.QueryTransactionVin(txId)
	if err != nil {
		return nil, err
	}
	for _, in := range dbVin {
		vin = append(vin, types.ToVinResp(in))
		if in.Address == address {
			change, ok := changeMap[in.CoinId]
			if ok {
				change.UTotalVin += in.Amount
			} else {
				changeMap[in.CoinId] = &types.TransferChange{
					UTotalVin: in.Amount,
				}
			}
		}

		fee, ok := feesMap[in.CoinId]
		if ok {
			fee.UTotalVin += in.Amount
		} else {
			feesMap[in.CoinId] = &types.Fees{
				UTotalVin: in.Amount,
			}
		}
	}
	dbVout, err := c.storage.QueryTransactionVout(txId)
	if err != nil {
		return nil, err
	}
	for _, out := range dbVout {
		vout = append(vout, types.ToVoutResp(out))
		if out.Address == address {
			change, ok := changeMap[out.CoinId]
			if ok {
				change.UTotalVout += out.Amount
			} else {
				changeMap[out.CoinId] = &types.TransferChange{
					UTotalVout: out.Amount,
				}
			}
		}

		fee, ok := feesMap[out.CoinId]
		if ok {
			fee.UTotalVout += out.Amount
		} else {
			feesMap[out.CoinId] = &types.Fees{
				UTotalVout: out.Amount,
			}
		}
	}
	if address != "" {
		for coinId, change := range changeMap {
			change.TotalVin = types2.Amount(change.UTotalVin).ToCoin()
			change.TotalVout = types2.Amount(change.UTotalVout).ToCoin()
			change.CoinId = coinId
			change.Change = types2.Amount(change.UTotalVin - change.UTotalVout).ToCoin()
			changeList = append(changeList, change)
		}
	}

	for coinId, fees := range feesMap {
		fees.CoinId = coinId
		if fees.UTotalVin > fees.UTotalVout {
			fees.TotalVin = types2.Amount(fees.UTotalVin).ToCoin()
			fees.TotalVout = types2.Amount(fees.UTotalVout).ToCoin()
			fees.Amount = types2.Amount(fees.UTotalVin - fees.UTotalVout).ToCoin()
		} else {
			fees.Amount = 0
		}
		feesList = append(feesList, fees)
	}
	header.Changes = changeList
	header.Fees = feesList
	return &types.TransactionDetailResp{Header: header, Vout: vout, Vin: vin}, nil
}
