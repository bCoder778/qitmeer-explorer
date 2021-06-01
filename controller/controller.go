package controller

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/cache"
	"github.com/bCoder778/qitmeer-explorer/conf"
	"github.com/bCoder778/qitmeer-explorer/controller/qitmeer"
	"github.com/bCoder778/qitmeer-explorer/db"
	"github.com/bCoder778/qitmeer-explorer/db/sqldb"
	"github.com/bCoder778/qitmeer-sync/config"
	"github.com/bCoder778/qitmeer-sync/rpc"
)

type Controller struct {
	storage   db.IDB
	qitmeer   IQitmeer
	cache     *cache.MemCache
	rpcClient *rpc.Client
}

func NewController(conf *conf.Config) (*Controller, error) {
	storage, err := db.ConnectDB(conf)
	if err != nil {
		return nil, err
	}
	rpcCli := rpc.NewClient(&config.Rpc{Host: conf.Rpc.Host, Admin: conf.Rpc.Admin, Password: conf.Rpc.Password})
	_, err = rpcCli.GetBlockCount()
	if err != nil {
		return nil, fmt.Errorf("connect rpc failed, %s", err.Error())
	}
	qitmeer, err := newQitmeer(conf.Qitmeer, storage, rpcCli)
	if err != nil {
		return nil, err
	}
	go qitmeer.StartFindPeer()
	return &Controller{storage: storage, qitmeer: qitmeer, cache: cache.NewMemCache()}, nil
}

func (c *Controller) Close() {
	c.storage.Close()
	c.qitmeer.StopFindPeer()
}

func newQitmeer(c *conf.Qitmeer, storage *sqldb.DB, rpcClient *rpc.Client) (IQitmeer, error) {
	switch c.Version {
	//case "0.9":
	//	return qitmeer.NewQitmeerV0_9(c.Network, storage, storage, rpcClient), nil
	case "0.10":
		return qitmeer.NewQitmeerV0_10(c.Network, storage), nil
	}
	return nil, fmt.Errorf("wrong version of Qitmeer %s", c.Version)
}
