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
	GetLastUnconfirmedOrder() (uint64, error)
	GetTransaction(txId string, blockHash string) (*types.Transaction, error)
	GetTransactionByTxId(txId string) ([]*types.Transaction, error)
	GetVout(txId string, vout int) (*types.Vinout, error)
	GetAllUtxo() float64
	GetConfirmedBlockCount() int64
	GetBlockCount() (int64, error)
	GetTransactionCount() (int64, error)
	GetAddressTransactionCount(address string) (int64, error)
	GetBlock(hash string) (*types.Block, error)
	GetAddressCount() (int64, error)
	GetUsableAmount(address string) (float64, error)
	GetLockedAmount(address string) (float64, error)
	GetLastMinerBlock(address string) *types.Block
	GetLastAlgorithmBlock(algorithm string, edgeBits int) (*types.Block, error)
}

type IQuery interface {
	QueryUnconfirmedTranslateTransaction() ([]types.Transaction, error)
	QueryMemTransaction() ([]types.Transaction, error)
	QueryUnConfirmedOrders() ([]uint64, error)
	QueryTransactions(txId string) ([]types.Transaction, error)
	QueryTransactionsByBlockHash(hash string) ([]types.Transaction, error)
	QueryTransactionVout(txId string) ([]*types.Vinout, error)
	QueryTransactionVin(txId string) ([]*types.Vinout, error)
}

type IList interface {
	LastBlocks(page, size int) ([]*types.Block, error)
	LastTransactions(page, size int) ([]*types.Transaction, error)
	LastAddressTxId(page, size int, address string) ([]string, error)
	BalanceTop(page, size int) ([]*dbtypes.Address, error)
}

type IStatus interface {
	BlocksDistribution() []*dbtypes.MinerStatus
}

func ConnectDB(setting *conf.Config) (IDB, error) {
	var (
		db  IDB
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
