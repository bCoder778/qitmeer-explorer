package sqldb

import (
	"fmt"
	"github.com/bCoder778/qitmeer-explorer/conf"
	dbtypes "github.com/bCoder778/qitmeer-explorer/db/types"
	"github.com/bCoder778/qitmeer-sync/storage/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	//"github.com/xormplus/xorm"
	"strings"

	//_ "github.com/lunny/godbc"
	_ "github.com/denisenkom/go-mssqldb"
)

type DB struct {
	engine *xorm.Engine
}

func ConnectMysql(conf *conf.DB) (*DB, error) {
	path := strings.Join([]string{conf.User, ":", conf.Password, "@tcp(", conf.Address, ")/", conf.DBName}, "")
	engine, err := xorm.NewEngine("mysql", path)
	if err != nil {
		return nil, err
	}
	engine.ShowSQL(true)

	if err = engine.Sync2(
		new(dbtypes.Peer),
		new(dbtypes.Location),
	); err != nil {
		return nil, err
	}
	return &DB{engine: engine}, nil
}

func ConnectSqlServer(conf *conf.DB) (*DB, error) {
	path := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", conf.Address, conf.User, conf.Password, conf.DBName)
	engine, err := xorm.NewEngine("mssql", path)

	if err != nil {
		return nil, err
	}
	engine.ShowSQL(false)
	if err = engine.Sync2(
		new(types.Block),
		new(types.Transaction),
		new(types.Vin),
		new(types.Vout),
		new(types.Transfer),
		new(dbtypes.Peer),
		new(dbtypes.Location),
	); err != nil {
		return nil, err
	}

	return &DB{engine}, nil
}

func (d *DB) Close() error {
	return d.engine.Close()
}
