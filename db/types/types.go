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

type Package struct {
	MaxTime int64 `xorm:"bigInt maxTime" json:"maxTime"`
	MinTime int64 `xorm:"bigInt minTime" json:"minTime"`
	AvgTime float64 `xorm:"bigInt avgTime" json:"avgTime"`
}