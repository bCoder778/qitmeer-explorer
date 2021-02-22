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
	err := d.engine.Table(new(types.Transaction)).Where("tx_id = ?", txId).Find(&txs)
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

func (d *DB) GetBlockCount() (int64, error) {
	return d.engine.Table(new(types.Block)).Count()
}

func (d *DB) GetValidBlockCount() (int64, error) {
	return d.engine.Table(new(types.Block)).Where("stat in (?, ?)", stat.Block_Confirmed, stat.Block_Unconfirmed).Count()
}

func (d *DB) GetTransactionCount() (int64, error) {
	return d.engine.Table(new(types.Transaction)).Count()
}

func (d *DB) GetAddressTransactionCount(address string) (int64, error) {
	return d.engine.Table(new(types.Transfer)).Where("address = ? ", address).Count()
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

func (d *DB) GetAddressCount() (int64, error) {
	return d.engine.Table(new(types.Vout)).
		Where("spent_tx = ? and stat in (?, ?)", "", stat.TX_Confirmed, stat.TX_Unconfirmed).
		GroupBy("address").Count()
}

func (d *DB) GetUsableAmount(address string) (float64, error) {
	return d.engine.Table(new(types.Vout)).Where("address = ? and spent_tx = ? and stat = ?",
		address, "", stat.TX_Confirmed).
		Sum(new(types.Vout), "amount")
}

func (d *DB) GetLockedAmount(address string) (float64, error) {
	return d.engine.Table(new(types.Vout)).Where("address = ? and spent_tx = ? and stat = ?",
		address, "", stat.TX_Unconfirmed).
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
