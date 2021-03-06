package controller

import (
	"github.com/bCoder778/qitmeer-explorer/controller/types"
)

type IQitmeer interface {
	AlgorithmList() []*types.AlgorithmResp
	AlgorithmLine(algorithm string, sec int) *types.AlgorithmLineResp
	StartFindPeer() error
	StopFindPeer() error
	PeerList() []*types.PeerResp
	CoinIdList() []string
}
