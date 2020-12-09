package controller

import "github.com/bCoder778/qitmeer-explorer/controller/types"

func (c *Controller) AlgorithmList() []*types.AlgorithmResp {
	return c.qitmeer.AlgorithmList()
}
