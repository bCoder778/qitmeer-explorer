package qitmeer

import (
	"errors"
	"github.com/Qitmeer/qitmeer/core/blockdag"
	"github.com/Qitmeer/qitmeer/core/message"
	"github.com/Qitmeer/qitmeer/core/protocol"
	"github.com/Qitmeer/qitmeer/p2p/connmgr"
	"github.com/Qitmeer/qitmeer/p2p/peer"
	"github.com/Qitmeer/qitmeer/params"
	"github.com/bCoder778/qitmeer-explorer/db/types"
	"github.com/bCoder778/qitmeer-sync/rpc"
	"net"
	"strings"
	"time"
)

const (
	peerFindInterval   int64 = 60 * 60
	defaultNodeTimeout       = time.Second * 10
)

var seedList = []string{
	"121.196.55.29",
	"121.196.28.213",
	"121.196.54.163",
	"47.114.183.16",
	"47.114.184.240",
	"47.57.148.48",
	"120.79.138.225",
	"121.42.12.225",
	"47.105.104.177",
	"118.190.97.29",
	"47.111.233.176",
	"47.241.77.114",
	"47.253.43.133",
	"47.242.150.246",
	"103.231.255.228",
	"8.210.82.0",
	"47.102.129.83",
	"8.210.117.240",
	"47.242.21.98",
}

type IPeerDB interface {
	GetPeer(address string) (*types.Peer, error)
	QueryPeers() []*types.Peer
	UpdatePeer(peer *types.Peer) error
}

type PeerManager struct {
	activeNetParams *params.Params
	rpcClient       *rpc.Client
	db              IPeerDB
	stop            chan struct{}
}

func NewPeerManager(network string, rpcClient *rpc.Client, db IPeerDB) *PeerManager {
	actParams := &params.MainNetParams
	switch network {
	case "mainnet":
	case "testnet":
		actParams = &params.TestNetParams
	case "mixnet":
		actParams = &params.MixNetParams
	}

	return &PeerManager{
		activeNetParams: actParams,
		stop:            make(chan struct{}),
		rpcClient:       rpcClient,
		db:              db,
	}
}

func (p *PeerManager) Close() error {
	return nil
}

func (p *PeerManager) AddPeer(peer *types.Peer) bool {
	oldPeer, _ := p.db.GetPeer(peer.Address)
	if oldPeer.Address != "" {
		if peer.FindTime < oldPeer.FindTime {
			return true
		}
	}
	p.db.UpdatePeer(peer)
	return true
}

func (p *PeerManager) AddPeers(peers []*types.Peer) {
	for _, peer := range peers {
		p.AddPeer(peer)
	}
}

func (p *PeerManager) Peers() []*types.Peer {
	return p.db.QueryPeers()
}

func (p *PeerManager) Find() {
	for _, seed := range seedList {
		p.AddPeer(&types.Peer{Address: seed})
	}

	t := time.NewTicker(time.Second * 60 * 10)
	defer t.Stop()

	for {
		select {
		case <-p.stop:
			return
		case <-t.C:
			if ps, _ := p.rpcClient.GetPeerInfo(); len(ps) != 0 {
				for _, peer := range ps {
					p.AddPeer(&types.Peer{Address: peer.Addr})
				}
			}
			peers := p.getFindPeer()
			if len(peers) != 0 {
				for _, peer := range peers {
					select {
					case <-p.stop:
						return
					default:
						p.creepOne(net.ParseIP(Ip(peer.Address)))
						peer.FindTime = time.Now().Unix()
						p.AddPeer(peer)
					}
				}
			}
		}
	}
}

func (p *PeerManager) getFindPeer() []*types.Peer {
	now := time.Now().Unix()

	rs := []*types.Peer{}
	peers := p.db.QueryPeers()
	for _, peer := range peers {
		if peer.FindTime+peerFindInterval < now {
			rs = append(rs, peer)
		}
	}
	return rs
}

func (p *PeerManager) creepOne(ip net.IP) error {

	onaddr := make(chan struct{})
	verack := make(chan struct{})

	newestGSFunc := func() (gs *blockdag.GraphState, err error) {
		gs = blockdag.NewGraphState()
		gs.GetTips().Add(p.activeNetParams.GenesisHash)
		gs.SetTotal(1)
		return gs, err
	}

	onAddrFunc := func(peer *peer.Peer, msg *message.MsgAddr) {
		for _, addr := range msg.AddrList {
			p.AddPeer(&types.Peer{
				Address: addr.IP.String(),
			})
		}
		onaddr <- struct{}{}
	}

	onVerAckFunc := func(p *peer.Peer, msg *message.MsgVerAck) {
		verack <- struct{}{}
	}

	messageListener := peer.MessageListeners{
		OnAddr:   onAddrFunc,
		OnVerAck: onVerAckFunc,
	}

	peerConfig := peer.Config{
		NewestGS:          newestGSFunc,
		UserAgentName:     "qitmeer-seeder",
		UserAgentVersion:  "0.3.1",
		UserAgentComments: []string{"qitmeer", "seeder"},
		ChainParams:       p.activeNetParams,
		DisableRelayTx:    true,
		Services:          protocol.Full,
		ProtocolVersion:   protocol.ProtocolVersion,
		Listeners:         messageListener,
	}
	host := net.JoinHostPort(ip.String(),
		p.activeNetParams.DefaultPort)
	peer, err := peer.NewOutboundPeer(&peerConfig, host)
	if err != nil {
	}
	conn, err := net.DialTimeout("tcp", peer.Addr(),
		defaultNodeTimeout)
	if err != nil {
		return err
	}
	c := connmgr.NewConnReq()
	c.SetConn(conn)
	peer.AssociateConnection(c)

	// Wait for the verack message or timeout in case of
	// failure.
	select {
	case <-verack:

		// Ask peer for some addresses.
		peer.QueueMessage(message.NewMsgGetAddr(), nil)

	case <-time.After(defaultNodeTimeout):
		peer.Disconnect()
		return errors.New("verack time out")
	}

	select {
	case <-onaddr:

	case <-time.After(defaultNodeTimeout):
		peer.Disconnect()
		return errors.New("onaddr time out")
	}
	peer.Disconnect()
	return nil
}

func Ip(address string) string {
	addrs := strings.Split(address, ":")
	if len(addrs) == 2 {
		return addrs[0]
	}
	return address
}
