package types

import (
	qitTypes "github.com/Qitmeer/qitmeer/core/types"
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

type BlockDetailResp struct {
	Header       *BlockResp               `json:"header"`
	Transactions []*TransactionDetailResp `json:"transactions"`
}

type TransactionDetailResp struct {
	Header *TransactionResp `json:"header"`
	Vin    []*VinResp       `json:"vin"`
	Vout   []*VoutResp      `json:"vout"`
}

type TransactionResp struct {
	Id            uint64            `json:"id"`
	TxId          string            `json:"txid"`
	BlockHash     string            `json:"blockhash"`
	BlockOrder    uint64            `json:"blockorder"`
	TxHash        string            `json:"txhash"`
	Size          int               `json:"size"`
	Version       uint32            `json:"version"`
	Locktime      uint64            `json:"locktime"`
	Timestamp     int64             `json:"timestamp"`
	Expire        uint64            `json:"expire"`
	Confirmations uint64            `json:"confirmations"`
	Txsvaild      bool              `json:"txsvaild"`
	IsCoinbase    bool              `json:"iscoinbase"`
	VinAmount    float64 		    `json:"vinamount"`
	VoutAmount    float64 		    `json:"voutamount"`
	VinAddress    string 		    `json:"vinaddress"`
	VoutAddress   string 		    `json:"voutaddress"`
	Vins          int               `json:"vin"`
	Vouts         int               `json:"vout"`
	Fees          []*Fees           `json:"fees"`
	Changes       []*TransferChange `json:"changes"`
	Duplicate     bool              `json:"duplicate"`
	Miner         *MinerPool        `json:"miner"`
	Stat          stat.TxStat       `json:"stat"`
}

type TransferChange struct {
	CoinId     string  `json:"coinid"`
	Change     float64 `json:"change"`
	TotalVin   float64 `json:"totalvin"`
	TotalVout  float64 `json:"totalvout"`
	UTotalVin  uint64  `json:"utotalvin"`
	UTotalVout uint64  `json:"utotalvout"`
}

type Fees struct {
	CoinId     string  `json:"coinid"`
	TotalVin   float64 `json:"totalvin"`
	TotalVout  float64 `json:"totalvout"`
	UTotalVin  uint64  `json:"utotalvin"`
	UTotalVout uint64  `json:"utotalvout"`
	Amount     float64 `json:"amount"`
}

type VinoutResp struct {
	Id                     uint64              `json:"id"`
	TxId                   string              `json:"txid"`
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

type VinResp struct {
	Id        uint64           `json:"id"`
	TxId      string           `json:"txid"`
	Number    int              `json:"number"`
	Order     uint64           `json:"order"`
	Timestamp int64            `json:"timestamp"`
	Address   string           `json:"address"`
	CoinId    string           `json:"coinid"`
	Amount    float64          `json:"amount"`
	SpentedTx string           `json:"spentedtx"`
	Vout      int              `json:"vout"`
	Sequence  uint64           `json:"sequence"`
	ScriptSig *types.ScriptSig `json:"scriptsig"`
	Stat      stat.TxStat      `json:"stat"`
}

type VoutResp struct {
	Id           uint64              `json:"id"`
	TxId         string              `json:"txid"`
	Number       int                 `json:"number"`
	Order        uint64              `json:"order"`
	Timestamp    int64               `json:"timestamp"`
	Address      string              `json:"address"`
	CoinId       string              `json:"coinid"`
	Amount       float64             `json:"amount"`
	ScriptPubKey *types.ScriptPubKey `json:"scriptpubkey"`
	SpentTx      string              `json:"spenttx"`
	Locked       bool                `json:"locked"`
	LockHeight   uint64              `json:"lockheight"`
	Stat         stat.TxStat         `json:"stat"`
}

type BlockResp struct {
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
	Nonce         string         `json:"nonce"`
	EdgeBits      int            `json:"edgebits"`
	CircleNonces  string         `json:"circlenonces"`
	Address       string         `json:"address"`
	Amount        float64        `json:"amount"`
	Miner         *MinerPool     `json:"miner"`
	Color         stat.Color     `json:"color"`
	Stat          stat.BlockStat `json:"stat"`
}

func ToTransactionResp(tx *types.Transaction) *TransactionResp {
	vinAmount := qitTypes.Amount{
		Id:    qitTypes.MEERID,
		Value: int64(tx.VinAmount),
	}
	voutAmount := qitTypes.Amount{
		Id:    qitTypes.MEERID,
		Value: int64(tx.VoutAmount),
	}
	miner := &MinerPool{
		Name: "",
		Url:  "",
	}
	if tx.IsCoinbase{
		_, miner = Miners.Get(tx.VoutAddress)
	}
	t := &TransactionResp{
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
		VinAmount:     vinAmount.ToCoin(),
		VoutAmount:    voutAmount.ToCoin(),
		VinAddress:    tx.VinAddress,
		VoutAddress:   tx.VoutAddress,
		Vins:          tx.Vins,
		Vouts:         tx.Vouts,
		Fees:          nil,
		Changes:       nil,
		Duplicate:     tx.Duplicate,
		Miner: miner,
		Stat:          tx.Stat,
	}
	return t
}

func ToTransactionRespList(dbTxs []*types.Transaction) []*TransactionResp {
	var txs []*TransactionResp
	for _, tx := range dbTxs {
		txs = append(txs, ToTransactionResp(tx))
	}
	return txs
}

func ToVinResp(vinout *types.Vin) *VinResp {

	amount := qitTypes.Amount{
		Id:    qitTypes.NewCoinID(vinout.CoinId),
		Value: int64(vinout.Amount),
	}
	return &VinResp{
		Id:        vinout.Id,
		TxId:      vinout.TxId,
		Number:    vinout.Number,
		Order:     vinout.Order,
		Timestamp: vinout.Timestamp,
		Address:   vinout.Address,
		CoinId:    vinout.CoinId,
		Amount:    amount.ToCoin(),
		SpentedTx: vinout.SpentedTx,
		Vout:      vinout.Vout,
		Sequence:  vinout.Sequence,
		ScriptSig: vinout.ScriptSig,
		Stat:      vinout.Stat,
	}
}

func ToVoutListResp(vs []*types.Vout, height uint64) []*VoutResp {
	var outs []*VoutResp
	for _, item := range vs {
		outs = append(outs, ToVoutResp(item, height))
	}
	return outs
}

func ToVoutResp(vinout *types.Vout, height uint64) *VoutResp {
	amount := qitTypes.Amount{
		Id:    qitTypes.NewCoinID(vinout.CoinId),
		Value: int64(vinout.Amount),
	}

	return &VoutResp{
		Id:           vinout.Id,
		TxId:         vinout.TxId,
		Number:       vinout.Number,
		Order:        vinout.Order,
		Timestamp:    vinout.Timestamp,
		Address:      vinout.Address,
		CoinId:       vinout.CoinId,
		Amount:       amount.ToCoin(),
		ScriptPubKey: vinout.ScriptPubKey,
		SpentTx:      vinout.SpentTx,
		LockHeight:   vinout.Lock,
		Locked:       height >= vinout.Lock,
		Stat:  		  vinout.Stat,
	}
}

func ToBlockResp(block *types.Block) *BlockResp {
	_, miner := Miners.Get(block.Address)
	amount := qitTypes.Amount{
		Id:    qitTypes.MEERID,
		Value: int64(block.Amount),
	}
	return &BlockResp{
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
		Amount:        amount.ToCoin(),
		Color:         block.Color,
		Stat:          block.Stat,
	}
}

func ToBlockRespList(dbBlocks []*types.Block) []*BlockResp {
	var blocks []*BlockResp
	for _, block := range dbBlocks {
		blocks = append(blocks, ToBlockResp(block))
	}
	return blocks
}

type AddressResp struct {
	Id      uint64  `json:"id"`
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
	Tag     string  `json:"tag"`
}

func ToAddressRespList(addrList []*types2.Address, start uint64) []*AddressResp {
	var addrRespList []*AddressResp
	for i, addr := range addrList {
		addrRespList = append(addrRespList, ToAddressResp(addr, start+uint64(i)+1))
	}
	return addrRespList
}

func ToAddressResp(addr *types2.Address, id uint64) *AddressResp {

	amount := qitTypes.Amount{
		Id:    qitTypes.NewCoinID(addr.CoinId),
		Value: int64(addr.Balance),
	}

	return &AddressResp{
		Id:      id,
		Address: addr.Address,
		Balance: amount.ToCoin(),
	}
}

type AddressStatusResp struct {
	Address    string  `json:"address"`
	Balance    float64 `json:"balance"`
	Usable     float64 `json:"usable"`
	Locked     float64 `json:"locaked"`
	Uncofirmed float64 `json:"uncofirmed"`
}

type AlgorithmResp struct {
	Name       string `json:"name"`
	HashRate   string `json:"hashrate"`
	Difficulty string `json:"difficulty"`
}

type AlgorithmAvg struct {
	Value string
	Time  int64
	Uint  string
}

type AlgorithmLineResp struct {
	Name string
	Sec  int
	Avgs []*AlgorithmAvg
}

type PeerResp struct {
	Id       uint64    `json:"id"`
	Addr     string    `json:"addr"`
	Other    string    `json:"other"`
	Location *Location `json:"location"`
}

type Location struct {
	City string  `json:"city"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type TipsResp struct {
	BlockAvg          string `json:"blockavg"`
	BlockInterval     string `json:"blockinterval"`
	MainBlockAvg      string `json:"mainblockavg"`
	MainBlockInterval string `json:"mainblockinterval"`
	ConcurrencyRate   string `json:"concurrencyrate"`
	BlockOrder        uint64 `json:"blockorder"`
	BlockHeight       uint64 `json:"blockheight"`
}

type Package struct {
	MaxTime int64 `json:"maxTime"`
	MinTime int64 `json:"minTime"`
	AvgTime int64 `json:"avgtime"`
}