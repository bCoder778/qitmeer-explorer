package qitmeer

import (
	"fmt"
	"github.com/Qitmeer/qitmeer/core/types/pow"
	"github.com/Qitmeer/qitmeer/params"
	"github.com/Qitmeer/qitmeer/qx"
	"math/big"
	"strconv"
)

const (
	H  = "H"
	KH = "K"
	MH = "M"
	GH = "G"
	TH = "T"
	PH = "P"
	EH = "E"
)

const (
	HIndex  = 0
	KHIndex = 1
	MHIndex = 2
	GHIndex = 3
	THIndex = 4
	PHIndex = 5
	EHIndex = 6
)

var uints = map[int]string{
	0: H,
	1: KH,
	2: MH,
	3: GH,
	4: TH,
	5: PH,
	6: EH,
}

func getHashUint(difficulty uint32, uintIndex int, blocktime int) string {
	uint, ok := uints[uintIndex]
	if !ok {
		return H
	}
	val, _ := compactToHashrate(difficulty, uint, blocktime)
	fVal, _ := strconv.ParseFloat(val, 64)
	if fVal < float64(100) {
		return getHashUint(difficulty, uintIndex-1, blocktime)
	}
	return uint
}

func compactToHashrate(input uint32, unit string, blocktime int) (string, string) {
	diffBig := pow.CompactToBig(input)
	maxBig, _ := new(big.Int).SetString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	maxBig.Div(maxBig, diffBig)
	maxBig.Div(maxBig, big.NewInt(int64(blocktime)))
	return qx.GetHashrate(maxBig, unit)
}

func compactToGPS(compact uint32, blockTime, scale int) float64 {
	u64Big := pow.CompactToBig(compact)
	if u64Big.Uint64() <= 0 {
		return 0
	}
	if scale <= 0 {
		return 0
	}
	if blockTime <= 0 {
		return 0
	}
	needGPS := float64(u64Big.Uint64()) / float64(scale) * 50.00 / float64(blockTime)
	return needGPS
}

func getNetWork(network string) *params.Params {
	switch network {
	case "testnet":
		return &params.TestNetParams
	case "privnet":
		return &params.PrivNetParams
	case "mainnet":
		return &params.MainNetParams
	case "mixnet":
		return &params.MixNetParams
	default:
		return &params.TestNetParams
	}
}

func getCuckooScale(powType string, p *params.Params, edgeBits, mheight int64) int {
	switch powType {
	case "cuckaroo":
		instance := &pow.Cuckaroo{}
		instance.SetMainHeight(pow.MainHeight(mheight))
		instance.SetEdgeBits(uint8(edgeBits))
		instance.SetParams(p.PowConfig)
		return int(instance.GraphWeight())
	case "cuckaroom":
		instance := &pow.Cuckaroom{}
		instance.SetMainHeight(pow.MainHeight(mheight))
		instance.SetEdgeBits(uint8(edgeBits))
		instance.SetParams(p.PowConfig)
		return int(instance.GraphWeight())
	case "cuckatoo":
		instance := &pow.Cuckaroo{}
		instance.SetMainHeight(pow.MainHeight(mheight))
		instance.SetEdgeBits(uint8(edgeBits))
		instance.SetParams(p.PowConfig)
		return int(instance.GraphWeight())
	}
	return 0
}

func concurrencyRate(blockTime, mainBlockTime float64) string {
	if mainBlockTime == 0 {
		return "00.00%"
	}
	return fmt.Sprintf("%.2f", mainBlockTime/blockTime*100)
}

func difficulty(diff uint32) uint64 {
	return compactToBig(diff).Uint64() / 30
}

func hashRate(diff uint64) string {
	fDiff24 := float64(diff) / 48 * 50

	return fmt.Sprintf("%.2f kgps", fDiff24/1000)
}

func compactToBig(compact uint32) *big.Int {
	// Extract the mantissa, sign bit, and exponent.
	mantissa := compact & 0x007fffff
	isNegative := compact&0x00800000 != 0
	exponent := uint(compact >> 24)

	// Since the base for the exponent is 256, the exponent can be treated
	// as the number of bytes to represent the full 256-bit number.  So,
	// treat the exponent as the number of bytes and shift the mantissa
	// right or left accordingly.  This is equivalent to:
	// N = mantissa * 256^(exponent-3)
	var bn *big.Int
	if exponent <= 3 {
		mantissa >>= 8 * (3 - exponent)
		bn = big.NewInt(int64(mantissa))
	} else {
		bn = big.NewInt(int64(mantissa))
		bn.Lsh(bn, 8*(exponent-3))
	}

	// Make it negative if the sign bit is set.
	if isNegative {
		bn = bn.Neg(bn)
	}

	return bn
}
