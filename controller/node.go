package controller

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	synctypes "github.com/bCoder778/qitmeer-sync/storage/types"
)

func (c *Controller) NodeList() interface{} {
	return c.qitmeer.PeerList()
}

func (c *Controller) Tips() *types.TipsResp {
	last, _ := c.storage.GetLastBlock()
	avg, _ := c.blockTimeAvg(last)
	tips := &types.TipsResp{
		BlockAvg:          fmt.Sprintf("%.2f", avg.OrderAvgTime),
		BlockInterval:     fmt.Sprintf("[%d - %d]", avg.OldOrder, avg.Order),
		MainBlockAvg:      fmt.Sprintf("%.2f", avg.MainAvgTime),
		MainBlockInterval: fmt.Sprintf("[%d - %d]", avg.OldHeight, avg.Height),
		ConcurrencyRate:   concurrencyRate(avg.OrderAvgTime, avg.MainAvgTime),
		BlockOrder:        last.Order,
		BlockHeight:       last.Height,
	}
	return tips
}

type avgTime struct {
	Order        uint64
	OldOrder     uint64
	OrderAvgTime float64
	Height       uint64
	OldHeight    uint64
	MainAvgTime  float64
}

func (c *Controller) blockTimeAvg(curBlock *synctypes.Block) (*avgTime, bool) {
	var perBlockOrder uint64
	if curBlock.Order >= 15 {
		perBlockOrder = curBlock.Order - 15
	}
	preBlock, _ := c.storage.GetBlockByOrder(perBlockOrder)

	return &avgTime{
		Order:        curBlock.Order,
		OldOrder:     perBlockOrder,
		OrderAvgTime: float64(curBlock.Timestamp-preBlock.Timestamp) / float64(curBlock.Order-perBlockOrder-1),
		Height:       curBlock.Height,
		OldHeight:    preBlock.Height,
		MainAvgTime:  float64(curBlock.Timestamp-preBlock.Timestamp) / float64(curBlock.Height-preBlock.Height-1),
	}, true
}

func concurrencyRate(blockTime, mainBlockTime float64) string {
	if mainBlockTime == 0 {
		return "00.00%"
	}
	return fmt.Sprintf("%.2f", mainBlockTime/blockTime*100)
}
