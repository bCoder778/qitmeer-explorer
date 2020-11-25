package controller

import (
	types2 "github.com/Qitmeer/qitmeer/core/types"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
)

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

func (c *Controller) AddressStatus(address string) (*types.AddressStatusResp, error) {
	usable, err := c.db.GetUsableAmount(address)
	if err != nil {
		return nil, err
	}
	locked, err := c.db.GetLockedAmount(address)
	if err != nil {
		return nil, err
	}
	usable = types2.Amount(uint64(usable)).ToCoin()
	locked = types2.Amount(uint64(locked)).ToCoin()
	return &types.AddressStatusResp{
		Address: address,
		Balance: usable + locked,
		Usable:  usable,
		Locked:  locked,
	}, nil
}
