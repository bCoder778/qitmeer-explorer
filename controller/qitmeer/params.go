package qitmeer

const (
	Cuckaroom29_Show   = "Cuckaroom29"
	Keccak256_Show     = "Keccak256"
	Cuckaroo_Show      = "Cuckaroo"
	Cryptonight_Show   = "Cryptonight"
	Blake2b_Show       = "Blake2b"
	MeerXkeccakV1_Show = "MeerKeccakV1"
)

const (
	Cuckaroom29_DB   = "cuckaroom"
	Keccak256_DB     = "qitmeer_keccak256"
	Cuckaroo_DB      = "cuckaroo"
	Cryptonight_DB   = "cryptonight"
	Blake2b_DB       = "blake2bd"
	MeerXkeccakV1_DB = "meer_xkeccak_v1"
)

type Params struct {
	AlgorithmList map[string]*Algorithm
}

type Algorithm struct {
	ShowName string
	DBName   string
	EdgeBits int
}

var Params0_9 = &Params{
	AlgorithmList: map[string]*Algorithm{
		Cuckaroom29_Show: {
			ShowName: Cuckaroom29_Show,
			DBName:   Cuckaroom29_DB,
			EdgeBits: 29,
		},
		Keccak256_Show: {
			ShowName: Keccak256_Show,
			DBName:   Keccak256_DB,
			EdgeBits: 0,
		},
	},
}

var Params0_10 = &Params{
	AlgorithmList: map[string]*Algorithm{
		Cuckaroo_Show: {
			ShowName: Cuckaroo_Show,
			DBName:   Cuckaroo_DB,
			EdgeBits: 24,
		},
		Keccak256_Show: {
			ShowName: Keccak256_Show,
			DBName:   Keccak256_DB,
			EdgeBits: 0,
		},
		Cryptonight_Show: {
			ShowName: Cryptonight_Show,
			DBName:   Cryptonight_DB,
			EdgeBits: 0,
		},
		Blake2b_Show: {
			ShowName: Blake2b_Show,
			DBName:   Blake2b_DB,
			EdgeBits: 0,
		},
		MeerXkeccakV1_Show: {
			ShowName: MeerXkeccakV1_Show,
			DBName:   MeerXkeccakV1_DB,
			EdgeBits: 0,
		},
	},
}
