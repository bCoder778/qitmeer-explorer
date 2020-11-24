package controller

import (
	"github.com/bCoder778/qitmeer-explorer/conf"
	"github.com/bCoder778/qitmeer-explorer/db"
)

type Controller struct {
	db db.IDB
}

func NewController(conf *conf.Config) (*Controller, error) {
	db, err := db.ConnectDB(conf)
	if err != nil {
		return nil, err
	}
	return &Controller{db}, nil
}
