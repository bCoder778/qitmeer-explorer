package sqldb

import (
	"github.com/bCoder778/qitmeer-sync/storage/types"
	"github.com/bCoder778/qitmeer-sync/verify/stat"
)

func (d *DB) QueryUnConfirmedOrders() ([]uint64, error) {
	orders := []uint64{}
	err := d.engine.Table(new(types.Block)).Where("stat = ?", stat.Block_Unconfirmed).Cols("order").Find(&orders)
	return orders, err
}

func (d *DB) QueryUnconfirmedTranslateTransaction() ([]types.Transaction, error) {
	txs := []types.Transaction{}
	err := d.engine.Where("is_coinbase = ?", 0).And("stat = ? or stat = ?", stat.TX_Unconfirmed, stat.TX_Memry).Find(&txs)
	return txs, err
}

func (d *DB) QueryMemTransaction() ([]types.Transaction, error) {
	txs := []types.Transaction{}
	err := d.engine.Where("stat = ?", stat.TX_Memry).Find(&txs)
	return txs, err
}

func (d *DB) QueryTransactions(txId string) ([]types.Transaction, error) {
	txs := []types.Transaction{}
	err := d.engine.Where("tx_id = ?", txId).Find(&txs)
	return txs, err
}

func (d *DB) QueryTransactionsByBlockHash(hash string) ([]types.Transaction, error) {
	txs := []types.Transaction{}
	err := d.engine.Where("block_hash = ?", hash).Find(&txs)
	return txs, err
}

func (d *DB) QueryTransactionVin(txId string) ([]*types.Vinout, error) {
	txs := []*types.Vinout{}
	err := d.engine.Where("tx_id = ? and type = ?", txId, stat.TX_Vin).Find(&txs)
	return txs, err
}

func (d *DB) QueryTransactionVout(txId string) ([]*types.Vinout, error) {
	txs := []*types.Vinout{}
	err := d.engine.Where("tx_id = ? and type = ?", txId, stat.TX_Vout).Find(&txs)
	return txs, err
}
