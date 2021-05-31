package controller

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	"time"
)

func (c *Controller) BlocksDistribution() []*types.DistributionResp {
	value, err := c.cache.Value("BlocksDistribution", "BlocksDistribution")
	if err != nil {
		distributions := c.blocksDistribution()
		c.cache.Add("BlocksDistribution", "BlocksDistribution", 60*60*time.Second, distributions)
		return distributions
	}
	return value.([]*types.DistributionResp)
}

func (c *Controller) blocksDistribution() []*types.DistributionResp {
	minerStatus := c.storage.BlocksDistribution()
	distributions := map[string]*types.DistributionResp{}
	rs := make([]*types.DistributionResp, 0)
	var all uint64
	for _, miner := range minerStatus {
		all += miner.Count
	}
	for _, miner := range minerStatus {
		distribution := &types.DistributionResp{}
		ok, pool := types.Miners.Get(miner.Address)
		if ok {
			distribution.Miner = pool.Name
		}
		block := c.storage.GetLastMinerBlock(miner.Address)
		distribution.Address = miner.Address
		distribution.Blocks = miner.Count
		distribution.LastOrder = block.Order
		distribution.LastTimestamp = block.Timestamp
		distribution.Proportion = blocksProportion(distribution.Blocks, all)
		distribution.ProportionNumber = blocksProportionNumber(distribution.Blocks, all)
		dt, ok := distributions[distribution.Address]
		if ok {
			if dt.LastOrder < distribution.LastOrder {
				dt.LastOrder = distribution.LastOrder
				dt.LastTimestamp = distribution.LastTimestamp
			}
			dt.Blocks += distribution.Blocks
			dt.Proportion = blocksProportion(distribution.Blocks, all)
			distributions[distribution.Address] = dt
		} else {
			distributions[distribution.Address] = distribution
		}
	}
	for _, dt := range distributions {
		rs = append(rs, dt)
	}
	return rs
}

func blocksProportion(blocks, all uint64) string {
	return fmt.Sprintf("%.5f %%", float64(blocks)/float64(all)*100)
}

func blocksProportionNumber(blocks, all uint64) float64 {
	return float64(blocks) / float64(all) * 100
}
