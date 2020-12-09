package qitmeer

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	db "github.com/bCoder778/qitmeer-explorer/db"
	"strconv"
)

type QitmeerV0_10 struct {
	network string
	storage db.IDB
	params  *Params
}

func NewQitmeerV0_10(network string, storage db.IDB) *QitmeerV0_10 {
	return &QitmeerV0_10{network: network, storage: storage, params: Params0_10}
}

func (q *QitmeerV0_10) AlgorithmList() []*types.AlgorithmResp {
	aRespList := []*types.AlgorithmResp{}
	for _, a := range q.params.AlgorithmList {
		block, _ := q.storage.GetLastAlgorithmBlock(a.Name, a.EdgeBits)
		if block != nil {
			ar := q.algorithmResp(block.Difficulty, a.ShowName)
			aRespList = append(aRespList, ar)
		}
	}
	return aRespList
}

func (q *QitmeerV0_10) algorithmResp(difficulty uint64, showName string) *types.AlgorithmResp {
	switch showName {
	case Cuckaroo_Show:
		return q.cuckaroo(difficulty)
	case Keccak256_Show:
		return q.hashRate(difficulty, showName)
	case Cryptonight_Show:
		return q.hashRate(difficulty, showName)
	case Blake2b_Show:
		return q.hashRate(difficulty, showName)
	default:
		return &types.AlgorithmResp{}
	}
	return &types.AlgorithmResp{}
}

func (q *QitmeerV0_10) cuckaroo(difficulty uint64) *types.AlgorithmResp {
	val := compactToGPS(uint32(difficulty), 43, getCuckooScale("cuckaroo", getNetWork(q.network), 24, 1))
	gps := fmt.Sprintf("%.2f GPS", val)
	return &types.AlgorithmResp{
		Name:       Cuckaroo_Show,
		HashRate:   gps,
		Difficulty: fmt.Sprintf("%d", difficulty),
	}
}

func (q *QitmeerV0_10) hashRate(difficulty uint64, showName string) *types.AlgorithmResp {
	if difficulty == 0 {
		return &types.AlgorithmResp{
			Name:       showName,
			HashRate:   "",
			Difficulty: "",
		}
	}
	blockTime := 30

	uint := getHashUint(uint32(difficulty), 6, 1)

	val := compactToHashrate(uint32(difficulty), uint, 1)
	fVal, _ := strconv.ParseFloat(val, 64)

	hashRateDiff := fmt.Sprintf("%.2f %s", fVal, uint)

	val = compactToHashrate(uint32(difficulty), uint, blockTime)
	fVal, _ = strconv.ParseFloat(val, 64)
	hashRate := fmt.Sprintf("%.2f %s", fVal, uint)
	return &types.AlgorithmResp{
		Name:       showName,
		HashRate:   hashRate,
		Difficulty: hashRateDiff,
	}
}
