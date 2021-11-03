package controller

import (
	"fmt"
	"time"
)


func (c *Controller) Volume(before int64) float64 {
	key := fmt.Sprintf("%d", before)
	value, err := c.cache.Value("Volume", key)
	if err != nil {
		v := c.volume(before)
		c.cache.Add("Volume", key, 10*time.Second, v)
		return v
	}
	return value.(float64)
}

func (c *Controller) volume(before int64) float64 {
	v := c.storage.GetChainVolume(before)
	return float64(-v) / 1e8
}
