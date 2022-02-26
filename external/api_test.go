package external

import (
	"fmt"
	"testing"
)

func TestGetIpInfo(t *testing.T) {
	fmt.Println(GetIpInfo("47.93.20.102"))
}

func TestGetPrice(t *testing.T) {
	rs, err := GetPrice()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rs.LastPrice)
}
