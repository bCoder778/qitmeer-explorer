package sqldb

import (
	"github.com/bCoder778/qitmeer-sync/storage/types"
	"github.com/bCoder778/qitmeer-sync/verify/stat"
)

func (d *DB) GetTransaction(txId string, blockHash string) (*types.Transaction, error) {
	return nil, nil
}

func (d *DB) GetTransactionByTxId(txId string) ([]*types.Transaction, error) {
	txs := []*types.Transaction{}
	err := d.engine.Table(new(types.Transaction)).Where("tx_id = ?", txId).Find(&txs)
	return txs, err
}

func (d *DB) GetVout(txId string, vout int) (*types.Vinout, error) {
	vinout := &types.Vinout{}
	_, err := d.engine.Where("tx_id = ? and type = ? and number = ?", txId, stat.TX_Vout, vout).Get(vinout)
	return vinout, err
}

func (d *DB) GetLastOrder() (uint64, error) {
	var block = &types.Block{}
	_, err := d.engine.Table(new(types.Block)).Desc("order").Get(block)
	return block.Order, err
}

func (d *DB) GetLastUnconfirmedOrder() (uint64, error) {
	var block = &types.Block{}
	_, err := d.engine.Table(new(types.Block)).Where("stat = ?", stat.Block_Unconfirmed).OrderBy("`order`").Get(block)
	return block.Order, err
}

func (d *DB) GetAllUtxo() float64 {
	amount, _ := d.engine.Where("spent_tx = ? and type = ? and stat = ?", "", stat.TX_Vout, stat.TX_Confirmed).Sum(new(types.Vinout), "amount")
	return amount
}

func (d *DB) GetConfirmedBlockCount() int64 {
	count, _ := d.engine.Table(new(types.Block)).Where("stat = ?", stat.Block_Confirmed).Count()
	return count
}

func (d *DB) GetBlockCount() (int64, error) {
	return d.engine.Table(new(types.Block)).Count()
}

func (d *DB) GetTransactionCount() (int64, error) {
	return d.engine.Table(new(types.Transaction)).Count()
}

func (d *DB) GetBlock(hash string) (*types.Block, error) {
	block := &types.Block{}
	_, err := d.engine.Table(block).Where("hash = ?", hash).Get(block)
	return block, err
}
