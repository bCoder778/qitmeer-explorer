package types

type Address struct {
	Address string `xorm:"address" json:"address"`
	Balance uint64 `xorm:"balance" json:"amount"`
}

type MinerStatus struct {
	Address string
	Count   uint64
}
