package controller

import (
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	dbtypes "github.com/bCoder778/qitmeer-explorer/db/types"
)

type IQitmeer interface {
	AlgorithmList() []*types.AlgorithmResp
	AlgorithmLine(algorithm string, sec int) *types.AlgorithmLineResp
	StartFindPeer() error
	StopFindPeer() error
	PeerList() []*dbtypes.Peer
}
