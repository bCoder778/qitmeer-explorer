package controller

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	"sort"
	"time"
)

func (c *Controller) BlocksDistribution(page, size int) *types.DistributionsResp {
	key := fmt.Sprintf("'BlocksDistribution'-%d-%d", page, size)
	value, err := c.cache.Value("BlocksDistribution", key)
	if err != nil {
		distributions := c.blocksDistribution(page, size)
		c.cache.Add("BlocksDistribution", key, 60*time.Second, distributions)
		return distributions
	}
	return value.(*types.DistributionsResp)
}

func (c *Controller) blocksDistribution(page, size int) *types.DistributionsResp {
	minerStatus := c.storage.BlocksDistribution(page, size)
	count := c.storage.BlocksDistributionCount()
	allBlocks := c.storage.BlocksCount()
	distributions := map[string]*types.DistributionResp{}
	rs := make([]*types.DistributionResp, 0)

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
		distribution.Proportion = blocksProportion(distribution.Blocks, uint64(allBlocks))
		distribution.ProportionNumber = blocksProportionNumber(distribution.Blocks, uint64(allBlocks))
		dt, ok := distributions[distribution.Address]
		if ok {
			if dt.LastOrder < distribution.LastOrder {
				dt.LastOrder = distribution.LastOrder
				dt.LastTimestamp = distribution.LastTimestamp
			}
			dt.Blocks += distribution.Blocks
			dt.Proportion = blocksProportion(distribution.Blocks, uint64(allBlocks))
			distributions[distribution.Address] = dt
		} else {
			distributions[distribution.Address] = distribution
		}
	}
	for _, dt := range distributions {
		rs = append(rs, dt)
	}
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Blocks > rs[j].Blocks
	})
	return &types.DistributionsResp{
		Page:  page,
		Size:  size,
		List:  rs,
		Count: count,
	}
}

func blocksProportion(blocks, all uint64) string {
	return fmt.Sprintf("%.5f %%", float64(blocks)/float64(all)*100)
}

func blocksProportionNumber(blocks, all uint64) float64 {
	return float64(blocks) / float64(all) * 100
}
