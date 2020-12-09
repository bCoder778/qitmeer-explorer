package qitmeer

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	"github.com/bCoder778/qitmeer-explorer/db"
	"strconv"
)

type QitmeerV0_9 struct {
	network string
	storage db.IDB
	params  *Params
}

func NewQitmeerV0_9(network string, storage db.IDB) *QitmeerV0_9 {
	return &QitmeerV0_9{network: network, storage: storage, params: Params0_9}
}

func (q *QitmeerV0_9) AlgorithmList() []*types.AlgorithmResp {
	aRespList := []*types.AlgorithmResp{}
	for _, a := range q.params.AlgorithmList {
		block, _ := q.storage.GetLastAlgorithmBlock(a.Name, a.EdgeBits)
		if block != nil {
			ar := q.algorithmResp(block.Difficulty, block.Height, a.ShowName)
			aRespList = append(aRespList, ar)
		}
	}
	return aRespList
}

func (q *QitmeerV0_9) algorithmResp(difficulty, mainHeight uint64, showName string) *types.AlgorithmResp {
	switch showName {
	case Cuckaroom29_Show:
		return q.cuckaroom29(difficulty, mainHeight)
	case Keccak256_Show:
		return q.keccak256(difficulty, mainHeight)
	default:
		return &types.AlgorithmResp{}
	}
	return &types.AlgorithmResp{}
}

func (q *QitmeerV0_9) cuckaroom29(difficulty, mainHeight uint64) *types.AlgorithmResp {
	if mainHeight > 238522 {
		return &types.AlgorithmResp{
			Name:       Cuckaroom29_Show,
			HashRate:   "",
			Difficulty: "",
		}
	}
	val := compactToGPS(uint32(difficulty), 43, getCuckooScale("cuckaroom", getNetWork(q.network), 29, 1))
	gps := fmt.Sprintf("%.2f GPS", val)
	return &types.AlgorithmResp{
		Name:       Cuckaroom29_Show,
		HashRate:   gps,
		Difficulty: fmt.Sprintf("%d", difficulty),
	}
}

func (q *QitmeerV0_9) keccak256(difficulty, mainHeight uint64) *types.AlgorithmResp {
	if difficulty == 0 {
		return &types.AlgorithmResp{
			Name:       Keccak256_Show,
			HashRate:   "",
			Difficulty: "",
		}
	}
	blockTime := 100
	// 大于238522主链高度后，Keccak256 blocktime改变，无cuckaroom29
	if mainHeight > 238522 {
		blockTime = 30
	}

	uint := getHashUint(uint32(difficulty), 6, 1)

	val := compactToHashrate(uint32(difficulty), uint, 1)
	fVal, _ := strconv.ParseFloat(val, 64)

	keccak256Diff := fmt.Sprintf("%.2f %s", fVal, uint)

	val = compactToHashrate(uint32(difficulty), uint, blockTime)
	fVal, _ = strconv.ParseFloat(val, 64)
	hashRate := fmt.Sprintf("%.2f %s", fVal, uint)
	return &types.AlgorithmResp{
		Name:       Keccak256_Show,
		HashRate:   hashRate,
		Difficulty: keccak256Diff,
	}
}
