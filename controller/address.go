package controller

import (
	"fmt"
	types2 "github.com/Qitmeer/qitmeer/core/types"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	"time"
)

func (c *Controller) BalanceTop(page, size int) (*types.ListResp, error) {
	key := fmt.Sprintf("%d-%d", page, size)
	value, err := c.cache.Value("BalanceTop", key)
	if err != nil {
		balances, err := c.balanceTop(page, size)
		if err != nil {
			return nil, err
		}
		c.cache.Add("BalanceTop", key, 6*60*60*time.Second, balances)
		return balances, nil
	}
	return value.(*types.ListResp), nil
}

func (c *Controller) balanceTop(page, size int) (*types.ListResp, error) {
	address, err := c.storage.BalanceTop(page, size)
	if err != nil {
		return nil, err
	}
	count, err := c.storage.GetAddressCount()
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
	value, err := c.cache.Value("AddressStatus", address)
	if err != nil {
		status, err := c.addressStatus(address)
		if err != nil {
			return nil, err
		}
		c.cache.Add("AddressStatus", address, 60*time.Second, status)
		return status, nil
	}
	return value.(*types.AddressStatusResp), nil
}

func (c *Controller) addressStatus(address string) (*types.AddressStatusResp, error) {
	usable, err := c.storage.GetUsableAmount(address)
	if err != nil {
		return nil, err
	}
	locked, err := c.storage.GetLockedAmount(address)
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
