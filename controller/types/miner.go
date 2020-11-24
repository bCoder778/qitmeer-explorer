package types

type MinerPool struct {
	Name string
	Url  string
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
	"TmPrXkjpjSUBiFG9RZKPjfdsAPbiaar94Ta": &MinerPool{
		Name: "666pool.cn",
		Url:  "https://www.666pool.cn",
	},
	"TmVfDq18VqSg735ko9aAo36tFwYww4PBGMC": {
		Name: "meerpool.com",
		Url:  "https://www.meerpool.com",
	},
	"TmekMwXHgk6NHD2i9ZtHeWfnC8ypfPvxvgf": {
		Name: "meerpool.com",
		Url:  "https://www.meerpool.com",
	},
	"TmUHh6bAdLbto9AYhodEwGZi9WY77CoBFXr": {
		Name: "Hashpool",
		Url:  "https://hashpool.com/",
	},
	"TmRfnZPT3r93WG1LhHdEgeJi36gmPN4MytD": {
		Name: "uupool.cn",
		Url:  "https://uupool.cn",
	},
	"TmX7F5yEq65Yb5x3uHmNNEVZ7DwwqZxTKtg": {
		Name: "hpt.com",
		Url:  "https://hpt.com",
	},
}
