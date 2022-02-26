package api

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/conf"
	"github.com/bCoder778/qitmeer-explorer/controller"
	"github.com/bCoder778/qitmeer-explorer/controller/types"
	"github.com/bCoder778/qitmeer-sync/verify/stat"
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
		GetSub("block", a.queryBLock).
		GetSub("transaction", a.queryTransaction).
		GetSub("address/transaction", a.lastAddressTransactions).
		GetSub("top/address", a.balanceTop).
		GetSub("node", a.nodeList)

	a.rest.AuthRouteSet("api/v1/detail").
		GetSub("block", a.getBlock).
		GetSub("transaction", a.getTransaction)

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

	// 交易所使用，必须保证和原有API一致
	a.rest.AuthRouteSet("api/v1/explorer").
		GetSub("price", a.getPrice).
		GetSub("circulating", a.getCirculating).
		GetSub("circulatingfloat", a.getCirculatingFloat).
		GetSub("max", a.getMax).
		GetSub("maxfloat", a.getMaxFloat)

	// Api V2
	a.rest.AuthRouteSet("api/v2/block").
		GetSub("detail", a.getBlock).
		GetSub("list", a.queryBLock).
		GetSub("txs", a.queryTransactionsByBlockHash).
		GetSub("pending", a.queryBlockPending).
		GetSub("completed", a.queryBlockCompleted)

	a.rest.AuthRouteSet("api/v2/tx").
		GetSub("detail", a.getTransaction).
		GetSub("list", a.queryTransaction).
		GetSub("pending", a.queryTransactionPending).
		GetSub("completed", a.queryTransactionCompleted).
		GetSub("transfer", a.queryTransfer).
		GetSub("coinbase", a.queryCoinbase)

	a.rest.AuthRouteSet("api/v2/token").
		GetSub("list", a.coinIdList).
		GetSub("transfers", a.queryTokenTxs)

	a.rest.AuthRouteSet("api/v2/address").
		GetSub("list", a.balanceTop).
		GetSub("detail", a.addressStatus).
		GetSub("transfers", a.lastAddressTransactions)

	a.rest.AuthRouteSet("api/v2/chain").
		GetSub("search", a.searchV2).
		GetSub("tips", a.tips).
		GetSub("txWaitTime", a.packageTime).
		GetSub("algorithm/list", a.algorithmList).
		GetSub("algorithm/line", a.algorithmLine).
		GetSub("miner", a.blocksDistribution).
		GetSub("circulating", a.getCirculating).
		GetSub("circulatingfloat", a.getCirculatingFloat).
		GetSub("max", a.getMax).
		GetSub("maxfloat", a.getMaxFloat).
		GetSub("price", a.meerPrice).
		GetSub("volume", a.volume)
}

func (a *Api) queryBLock(ct *Context) (interface{}, *Error) {
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

func (a *Api) queryBlockPending(ct *Context) (interface{}, *Error) {
	page, size, err := a.parseListParam(ct)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}

	blocks, err := a.controller.QueryBlockByStatus(page, size, fmt.Sprintf("%v", stat.Block_Unconfirmed))
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	return blocks, nil
}

func (a *Api) queryBlockCompleted(ct *Context) (interface{}, *Error) {
	page, size, err := a.parseListParam(ct)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}

	blocks, err := a.controller.QueryBlockByStatus(page, size, fmt.Sprintf("%v", stat.Block_Confirmed))
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	return blocks, nil
}

func (a *Api) queryTransaction(ct *Context) (interface{}, *Error) {
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

func (a *Api) queryTransactionPending(ct *Context) (interface{}, *Error) {
	page, size, err := a.parseListParam(ct)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	txs, err := a.controller.QueryBlockByStatus(page, size, fmt.Sprintf("%v,%v", stat.TX_Unconfirmed, stat.TX_Memry))
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	return txs, nil
}

func (a *Api) queryTransactionCompleted(ct *Context) (interface{}, *Error) {
	page, size, err := a.parseListParam(ct)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	txs, err := a.controller.QueryBlockByStatus(page, size, fmt.Sprintf("%v", stat.TX_Confirmed))
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	return txs, nil
}

func (a *Api) queryTransfer(ct *Context) (interface{}, *Error) {
	page, size, err := a.parseListParam(ct)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	txs, err := a.controller.QueryTransfer(page, size)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	return txs, nil
}

func (a *Api) queryCoinbase(ct *Context) (interface{}, *Error) {
	page, size, err := a.parseListParam(ct)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	txs, err := a.controller.QueryCoinbase(page, size)
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
	coin := ct.Query["coin"]
	if coin == "" {
		coin = "MEER"
	}
	txs, err := a.controller.LastAddressTransactions(page, size, ct.Query["address"], coin)
	if err != nil {
		return txs, ParseError(err)
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

func (a *Api) getBlock(ct *Context) (interface{}, *Error) {
	var block *types.BlockDetailResp
	var err error
	hash := ct.Query["hash"]
	if hash != "" {
		block, err = a.controller.BlockDetail(hash)
	} else {
		block, err = a.controller.BlockDetail(ct.Query["order"])
	}
	if err != nil {
		return nil, &Error{
			Code:    ERROR_UNKNOWN,
			Message: err.Error(),
		}
	}
	return block, nil
}

func (a *Api) getTransaction(ct *Context) (interface{}, *Error) {
	block, err := a.controller.TransactionDetail(ct.Query["txid"], ct.Query["blockhash"], ct.Query["address"])
	if err != nil {
		return nil, &Error{
			Code:    ERROR_UNKNOWN,
			Message: err.Error(),
		}
	}
	return block, nil
}

func (a *Api) addressStatus(ct *Context) (interface{}, *Error) {
	status, err := a.controller.AddressStatus(ct.Query["address"], ct.Query["coin"])
	if err != nil {
		return status, ParseError(err)
	}
	return status, nil
}

func (a *Api) blocksDistribution(ct *Context) (interface{}, *Error) {
	page, size, err := a.parseListParam(ct)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	distribution := a.controller.BlocksDistribution(page, size)
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

func (a *Api) packageTime(ct *Context) (interface{}, *Error) {
	count := ct.Query["count"]
	packgeInfo := a.controller.PackageTime(count)
	return packgeInfo, nil
}

func (a *Api) coinIdList(ct *Context) (interface{}, *Error) {
	tokens := a.controller.GetCoinIds()
	return tokens, nil
}

func (a *Api) volume(ct *Context) (interface{}, *Error) {
	var before int64
	var err error
	b := ct.Query["before"]
	if len(b) == 0 {
		before = 0
	} else {
		before, err = strconv.ParseInt(b, 10, 64)
		if err != nil {
			return "", &Error{ERROR_UNKNOWN, err.Error()}
		}
	}
	v := a.controller.Volume(before)
	return v, nil
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

func (a *Api) queryTokenTxs(ct *Context) (interface{}, *Error) {
	page, size, err := a.parseListParam(ct)
	coinId := ct.Query["coinId"]
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}

	vs, err := a.controller.QueryTokenTxs(page, size, coinId, "")
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	return vs, nil
}

func (a *Api) queryTransactionsByBlockHash(ct *Context) (interface{}, *Error) {
	page, size, err := a.parseListParam(ct)
	hash := ct.Query["hash"]

	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	rs, err := a.controller.BlockTransactions(hash, size, page)
	if err != nil {
		return nil, &Error{Code: ERROR_UNKNOWN, Message: err.Error()}
	}
	return rs, nil
}

func (a *Api) searchV2(ct *Context) (interface{}, *Error) {
	val := ct.Query["value"]
	if len(val) == 0 {
		return "", &Error{ERROR_PARAM, "condition is required"}
	}
	rs, err := a.controller.SearchV2(val)
	if err != nil {
		return "", ParseError(err)
	}
	return rs, nil
}

func (a *Api) meerPrice(ct *Context) (interface{}, *Error) {
	rs := a.controller.GetMeerPrice()
	return rs, nil
}
