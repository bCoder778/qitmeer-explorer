package sqldb

import (
	dbtype "github.com/bCoder778/qitmeer-explorer/db/types"
	"github.com/bCoder778/qitmeer-sync/storage/types"
	"github.com/bCoder778/qitmeer-sync/verify/stat"
	"strconv"
)

func (d *DB) BlocksDistribution() []*dbtype.MinerStatus {
	status := []*dbtype.MinerStatus{}
	d.engine.Table(new(types.Block)).Select("address, count(*) as count").Where("`order` > ? and stat in (?, ?)", 0, stat.Block_Confirmed, stat.Block_Unconfirmed).GroupBy("address").Find(&status)
	return status
}


func (d *DB)PackageTime() *dbtype.Package{
	paInfo := &dbtype.Package{}
	rs, err := d.engine.QueryString("select max(block.timestamp-transaction.timestamp) as maxTime,min(block.timestamp-transaction.timestamp) as minTime,sum(block.timestamp-transaction.timestamp)/count(*) as avgTime  from transaction INNER JOIN block on transaction.block_hash = block.hash where transaction.is_coinbase=0; ")
	if err != nil{
		return paInfo
	}
	for _, value := range rs{
		paInfo.AvgTime, _ = strconv.ParseFloat(value["avgTime"], 64)
		paInfo.MaxTime, _ = strconv.ParseInt(value["maxTime"], 10, 64)
		paInfo.MinTime, _ = strconv.ParseInt(value["minTime"], 10, 64)
	}
	return paInfo
}
