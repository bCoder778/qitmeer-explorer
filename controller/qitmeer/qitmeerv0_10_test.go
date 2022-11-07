package qitmeer

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/rpc"
	"testing"
)

func TestQitmeerV0_10_PeerList(t *testing.T) {
	client := rpc.NewClient(node_rpc_host, node_rpc_user, node_rpc_pas)
	peers, err := client.GetNodeList()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(peers)
}
