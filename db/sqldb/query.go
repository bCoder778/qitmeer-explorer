package sqldb

import (
	dbtypes "github.com/bCoder778/qitmeer-explorer/db/types"
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

func (d *DB) QueryTransactionsByBlockHash(hash string, size, p int) ([]types.Transaction, error) {
	var txs []types.Transaction
	start := (p - 1) * size
	err := d.engine.Where("block_hash = ?", hash).Limit(size, start).Find(&txs)
	return txs, err
}

func (d *DB) QueryTransactionVin(txId string) ([]*types.Vin, error) {
	txs := []*types.Vin{}
	err := d.engine.Where("tx_id = ?", txId).Find(&txs)
	return txs, err
}

func (d *DB) QueryTransactionVout(txId string) ([]*types.Vout, error) {
	txs := []*types.Vout{}
	err := d.engine.Where("tx_id = ?", txId).Find(&txs)
	return txs, err
}

func (d *DB) QueryAlgorithmDiffInTime(algorithm string, edgeBits int, max int64, min int64) []*types.Block {
	blocks := []*types.Block{}
	d.engine.Table(new(types.Block)).Where("pow_name = ? and edge_bits = ? and timestamp between ? and ?", algorithm, edgeBits, min, max).Find(&blocks)
	return blocks
}

func (d *DB) QueryPeers() []*dbtypes.Peer {
	peers := []*dbtypes.Peer{}
	d.engine.Table(new(dbtypes.Peer)).Find(&peers)
	return peers
}

func (d *DB) QueryTokens() []string {
	tokens := []string{}
	d.engine.Table(new(types.Vout)).Distinct("coinid").Find(&tokens)
	return tokens
}

func (d *DB) QueryLocation() []*dbtypes.Location {
	var locals []*dbtypes.Location
	_ = d.engine.Table(new(dbtypes.Location)).Find(&locals)
	return locals
}
