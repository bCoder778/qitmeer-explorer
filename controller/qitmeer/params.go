package qitmeer

const (
	Cuckaroom29_Show = "Cuckaroom 29"
	Keccak256_Show   = "Keccak256"
	Cuckaroo_Show    = "Cuckaroo"
	Cryptonight_Show = "Cryptonight"
	Blake2b_Show     = "Blake2b"
)

const (
	Cuckaroom29_DB = "cuckaroom"
	Keccak256_DB   = "qitmeer_keccak256"
	Cuckaroo_DB    = "cuckaroo"
	Cryptonight_DB = "cryptonight"
	Blake2b_DB     = "blake2bd"
)

type Params struct {
	AlgorithmList []*Algorithm
}

type Algorithm struct {
	ShowName string
	Name     string
	EdgeBits int
}

var Params0_9 = &Params{
	AlgorithmList: []*Algorithm{
		{
			ShowName: Cuckaroom29_Show,
			Name:     Cuckaroom29_DB,
			EdgeBits: 29,
		},
		{
			ShowName: Keccak256_Show,
			Name:     Keccak256_DB,
			EdgeBits: 0,
		},
	},
}

var Params0_10 = &Params{
	AlgorithmList: []*Algorithm{
		{
			ShowName: Cuckaroo_Show,
			Name:     Cuckaroo_DB,
			EdgeBits: 24,
		},
		{
			ShowName: Keccak256_Show,
			Name:     Keccak256_DB,
			EdgeBits: 0,
		},
		{
			ShowName: Cryptonight_Show,
			Name:     Cryptonight_DB,
			EdgeBits: 0,
		},
		{
			ShowName: Blake2b_Show,
			Name:     Blake2b_DB,
			EdgeBits: 0,
		},
	},
}
