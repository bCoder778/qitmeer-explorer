package controller

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	dbtype "github.com/bCoder778/qitmeer-explorer/db/types"
	synctypes "github.com/bCoder778/qitmeer-sync/storage/types"
	"strconv"
	"time"
)

func (c *Controller) NodeList() interface{} {
	value, err := c.cache.Value("NodeList", "NodeList")
	if err != nil {
		list := c.queryNode()
		c.cache.Add("NodeList", "NodeList", time.Minute*60, list)
		return list
	}
	return value.([]*types.PeerResp)
}

func (c *Controller) queryNode() []*types.PeerResp {
	rs := c.storage.QueryLocation()

	list := make([]*types.PeerResp, 0)
	for _, item := range rs {
		list = append(list, &types.PeerResp{
			Id:    uint64(item.Id),
			Other: item.Other,
			Addr:  item.IpAddress,
			Location: &types.Location{
				City: item.City,
				Lon:  item.Lon,
				Lat:  item.Lat,
			},
		})
	}
	go c.qitmeer.PeerList()
	return list
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
	var blockCount uint64 = 15
	var perBlockOrder uint64
	var preBlock *synctypes.Block
	var err error
	for blockCount < 200 {
		if curBlock.Order >= blockCount {
			perBlockOrder = curBlock.Order - blockCount
		}
		preBlock, err = c.storage.GetBlockByOrder(perBlockOrder)
		if err != nil || preBlock.Order == 0 {
			blockCount++
		} else {
			break
		}
	}

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

func (c *Controller) PackageTime(count string) *dbtype.Package {
	value, err := c.cache.Value("packageTime", count)
	if err != nil {
		info := c.packageTime(count)
		c.cache.Add("packageTime", count, 1*time.Second*5, info)
		return info
	}
	return value.(*dbtype.Package)
}

func (c *Controller) packageTime(count string) *dbtype.Package {
	iCount, _ := strconv.Atoi(count)
	return c.storage.PackageTime(iCount)
}

func (c *Controller) AvgBlocks(sec uint64) float64 {
	return c.storage.AvgBlocks(sec)
}
