package types

import (
	qittypes "github.com/Qitmeer/qitmeer/core/types"
	types2 "github.com/bCoder778/qitmeer-explorer/db/types"
	"github.com/bCoder778/qitmeer-sync/storage/types"
	"github.com/bCoder778/qitmeer-sync/verify/stat"
)

type ListResp struct {
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Count int64       `json:"count"`
	List  interface{} `json:"list"`
}

type BlockDetail struct {
	Header       *Block               `json:"header"`
	Transactions []*TransactionDetail `json:"transactions"`
}

type TransactionDetail struct {
	Header *Transaction `json:"header"`
	Vin    []*Vinout    `json:"vin"`
	Vout   []*Vinout    `json:"vout"`
}

type Transaction struct {
	Id            uint64      `json:"id"`
	TxId          string      `json:"txid"`
	BlockHash     string      `json:"blockhash"`
	BlockOrder    uint64      `json:"blockorder"`
	TxHash        string      `json:"txhash"`
	Size          int         `json:"size"`
	Version       uint32      `json:"version"`
	Locktime      uint64      `json:"locktime"`
	Timestamp     int64       `json:"timestamp"`
	Expire        uint64      `json:"expire"`
	Confirmations uint64      `json:"confirmations"`
	Txsvaild      bool        `json:"txsvaild"`
	IsCoinbase    bool        `json:"iscoinbase"`
	Vins          int         `json:"vin"`
	Vouts         int         `json:"vout"`
	TotalVin      float64     `json:"totalvin"`
	TotalVout     float64     `json:"totalvout"`
	AddressChange float64     `json:"addresschange"`
	Fees          float64     `json:"fees"`
	Duplicate     bool        `json:"duplicate"`
	Stat          stat.TxStat `json:"stat"`
}

type Vinout struct {
	Id                     uint64              `json:"id"`
	TxId                   string              `json:"txid"`
	Type                   stat.TxType         `json:"type"`
	Number                 int                 `json:"number"`
	Order                  uint64              `json:"order"`
	Timestamp              int64               `json:"timestamp"`
	Address                string              `json:"address"`
	Amount                 float64             `json:"amount"`
	ScriptPubKey           *types.ScriptPubKey `json:"scriptpubkey"`
	SpentTx                string              `json:"spenttx"`
	SpentNumber            int                 `json:"spentnumber"`
	UnconfirmedSpentTx     string              `json:"unconfirmedspenttx"`
	UnconfirmedSpentNumber int                 `json:"unconfirmedspentnumber"`
	SpentedTx              string              `json:"spentedtx"`
	Vout                   int                 `json:"vout"`
	Sequence               uint64              `json:"sequence"`
	ScriptSig              *types.ScriptSig    `json:"scriptsig"`
	Stat                   stat.TxStat         `json:"stat"`
}

type Block struct {
	Id            uint64         `json:"id"`
	Hash          string         `json:"hash"`
	Txvalid       bool           `json:"txvalid"`
	Confirmations uint64         `json:"confirmation"`
	Version       uint32         `json:"version"`
	Weight        uint64         `json:"weight"`
	Height        uint64         `json:"height"`
	TxRoot        string         `json:"txroot"`
	Order         uint64         `json:"order"`
	Transactions  int            `json:"transactions"`
	StateRoot     string         `json:"stateroot"`
	Bits          string         `json:"bits"`
	Timestamp     int64          `json:"timestamp"`
	ParentRoot    string         `json:"parantroot"`
	Parents       []string       `json:"parents"`
	Children      []string       `json:"children"`
	Difficulty    uint64         `json:"difficulty"`
	PowName       string         `json:"powname"`
	PowType       int            `json:"powtype"`
	Nonce         uint64         `json:"nonce"`
	EdgeBits      int            `json:"edgebits"`
	CircleNonces  string         `json:"circlenonces"`
	Address       string         `json:"address"`
	Amount        float64        `json:"amount"`
	Miner         *MinerPool     `json:"miner"`
	Stat          stat.BlockStat `json:"stat"`
}

func DBTransactionToTransaction(tx *types.Transaction) *Transaction {
	return &Transaction{
		Id:            tx.Id,
		TxId:          tx.TxId,
		BlockHash:     tx.BlockHash,
		BlockOrder:    tx.BlockOrder,
		TxHash:        tx.TxHash,
		Size:          tx.Size,
		Version:       tx.Version,
		Locktime:      tx.Locktime,
		Timestamp:     tx.Timestamp,
		Expire:        tx.Expire,
		Confirmations: tx.Confirmations,
		Txsvaild:      tx.Txsvaild,
		IsCoinbase:    tx.IsCoinbase,
		Vins:          tx.Vins,
		Vouts:         tx.Vouts,
		TotalVin:      qittypes.Amount(tx.TotalVin).ToCoin(),
		TotalVout:     qittypes.Amount(tx.TotalVout).ToCoin(),
		Fees:          qittypes.Amount(tx.Fees).ToCoin(),
		Duplicate:     false,
		Stat:          0,
	}
}

func DBTransactionsToTransactions(dbTxs []*types.Transaction) []*Transaction {
	txs := []*Transaction{}
	for _, tx := range dbTxs {
		txs = append(txs, DBTransactionToTransaction(tx))
	}
	return txs
}

func DBVinoutToVinout(vinout *types.Vinout) *Vinout {
	return &Vinout{
		Id:                     vinout.Id,
		TxId:                   vinout.TxId,
		Type:                   vinout.Type,
		Number:                 vinout.Number,
		Order:                  vinout.Order,
		Timestamp:              vinout.Timestamp,
		Address:                vinout.Address,
		Amount:                 qittypes.Amount(vinout.Amount).ToCoin(),
		ScriptPubKey:           vinout.ScriptPubKey,
		SpentTx:                vinout.SpentTx,
		SpentNumber:            vinout.SpentNumber,
		UnconfirmedSpentTx:     vinout.UnconfirmedSpentTx,
		UnconfirmedSpentNumber: vinout.UnconfirmedSpentNumber,
		SpentedTx:              vinout.SpentedTx,
		Vout:                   vinout.Vout,
		Sequence:               vinout.Sequence,
		ScriptSig:              vinout.ScriptSig,
		Stat:                   vinout.Stat,
	}
}

func DBBlockToBlock(block *types.Block) *Block {
	_, miner := Miners.Get(block.Address)
	return &Block{
		Id:            block.Id,
		Hash:          block.Hash,
		Txvalid:       block.Txvalid,
		Confirmations: block.Confirmations,
		Version:       block.Version,
		Weight:        block.Weight,
		Height:        block.Height,
		TxRoot:        block.TxRoot,
		Order:         block.Order,
		Transactions:  block.Transactions,
		StateRoot:     block.StateRoot,
		Bits:          block.Bits,
		Timestamp:     block.Timestamp,
		ParentRoot:    block.ParentRoot,
		Parents:       block.Parents,
		Children:      block.Children,
		Difficulty:    block.Difficulty,
		PowName:       block.PowName,
		PowType:       block.PowType,
		Nonce:         block.Nonce,
		EdgeBits:      block.EdgeBits,
		CircleNonces:  block.CircleNonces,
		Address:       block.Address,
		Miner:         miner,
		Amount:        qittypes.Amount(block.Amount).ToCoin(),
		Stat:          block.Stat,
	}
}

func DBBlocksToBlocks(dbBlocks []*types.Block) []*Block {
	blocks := []*Block{}
	for _, block := range dbBlocks {
		blocks = append(blocks, DBBlockToBlock(block))
	}
	return blocks
}

type AddressResp struct {
	Id      uint64  `json:"id"`
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}

func ToAddressRespList(addrList []*types2.Address, start uint64) []*AddressResp {
	addrRespList := []*AddressResp{}
	for i, addr := range addrList {
		addrRespList = append(addrRespList, ToAddressResp(addr, start+uint64(i)+1))
	}
	return addrRespList
}

func ToAddressResp(addr *types2.Address, id uint64) *AddressResp {
	return &AddressResp{
		Id:      id,
		Address: addr.Address,
		Balance: qittypes.Amount(addr.Balance).ToCoin(),
	}
}

type AddressStatusResp struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
	Usable  float64 `json:"usable"`
	Locked  float64 `json:"locaked"`
}
