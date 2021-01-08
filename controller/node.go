package controller

import "github.com/bCoder778/qitmeer-explorer/controller/types"

func (c *Controller) NodeList() interface{} {
	return c.qitmeer.PeerList()
}

func (c *Controller) Tips() *types.TipsResp {
	tips := &types.TipsResp{
		BlockAvg:          "",
		BlockInterval:     "",
		MainBlockAvg:      "",
		MainBlockInterval: "",
		ConcurrencyRate:   "",
		BlockOrder:        0,
		BlockHeight:       0,
	}
	return tips
}
