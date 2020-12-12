package qitmeer

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	"github.com/bCoder778/qitmeer-explorer/db"
	dbtypes "github.com/bCoder778/qitmeer-sync/storage/types"
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
		block, _ := q.storage.GetLastAlgorithmBlock(a.DBName, a.EdgeBits)
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

	val, _ := compactToHashrate(uint32(difficulty), uint, 1)
	fVal, _ := strconv.ParseFloat(val, 64)

	keccak256Diff := fmt.Sprintf("%.2f %s", fVal, uint)

	val, finalUint := compactToHashrate(uint32(difficulty), uint, blockTime)
	fVal, _ = strconv.ParseFloat(val, 64)
	hashRate := fmt.Sprintf("%.2f %s", fVal, finalUint)
	return &types.AlgorithmResp{
		Name:       Keccak256_Show,
		HashRate:   hashRate,
		Difficulty: keccak256Diff,
	}
}

func (q *QitmeerV0_9) AlgorithmLine(algorithm string, sec int) *types.AlgorithmLineResp {
	alr := &types.AlgorithmLineResp{Name: algorithm, Sec: sec, Avgs: []*types.AlgorithmAvg{}}
	a, ok := q.params.AlgorithmList[algorithm]
	if !ok {
		return alr
	}
	last, err := q.storage.GetLastBlock()
	if err != nil {
		return &types.AlgorithmLineResp{}
	}
	lastTime := last.Timestamp

	avgValue := []*types.AlgorithmAvg{}

	max := lastTime
	min := lastTime - int64(16*sec)
	blocks := q.storage.QueryAlgorithmDiffInTime(a.DBName, a.EdgeBits, max, min)
	if len(blocks) == 0 {
		return alr
	}
	for i := 16; i > 0; i-- {
		maxTime := lastTime - int64((i-1)*sec)
		minTime := lastTime - int64(i*sec)
		blockList := getBlockList(blocks, maxTime, minTime)
		value, uint := q.avgAlgorithmRate(blockList, algorithm, blocks[len(blocks)-1].Difficulty)
		avgValue = append(avgValue, &types.AlgorithmAvg{
			Value: value,
			Uint:  uint,
			Time:  maxTime,
		})
	}
	alr.Avgs = avgValue
	return alr
}

func (q *QitmeerV0_9) avgAlgorithmRate(blocks []*dbtypes.Block, algorithm string, lastDiff uint64) (string, string) {
	switch algorithm {
	case Cuckaroom29_Show:
		return q.avgGPS(blocks, 29)
	case Keccak256_Show:
		return q.avgHashRate(blocks, lastDiff)
	}
	return "", ""
}

func (q *QitmeerV0_9) avgGPS(blocks []*dbtypes.Block, edgeBits int64) (string, string) {
	var sum float64
	for _, block := range blocks {
		sum += compactToGPS(uint32(block.Difficulty), 43, getCuckooScale("cuckaroom", getNetWork(q.network), edgeBits, 1))
	}
	if len(blocks) == 0 {
		return "0", "GPS"
	}
	return fmt.Sprintf("%.2f", sum/float64(len(blocks))), "GPS"
}

func (q *QitmeerV0_9) avgHashRate(blocks []*dbtypes.Block, lastDiff uint64) (string, string) {
	var sum float64
	length := len(blocks)
	if length == 0 {
		return "0", "H"
	}
	finalUint := ""
	val := ""
	uint := getHashUint(uint32(lastDiff), 6, 1)
	for _, block := range blocks {
		blockTime := 30
		if block.Height <= 238522 {
			blockTime = 100
		}
		val, finalUint = compactToHashrate(uint32(block.Difficulty), uint, blockTime)
		fVal, _ := strconv.ParseFloat(val, 64)
		sum += fVal
	}
	return fmt.Sprintf("%.2f", sum/float64(length)), finalUint
}

func getBlockList(blocks []*dbtypes.Block, max, min int64) []*dbtypes.Block {
	blockList := []*dbtypes.Block{}
	for _, b := range blocks {
		if b.Timestamp <= max && b.Timestamp > min {
			blockList = append(blockList, b)
		}
	}
	return blockList
}

func (q *QitmeerV0_9) NodeList() {

}
