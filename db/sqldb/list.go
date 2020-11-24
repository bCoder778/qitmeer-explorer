package sqldb

import (
	dbtype "github.com/bCoder778/qitmeer-explorer/db/types"
	"github.com/bCoder778/qitmeer-sync/storage/types"
	"github.com/bCoder778/qitmeer-sync/verify/stat"
)

func (d *DB) LastBlocks(page, size int) ([]*types.Block, error) {
	page -= 1
	start := page * size

	blocks := []*types.Block{}
	err := d.engine.Table(new(types.Block)).Desc("order").Limit(size, start).Find(&blocks)
	return blocks, err
}

func (d *DB) LastTransactions(page, size int) ([]*types.Transaction, error) {
	page -= 1
	start := page * size

	txs := []*types.Transaction{}
	err := d.engine.Table(new(types.Transaction)).Desc("id").Limit(size, start).Find(&txs)
	return txs, err
}

//select address, sum(Amount) as sumamount from %s WHERE spent_tx_id = '' and stat < %d GROUP BY address ORDER BY sumamount Desc Limit %d, %d
func (d *DB) MaxBalanceAddress(page, size int) ([]*dbtype.Address, error) {
	page -= 1
	start := page * size

	addrs := []*dbtype.Address{}
	err := d.engine.Table(new(types.Vinout)).
		Select("address, sum(amount) as balance").
		Where("type = ? and spent_tx = ? and stat = ?", stat.TX_Vout, "", stat.TX_Confirmed).
		GroupBy("address").Desc("balance").Limit(size, start).Find(&addrs)

	return addrs, err
}
