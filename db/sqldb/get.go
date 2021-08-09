package sqldb

import (
	dbtypes "github.com/bCoder778/qitmeer-explorer/db/types"
	"github.com/bCoder778/qitmeer-sync/storage/types"
	"github.com/bCoder778/qitmeer-sync/verify/stat"
)

func (d *DB) GetTransaction(txId string, blockHash string) (*types.Transaction, error) {
	return nil, nil
}

func (d *DB) GetTransactionByTxId(txId string) ([]*types.Transaction, error) {
	txs := []*types.Transaction{}

	err := d.engine.Table(new(types.Transaction)).Where("tx_id = ? and duplicate != ?", txId, 1).Find(&txs)
	return txs, err
}

func (d *DB) GetVout(txId string, vout int) (*types.Vout, error) {
	vinout := &types.Vout{}
	_, err := d.engine.Where("tx_id = ? and number = ?", txId, vout).Get(vinout)
	return vinout, err
}

func (d *DB) GetLastOrder() (uint64, error) {
	var block = &types.Block{}
	_, err := d.engine.Table(new(types.Block)).Desc("order").Get(block)
	return block.Order, err
}

func (d *DB) GetLastHeight() (uint64, error) {
	var block = &types.Block{}
	_, err := d.engine.Table(new(types.Block)).Desc("height").Get(block)
	return block.Height, err
}

func (d *DB) GetLastUnconfirmedOrder() (uint64, error) {
	var block = &types.Block{}
	_, err := d.engine.Table(new(types.Block)).Where("stat = ?", stat.Block_Unconfirmed).OrderBy("`order`").Get(block)
	return block.Order, err
}

func (d *DB) GetAllUtxo() float64 {
	amount, _ := d.engine.Where("spent_tx = ? and stat = ?", "", stat.TX_Confirmed).Sum(new(types.Vout), "amount")
	return amount
}

func (d *DB) GetConfirmedBlockCount() int64 {
	count, _ := d.engine.Table(new(types.Block)).Where("stat = ?", stat.Block_Confirmed).Count()
	return count
}

func (d *DB) GetBlockCount(stat string) (int64, error) {

	sql := d.engine.Table(new(types.Block))
	if len(stat) > 0 {
		sql.Where("find_in_set(stat, ?)", stat)
	}

	return sql.Count()
}

func (d *DB) GetValidBlockCount() (int64, error) {
	return d.engine.Table(new(types.Block)).Where("stat in (?, ?)", stat.Block_Confirmed, stat.Block_Unconfirmed).Count()
}

func (d *DB) GetTransactionCount(stat string) (int64, error) {
	sql := d.engine.Table(new(types.Transaction))
	if len(stat) > 0 {
		sql.Where("find_in_set(stat, ?)", stat)
	}

	return sql.Count()
}

func (d *DB) GetAddressTransactionCount(address, coin string) (int64, error) {
	count, err := d.engine.Table(new(types.Transfer)).Where("address = ? and coin_id = ? and is_coinbase = ?", address, coin, 0).
		Or("address = ? and coin_id = ? and is_coinbase = ? and is_blue = ?", address, coin, 1, 1).Distinct("tx_id").Count()
	return count, err
}

func (d *DB) GetBlock(hash string) (*types.Block, error) {
	block := &types.Block{}
	_, err := d.engine.Table(block).Where("hash = ?", hash).Get(block)
	return block, err
}

func (d *DB) GetBlockByOrder(order uint64) (*types.Block, error) {
	block := &types.Block{}
	_, err := d.engine.Table(block).Where("`order` = ?", order).Get(block)
	return block, err
}

func (d *DB) GetLastBlock() (*types.Block, error) {
	var block = &types.Block{}
	_, err := d.engine.Table(new(types.Block)).Desc("order").Get(block)
	return block, err
}

func (d *DB) GetAddressCount(coin string) (int64, error) {
	return d.engine.Table(new(types.Vout)).
		Where("spent_tx = ?  and coin_id = ? and stat in (?, ?)", "", coin, stat.TX_Confirmed, stat.TX_Unconfirmed).
		Distinct("address").Count()
}

func (d *DB) GetUsableAmount(address string, coinId string, height uint64) (float64, error) {
	return d.engine.Table(new(types.Vout)).Where("vout.lock <= ?  and address = ? and coin_id = ? and spent_tx = ? and stat = ?",
		height, address, coinId, "", stat.TX_Confirmed).
		Sum(new(types.Vout), "amount")
}

func (d *DB) GetUnconfirmedAmount(address string, coinId string) (float64, error) {
	return d.engine.Table(new(types.Vout)).Where("amount != 0 and address = ? and coin_id = ? and spent_tx = ? and (is_coinbase = ? or (is_coinbase = ? and is_blue = ?)) and stat in (?, ?)",
		address, coinId, "", 0, 1, 1, stat.TX_Unconfirmed, stat.TX_Memry).
		Sum(new(types.Vout), "amount")
}

func (d *DB) GetLockedAmount(address string, coinId string, height uint64) (float64, error) {
	return d.engine.Table(new(types.Vout)).Where("vout.lock > ? and address = ? and coin_id = ? and spent_tx = ? and stat = ?",
		height, address, coinId, "", stat.TX_Confirmed).
		Sum(new(types.Vout), "amount")
}

func (d *DB) GetLastMinerBlock(address string) *types.Block {
	block := types.Block{}
	d.engine.Table(new(types.Block)).Where("address = ? and stat in (?, ?)", address, stat.Block_Confirmed, stat.Block_Unconfirmed).Desc("order").Get(&block)
	return &block
}

func (d *DB) GetLastAlgorithmBlock(algorithm string, edgeBits int) (*types.Block, error) {
	block := new(types.Block)
	_, err := d.engine.Where("pow_name = ? and edge_bits = ?", algorithm, edgeBits).Desc("order").Get(block)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (d *DB) GetPeer(address string) (*dbtypes.Peer, error) {
	peer := new(dbtypes.Peer)
	_, err := d.engine.Where("address = ?", address).Get(peer)
	if err != nil {
		return nil, err
	}
	return peer, nil
}

func (d *DB) GetTokenTransactionCount(coinId, stat string) (int64, error) {

	sql := d.engine.Table(new(types.Vout))

	if len(coinId) > 0 {
		sql.Where("coin_id = ?", coinId)
	}

	if len(stat) > 0 {
		sql.Where("find_in_set(stat, ?)", stat)
	}

	return sql.Count()
}

func (d *DB)QueryTransferCount()(int64, error){
	return d.engine.Table(new(types.Transaction)).Where("is_coinbase = ?", 0).Count()
}

func (d *DB)QueryCoinbaseCount()(int64, error){
	return d.engine.Table(new(types.Transaction)).Where("is_coinbase = ?", 1).Count()
}