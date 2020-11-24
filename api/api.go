package api

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/conf"
	"github.com/bCoder778/qitmeer-explorer/controller"
	"os"
	"os/signal"
	"strconv"
)

type Api struct {
	rest       *RestApi
	controller *controller.Controller
}

func NewApi(conf *conf.Config) (*Api, error) {
	controller, err := controller.NewController(conf)
	if err != nil {
		return nil, err
	}
	return &Api{
		rest:       NewRestApi(conf.Api.Listen),
		controller: controller,
	}, nil
}

func (a *Api) Run() error {
	a.listenInterrupt()
	a.addApi()
	return a.rest.Start()
}

func (a *Api) listenInterrupt() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)

	go func() {
		<-c
		a.rest.Stop()
	}()
}

func (a *Api) addApi() {
	a.rest.AuthRouteSet("api/v1/list").
		GetSub("block", a.lastBlocks).
		GetSub("transaction", a.lastTransaction).
		GetSub("address", a.maxBalanceAddress)

	a.rest.AuthRouteSet("api/v1/detail").
		GetSub("block", a.blockDetail).
		GetSub("transaction", a.transactionDetail)
}

func (a *Api) lastBlocks(ct *Context) (interface{}, *Error) {
	page, size, err := a.parseListParam(ct)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	blocks, err := a.controller.LastBlocks(page, size)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	return blocks, nil
}

func (a *Api) lastTransaction(ct *Context) (interface{}, *Error) {
	page, size, err := a.parseListParam(ct)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	txs, err := a.controller.LastTransactions(page, size)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	return txs, nil
}

func (a *Api) maxBalanceAddress(ct *Context) (interface{}, *Error) {
	page, size, err := a.parseListParam(ct)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	txs, err := a.controller.MaxBalanceAddress(page, size)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	return txs, nil
}

func (a *Api) blockDetail(ct *Context) (interface{}, *Error) {
	block, err := a.controller.BlockDetail(ct.Query["hash"])
	if err != nil {
		return nil, &Error{
			Code:    ERROR_UNKNOWN,
			Message: err.Error(),
		}
	}
	return block, nil
}

func (a *Api) transactionDetail(ct *Context) (interface{}, *Error) {
	block, err := a.controller.TransactionDetail(ct.Query["txid"])
	if err != nil {
		return nil, &Error{
			Code:    ERROR_UNKNOWN,
			Message: err.Error(),
		}
	}
	return block, nil
}

func (a *Api) parseListParam(ct *Context) (int, int, error) {
	page, err := strconv.Atoi(ct.Query["page"])
	if err != nil {
		return 0, 0, fmt.Errorf("page is required")
	}
	size, err := strconv.Atoi(ct.Query["size"])
	if err != nil {
		return 0, 0, fmt.Errorf("size is required")
	}
	if page < 0 {
		page = 1
	}
	if size < 0 {
		size = 10
	}
	if size > 500 {
		size = 500
	}
	return page, size, nil
}
