package external

import (
	"fmt"
	"testing"
)

func TestGetIpInfo(t *testing.T) {
	fmt.Println(GetIpInfo("47.93.20.102"))
}
