package qitmeer

import (
	"fmt"
	qts "github.com/Qitmeer/qitmeer/core/types"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	db "github.com/bCoder778/qitmeer-explorer/db"
	dbtype "github.com/bCoder778/qitmeer-explorer/db/types"
	"github.com/bCoder778/qitmeer-explorer/external"
	"github.com/bCoder778/qitmeer-explorer/rpc"
	dbtypes "github.com/bCoder778/qitmeer-sync/storage/types"
	"regexp"
	"strconv"
)

const (
	node_rpc_host = "https://testnet.meerscan.io/crawler"
	node_rpc_user = "admin"
	node_rpc_pas  = "123"
)

type QitmeerV0_10 struct {
	network string
	storage db.IDB
	params  *Params
	nodeRpc *rpc.Client
}

func NewQitmeerV0_10(network string, storage db.IDB) *QitmeerV0_10 {
	client := rpc.NewClient(node_rpc_host, node_rpc_user, node_rpc_pas)
	return &QitmeerV0_10{network: network, storage: storage, params: Params0_10, nodeRpc: client}
}

func (q *QitmeerV0_10) StartFindPeer() error {
	return nil
}

func (q *QitmeerV0_10) StopFindPeer() error {
	return nil
}

func (q *QitmeerV0_10) PeerList() []*types.PeerResp {
	peers, err := q.nodeRpc.GetNodeList()
	if err != nil {
		return nil
	}
	var rs []*types.PeerResp
	ipMap := map[string]bool{}
	for i, p := range peers {
		r, _ := regexp.Compile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`)
		ip := string(r.Find([]byte(p.Ip)))

		local := q.storage.GetLocation(ip)
		loc := &types.Location{
			City: "",
			Lat:  0,
			Lon:  0,
		}
		if local == nil || local.Id < 1 {
			loc = getLocation(ip)
		} else {
			loc.City = local.City
			loc.Lat = local.Lat
			loc.Lon = local.Lon
		}
		rs = append(rs, &types.PeerResp{
			Id:       uint64(i),
			Addr:     ip,
			Other:    p.Id,
			Location: loc,
		})
		ipMap[ip] = true

		if len(loc.City) > 0 {
			_ = q.storage.UpdateLocation(&dbtype.Location{
				IpAddress: ip,
				City:      loc.City,
				Lon:       loc.Lon,
				Lat:       loc.Lat,
				Other:     p.Id,
			})
		}
	}
	return rs
}

func getLocation(ip string) *types.Location {
	addr, err := external.GetIpInfo(ip)
	if err != nil {
		return &types.Location{}
	}
	return &types.Location{
		City: addr.City,
		Lon:  addr.Lon,
		Lat:  addr.Lat,
	}
}

func (q *QitmeerV0_10) AlgorithmList() []*types.AlgorithmResp {
	aRespList := []*types.AlgorithmResp{}
	for _, a := range q.params.AlgorithmList {
		block, _ := q.storage.GetLastAlgorithmBlock(a.DBName, a.EdgeBits)
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
	case MeerXkeccakV1_Show:
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

	val, _ := compactToHashrate(uint32(difficulty), uint, 1)
	fVal, _ := strconv.ParseFloat(val, 64)
	hashRateDiff := fmt.Sprintf("%.2f %s", fVal, uint)

	val, finalUint := compactToHashrate(uint32(difficulty), uint, blockTime)
	fVal, _ = strconv.ParseFloat(val, 64)
	hashRate := fmt.Sprintf("%.2f %s", fVal, finalUint)
	return &types.AlgorithmResp{
		Name:       showName,
		HashRate:   hashRate,
		Difficulty: hashRateDiff,
	}
}

func (q *QitmeerV0_10) AlgorithmLine(algorithm string, sec int) *types.AlgorithmLineResp {
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

func (q *QitmeerV0_10) avgAlgorithmRate(blocks []*dbtypes.Block, algorithm string, lastDiff uint64) (string, string) {
	switch algorithm {
	case Cuckaroo_Show:
		return q.avgGPS(blocks, 24)
	case Cryptonight_Show:
		return q.avgHashRate(blocks, lastDiff)
	case Keccak256_Show:
		return q.avgHashRate(blocks, lastDiff)
	case Blake2b_Show:
		return q.avgHashRate(blocks, lastDiff)
	case MeerXkeccakV1_Show:
		return q.avgHashRate(blocks, lastDiff)
	}
	return "", ""
}

func (q *QitmeerV0_10) avgGPS(blocks []*dbtypes.Block, edgeBits int64) (string, string) {
	var sum float64
	for _, block := range blocks {
		sum += compactToGPS(uint32(block.Difficulty), 43, getCuckooScale("cuckaroom", getNetWork(q.network), edgeBits, 1))
	}
	if len(blocks) == 0 {
		return "0", "GPS"
	}
	return fmt.Sprintf("%.2f", sum/float64(len(blocks))), "GPS"
}

func (q *QitmeerV0_10) avgHashRate(blocks []*dbtypes.Block, lastDiff uint64) (string, string) {
	var sum float64
	length := len(blocks)
	if length == 0 {
		return "0", "H"
	}
	finalUint := ""
	val := ""
	uint := getHashUint(uint32(lastDiff), 6, 1)
	blockTime := 30

	for _, block := range blocks {
		val, finalUint = compactToHashrate(uint32(block.Difficulty), uint, blockTime)
		fVal, _ := strconv.ParseFloat(val, 64)
		sum += fVal
	}

	return fmt.Sprintf("%.2f", sum/float64(length)), finalUint
}

func (q *QitmeerV0_10) CoinIdList() []string {
	var coins []string
	for _, item := range qts.CoinIDList {
		coins = append(coins, item.Name())
	}
	return coins
}
