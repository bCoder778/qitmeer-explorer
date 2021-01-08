package types

type Address struct {
	Address string `xorm:"address" json:"address"`
	Balance uint64 `xorm:"balance" json:"amount"`
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
