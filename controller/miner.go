package controller

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
)

func (c *Controller) BlocksDistribution() []*types.DistributionResp {
	minerstatus := c.storage.BlocksDistribution()
	distributions := map[string]*types.DistributionResp{}
	rs := []*types.DistributionResp{}
	var all uint64
	for _, miner := range minerstatus {
		all += miner.Count
	}
	for _, miner := range minerstatus {
		distributuon := &types.DistributionResp{}
		ok, pool := types.Miners.Get(miner.Address)
		if ok {
			distributuon.Miner = pool.Name
		} else {
			distributuon.Miner = miner.Address
		}
		block := c.storage.GetLastMinerBlock(miner.Address)
		distributuon.Blocks = miner.Count
		distributuon.LastOrder = block.Order
		distributuon.LastTimestamp = block.Timestamp
		distributuon.Proportion = blocksProportion(distributuon.Blocks, all)
		dt, ok := distributions[distributuon.Miner]
		if ok {
			if dt.LastOrder < distributuon.LastOrder {
				dt.LastOrder = distributuon.LastOrder
				dt.LastTimestamp = distributuon.LastTimestamp
			}
			dt.Blocks += distributuon.Blocks
			dt.Proportion = blocksProportion(distributuon.Blocks, all)
			distributions[distributuon.Miner] = dt
		} else {
			distributions[distributuon.Miner] = distributuon
		}
	}
	for _, dt := range distributions {
		rs = append(rs, dt)
	}
	return rs
}

func blocksProportion(blocks, all uint64) string {
	return fmt.Sprintf("%.2f %%", (float64(blocks) / float64(all) * 100))
}
