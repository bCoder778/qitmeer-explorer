package controller

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/conf"
	"github.com/bCoder778/qitmeer-explorer/controller/qitmeer"
	"github.com/bCoder778/qitmeer-explorer/db"
)

type Controller struct {
	storage db.IDB
	qitmeer IQitmeer
}

func NewController(conf *conf.Config) (*Controller, error) {
	storage, err := db.ConnectDB(conf)
	if err != nil {
		return nil, err
	}
	qitmeer, err := newQitmeer(conf.Qitmeer, storage)
	if err != nil {
		return nil, err
	}
	return &Controller{storage: storage, qitmeer: qitmeer}, nil
}

func newQitmeer(c *conf.Qitmeer, storage db.IDB) (IQitmeer, error) {
	switch c.Version {
	case "0.9":
		return qitmeer.NewQitmeerV0_9(c.Network, storage), nil
	case "0.10":
		return qitmeer.NewQitmeerV0_10(c.Network, storage), nil
	}
	return nil, fmt.Errorf("wrong version of Qitmeer %s", c.Version)
}
