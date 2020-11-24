package sqldb

import "github.com/bCoder778/qitmeer-sync/storage/types"

func (d *DB) LastBlocks(page, size int) ([]types.Block, error) {
	page -= 1
	start := page * size

	blocks := []types.Block{}
	err := d.engine.Table(new(types.Block)).Desc("order").Limit(size, start).Find(&blocks)
	return blocks, err
}

func (d *DB) LastTransactions(page, size int) ([]types.Transaction, error) {
	page -= 1
	start := page * size

	txs := []types.Transaction{}
	err := d.engine.Table(new(types.Transaction)).Desc("id").Limit(size, start).Find(&txs)
	return txs, err
}
