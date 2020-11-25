package controller

import "github.com/bCoder778/qitmeer-explorer/controller/types"

func (c *Controller) BalanceTop(page, size int) (*types.ListResp, error) {
	address, err := c.db.BalanceTop(page, size)
	if err != nil {
		return nil, err
	}
	count, err := c.db.GetAddressCount()
	if err != nil {
		return nil, err
	}
	start := (page - 1) * size
	return &types.ListResp{
		Page:  page,
		Size:  size,
		List:  types.ToAddressRespList(address, uint64(start)),
		Count: count,
	}, nil
}
