package types

type MinerPool struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type MinerPools map[string]*MinerPool

func (m MinerPools) Get(address string) (bool, *MinerPool) {
	miner, ok := m[address]
	if !ok {
		return false, &MinerPool{}
	}
	return ok, miner
}

var Miners = MinerPools{
	"MmPS2J6VqkJEotaxQSakvEkDsHHV6oJniaV": {
		Name: "MeerPool",
		Url:  "https://www.meerpool.com",
	},
	"MmQvwPfLRBZ4ujEwmdj8Fp7Ge6CFAEoiof5 ": {
		Name: "HashPool",
		Url:  "https://hashpool.com",
	},
	"MmcepjS9G4oNhQrhKs6J6XaVALsKAVEWz5F": {
		Name: "HuobiPool",
		Url:  "https://hpt.com",
	},
	"MmUt7mUmQSCSAU7FJPmKAEHdiGsKkrKhdoD": {
		Name: "F2Pool",
		Url:  "https://www.f2pool.com",
	},
}

type DistributionResp struct {
	Miner            string  `json:"miner"`
	Address          string  `json:"address"`
	Blocks           uint64  `json:"blocks"`
	Proportion       string  `json:"proportion"`
	ProportionNumber float64 `json:"proportionNum"`
	LastOrder        uint64  `json:"lastorder"`
	LastTimestamp    int64   `json:"lasttimestamp"`
}

type DistributionsResp struct {
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
	List  []*DistributionResp `json:"list"`
	Count int64               `json:"count"`
}
