package db

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/conf"
	"github.com/bCoder778/qitmeer-explorer/db/sqldb"
	dbtypes "github.com/bCoder778/qitmeer-explorer/db/types"
	"github.com/bCoder778/qitmeer-sync/storage/types"
)

type IDB interface {
	IGet
	IQuery
	IList
	IStatus
	Close() error
}

type IGet interface {
	GetLastOrder() (uint64, error)
	GetLastHeight() (uint64, error)
	GetLastUnconfirmedOrder() (uint64, error)
	GetTransaction(txId string, blockHash string) (*types.Transaction, error)
	GetTransactionByTxId(txId string) ([]*types.Transaction, error)
	GetTransactionByTxIdBlockHash(txId, hash string) (*types.Transaction, error)
	GetVout(txId string, vout int) (*types.Vout, error)
	GetAllUtxo() float64
	GetConfirmedBlockCount() int64
	GetBlockCount(stat string) (int64, error)
	GetValidBlockCount() (int64, error)
	GetTransactionCount(stat string) (int64, error)
	GetAddressTransactionCount(address, coin string) (int64, error)
	GetBlock(hash string) (*types.Block, error)
	GetBlockByOrder(order uint64) (*types.Block, error)
	GetLastBlock() (*types.Block, error)
	GetAddressCount(coin string) (int64, error)
	GetUsableAmount(address, coinId string, height uint64) (float64, error)
	GetLockedAmount(address, coinId string, height uint64) (float64, error)
	GetUnconfirmedAmount(address, coinId string) (float64, error)
	GetLastMinerBlock(address string) *types.Block
	GetLastAlgorithmBlock(algorithm string, edgeBits int) (*types.Block, error)
	GetTokenTransactionCount(coinId, stat string) (int64, error)
	QueryTransferCount()(int64, error)
	QueryCoinbaseCount()(int64, error)
}

type IQuery interface {
	QueryUnconfirmedTranslateTransaction() ([]types.Transaction, error)
	QueryMemTransaction() ([]types.Transaction, error)
	QueryUnConfirmedOrders() ([]uint64, error)
	QueryTransactions(txId string) ([]types.Transaction, error)
	QueryTransactionsByBlockHash(hash string) ([]types.Transaction, error)
	QueryTransactionVout(txId string) ([]*types.Vout, error)
	QueryTransactionVin(txId string) ([]*types.Vin, error)
	QueryAlgorithmDiffInTime(algorithm string, edgeBits int, max int64, min int64) []*types.Block
	QueryTokens() []string
}

type IList interface {
	LastBlocks(page, size int) ([]*types.Block, error)
	LastTransactions(page, size int) ([]*types.Transaction, error)
	LastAddressTxId(page, size int, address, coin string) ([]string, error)
	BalanceTop(page, size int, coinId string) ([]*dbtypes.Address, error)
	QueryBlock(page, size int, stat string) ([]*types.Block, error)
	QueryTransaction(page, size int, stat string) ([]*types.Transaction, error)
	QueryTokenTransaction(page, size int, coinId, stat string) ([]*types.Vout, error)
	QueryTransfer(page, size int)([]*types.Transaction, error)
	QueryCoinbase(page, size int)([]*types.Transaction, error)
}

type IStatus interface {
	BlocksDistribution(page, size int) []*dbtypes.MinerStatus
	PackageTime(count int) *dbtypes.Package
}

func ConnectDB(setting *conf.Config) (*sqldb.DB, error) {
	var (
		db  *sqldb.DB
		err error
	)
	switch setting.DB.DBType {
	case "mysql":
		if db, err = sqldb.ConnectMysql(setting.DB); err != nil {
			return nil, fmt.Errorf("failed to connect mysql, error:%v", err)
		}
	case "sqlserver":
		if db, err = sqldb.ConnectSqlServer(setting.DB); err != nil {
			return nil, fmt.Errorf("failed to connect mysql, error:%v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported database %s", setting.DB.DBType)
	}
	return db, nil
}
