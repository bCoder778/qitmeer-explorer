package types

type Address struct {
	Address string `xorm:"address" json:"address"`
	Balance uint64 `xorm:"balance" json:"amount"`
	CoinId  string `xorm:"coin_id"`
}

type MinerStatus struct {
	Address string
	Count   uint64
}

type Peer struct {
	Id       uint64 `xorm:"bigint autoincr pk" json:"id"`
	Address  string `xorm:"varchar(64)" json:"address"`
	FindTime int64  `xorm:"bigint" json:"findtime"`
	Other    string `xorm:"varchar(64)" json:"other"`
}

type TimeInfo struct {
	WaitTime int64 `json:"waitTime"`
	BlockHash string `json:"blockHash"`
	TxId string `json:"txId"`
}

type Package struct {
	MaxInfo *TimeInfo `json:"maxTime"`
	MinInfo *TimeInfo `json:"minTime"`
	AvgTime float64 `json:"avgTime"`
	TxCount int64 `json:"txCount"`
	SumTime int64 `json:"sumTime"`
}