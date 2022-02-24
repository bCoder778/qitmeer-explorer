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
	// 除去不被认可的order = 0的块
	err := d.engine.Table(new(types.Block)).Where("block.order != 0 or color != 2").Desc("order").Limit(size, start).Find(&blocks)
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
func (d *DB) LastAddressTxId(page, size int, address, coin string) ([]string, error) {
	page -= 1
	start := page * size
	txIds := []string{}
	transfers := []types.Transfer{}
	err := d.engine.Table(new(types.Transfer)).Select("DISTINCT(tx_id),`timestamp`").Where("address = ? and coin_id = ? and is_coinbase = ?", address, coin, 0).
		Or("address = ? and coin_id = ? and is_coinbase = ? and is_blue = ?", address, coin, 1, 1).Desc("timestamp").
		Limit(size, start).Find(&transfers)
	for _, trans := range transfers {
		txIds = append(txIds, trans.TxId)
	}
	return txIds, err
}

//select address, sum(Amount) as sumamount from %s WHERE spent_tx_id = '' and stat < %d GROUP BY address ORDER BY sumamount Desc Limit %d, %d
func (d *DB) BalanceTop(page, size int, coinId string) ([]*dbtype.Address, error) {
	page -= 1
	start := page * size

	addrs := []*dbtype.Address{}
	err := d.engine.Table(new(types.Vout)).
		Select("address, coin_id, sum(amount) as balance").
		Where("spent_tx = ? and coin_id = ? and (is_coinbase = 0 or (is_coinbase = 1 and is_blue = 1)) and stat in (?,?)", "", coinId, stat.TX_Confirmed, stat.TX_Unconfirmed).
		GroupBy("address").Desc("balance").Limit(size, start).Find(&addrs)

	return addrs, err
}

func (d *DB) QueryBlock(page, size int, stat string) ([]*types.Block, error) {
	page -= 1
	start := page * size

	var bs []*types.Block

	sql := d.engine.Table(new(types.Block)).Where("1 = 1")
	if len(stat) > 0 {
		sql.Where("stat in (?)", stat).And("block.order != 0 or color != 2")
	}

	err := sql.Desc("order").Limit(size, start).Find(&bs)
	return bs, err
}

func (d *DB) QueryTransaction(page, size int, stat string) ([]*types.Transaction, error) {
	page -= 1
	start := page * size

	var txs []*types.Transaction
	sql := d.engine.Table(new(types.Transaction))

	if len(stat) > 0 {
		sql.Where("stat in (?)", stat)
	}

	err := sql.Desc("id").Limit(size, start).Find(&txs)
	return txs, err
}

func (d *DB) QueryTokenTransaction(page, size int, coinId, stat string) ([]*types.Vout, error) {
	page -= 1
	start := page * size

	var vos []*types.Vout

	sql := d.engine.Table(new(types.Vout))

	if len(coinId) > 0 {
		sql.Where("coin_id = ?", coinId)
	}

	if len(stat) > 0 {
		sql.Where("find_in_set(stat, ?)", stat)
	}
	err := sql.Desc("id").Limit(size, start).Find(&vos)

	return vos, err
}

func (d *DB) QueryTransfer(page, size int) ([]*types.Transaction, error) {
	page -= 1
	start := page * size
	var txs []*types.Transaction
	err := d.engine.Table(new(types.Transaction)).Where("is_coinbase = ? and duplicate = ?", 0, 0).Desc("id").Limit(size, start).Find(&txs)
	return txs, err
}

func (d *DB) QueryCoinbase(page, size int) ([]*types.Transaction, error) {
	page -= 1
	start := page * size
	var txs []*types.Transaction
	err := d.engine.Table(new(types.Transaction)).Where("is_coinbase = ? and duplicate = ?", 1, 0).Desc("timestamp").Limit(size, start).Find(&txs)
	return txs, err
}
