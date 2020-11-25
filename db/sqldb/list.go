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

//select td.* from %s td left join (select address, tx_id from %s where address = '%s' group by tx_id) ot on ot.tx_id = td.tx_id left join (select `from`, spent_tx_id from %s where `from` = '%s' group by spent_tx_id) tv on tv.spent_tx_id = td.tx_id where ot.address = '%s' or tv.`from` = '%s' order by td.id desc limit %d, %d;
func (d *DB) LastAddressTxId(page, size int, address string) ([]string, error) {
	page -= 1
	start := page * size
	txIds := []string{}
	vinouts := []types.Vinout{}
	err := d.engine.Table(new(types.Vinout)).Select("DISTINCT(tx_id),`id`").Desc("id").Where("address = ?", address).
		Limit(size, start).Find(&vinouts)
	for _, vinout := range vinouts {
		txIds = append(txIds, vinout.TxId)
	}
	return txIds, err
}

//select address, sum(Amount) as sumamount from %s WHERE spent_tx_id = '' and stat < %d GROUP BY address ORDER BY sumamount Desc Limit %d, %d
func (d *DB) BalanceTop(page, size int) ([]*dbtype.Address, error) {
	page -= 1
	start := page * size

	addrs := []*dbtype.Address{}
	err := d.engine.Table(new(types.Vinout)).
		Select("address, sum(amount) as balance").
		Where("type = ? and spent_tx = ? and stat = ?", stat.TX_Vout, "", stat.TX_Confirmed).
		GroupBy("address").Desc("balance").Limit(size, start).Find(&addrs)

	return addrs, err
}
