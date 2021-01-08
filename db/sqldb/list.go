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
	transfers := []types.Transfer{}
	err := d.engine.Table(new(types.Transfer)).Select("tx_id,`timestamp`").Where("address = ?", address).Desc("timestamp").
		Limit(size, start).Find(&transfers)
	for _, trans := range transfers {
		txIds = append(txIds, trans.TxId)
	}
	return txIds, err
}

//select address, sum(Amount) as sumamount from %s WHERE spent_tx_id = '' and stat < %d GROUP BY address ORDER BY sumamount Desc Limit %d, %d
func (d *DB) BalanceTop(page, size int) ([]*dbtype.Address, error) {
	page -= 1
	start := page * size

	addrs := []*dbtype.Address{}
	err := d.engine.Table(new(types.Vout)).
		Select("address, sum(amount) as balance").
		Where("spent_tx = ? and unconfirmed_spent_tx = ? and stat in (?,?)", "", "", stat.TX_Confirmed, stat.TX_Unconfirmed).
		GroupBy("address").Desc("balance").Limit(size, start).Find(&addrs)

	return addrs, err
}
