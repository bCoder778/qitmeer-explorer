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
		a.controller.Close()
		a.rest.Stop()
	}()
}

func (a *Api) addApi() {
	a.rest.AuthRouteSet("api/v1/list").
		GetSub("block", a.lastBlocks).
		GetSub("transaction", a.lastTransactions).
		GetSub("address/transaction", a.lastAddressTransactions).
		GetSub("top/address", a.balanceTop).
		GetSub("node", a.nodeList)

	a.rest.AuthRouteSet("api/v1/detail").
		GetSub("block", a.blockDetail).
		GetSub("transaction", a.transactionDetail)

	a.rest.AuthRouteSet("api/v1/status").
		GetSub("address", a.addressStatus)

	a.rest.AuthRouteSet("api/v1/blocks").
		GetSub("distribution", a.blocksDistribution)

	a.rest.AuthRouteSet("api/v1/algorithm").
		GetSub("list", a.algorithmList).
		GetSub("line", a.algorithmLine)

	a.rest.AuthRouteSet("api/v1/export").
		GetSub("address/transaction", nil)

	a.rest.AuthRouteSet("api/v1/tips").Get(a.tips)

	// 交易所使用
	a.rest.AuthRouteSet("api/v1/explorer").
		GetSub("price", a.getPrice).
		GetSub("circulating", a.getCirculating).
		GetSub("circulatingfloat", a.getCirculatingFloat).
		GetSub("max", a.getMax).
		GetSub("maxfloat", a.getMaxFloat)

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

func (a *Api) lastTransactions(ct *Context) (interface{}, *Error) {
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

func (a *Api) lastAddressTransactions(ct *Context) (interface{}, *Error) {
	page, size, err := a.parseListParam(ct)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	txs, err := a.controller.LastAddressTransactions(page, size, ct.Query["address"])
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	return txs, nil
}

func (a *Api) balanceTop(ct *Context) (interface{}, *Error) {
	page, size, err := a.parseListParam(ct)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	txs, err := a.controller.BalanceTop(page, size)
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
	block, err := a.controller.TransactionDetail(ct.Query["txid"], "")
	if err != nil {
		return nil, &Error{
			Code:    ERROR_UNKNOWN,
			Message: err.Error(),
		}
	}
	return block, nil
}

func (a *Api) addressStatus(ct *Context) (interface{}, *Error) {
	status, err := a.controller.AddressStatus(ct.Query["address"])
	if err != nil {
		return nil, &Error{
			Code:    ERROR_UNKNOWN,
			Message: err.Error(),
		}
	}
	return status, nil
}

func (a *Api) blocksDistribution(ct *Context) (interface{}, *Error) {
	distribution := a.controller.BlocksDistribution()
	return distribution, nil
}

func (a *Api) algorithmList(ct *Context) (interface{}, *Error) {
	alist := a.controller.AlgorithmList()
	return alist, nil
}

func (a *Api) algorithmLine(ct *Context) (interface{}, *Error) {
	algorithm := ct.Query["algorithm"]
	if len(algorithm) == 0 {
		return nil, &Error{
			Code:    ERROR_UNKNOWN,
			Message: "algorithm is required",
		}
	}
	sec := ct.Query["sec"]
	if len(algorithm) == 0 {
		return nil, &Error{
			Code:    ERROR_UNKNOWN,
			Message: "sec is required",
		}
	}

	iSec, err := strconv.Atoi(sec)
	if len(algorithm) == 0 {
		return nil, &Error{
			Code:    ERROR_UNKNOWN,
			Message: err.Error(),
		}
	}
	alist := a.controller.AlgorithmLine(algorithm, iSec)
	return alist, nil
}

func (a *Api) getPrice(ct *Context) (interface{}, *Error) {
	price, err := a.controller.GetPrice()
	if err != nil {
		return nil, &Error{
			Code:    ERROR_UNKNOWN,
			Message: err.Error(),
		}
	}
	return price, nil
}

func (a *Api) getCirculatingFloat(ct *Context) (interface{}, *Error) {
	pMeer := a.controller.GetCirculatingFloat()
	return pMeer, nil
}

func (a *Api) getCirculating(ct *Context) (interface{}, *Error) {
	pMeer := a.controller.GetCirculating()
	return pMeer, nil
}

func (a *Api) getMax(ct *Context) (interface{}, *Error) {
	pMeer := a.controller.GetMaxPMeer()
	return pMeer, nil
}

func (a *Api) getMaxFloat(ct *Context) (interface{}, *Error) {
	pMeer := a.controller.GetMaxFloatPMeer()
	return pMeer, nil
}

func (a *Api) nodeList(ct *Context) (interface{}, *Error) {
	peer := a.controller.NodeList()
	return peer, nil
}

func (a *Api) tips(ct *Context) (interface{}, *Error) {
	tips := a.controller.Tips()
	return tips, nil
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
