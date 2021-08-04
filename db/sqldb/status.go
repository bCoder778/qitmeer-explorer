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
	maxRs, err := d.engine.QueryString("select block.timestamp-transaction.timestamp as minTime, tx_id, block_hash  from transaction INNER JOIN block on transaction.block_hash = block.hash where transaction.is_coinbase=0 and block.timestamp-transaction.timestamp > 0 order by block.timestamp-transaction.timestamp desc limit 1;")
	if err != nil{
		return paInfo
	}
	max := &dbtype.TimeInfo{
		WaitTime:  0,
		BlockHash: "",
		TxId:      "",
	}
	for _, value := range maxRs{
		max.WaitTime, _ = strconv.ParseInt(value["maxTime"], 10, 64)
		max.BlockHash, _ = value["block_hash"]
		max.TxId, _ = value["tx_id"]
	}

	minRs, err := d.engine.QueryString("select block.timestamp-transaction.timestamp as minTime, tx_id, block_hash  from transaction INNER JOIN block on transaction.block_hash = block.hash where transaction.is_coinbase=0 and block.timestamp-transaction.timestamp > 0 order by block.timestamp-transaction.timestamp  limit 1;")
	if err != nil{
		return paInfo
	}
	min := &dbtype.TimeInfo{
		WaitTime:  0,
		BlockHash: "",
		TxId:      "",
	}
	for _, value := range minRs{
		min.WaitTime, _ = strconv.ParseInt(value["minTime"], 10, 64)
		min.BlockHash, _ = value["block_hash"]
		min.TxId, _ = value["tx_id"]
	}

	avgRs, err := d.engine.QueryString("select sum(block.timestamp-transaction.timestamp) as sumTime, count(*) as count,sum(block.timestamp-transaction.timestamp)/count(*) as avgTime  from transaction INNER JOIN block on transaction.block_hash = block.hash where transaction.is_coinbase=0  and block.timestamp-transaction.timestamp>0;")
	if err != nil{
		return paInfo
	}
	for _, value := range avgRs{
		paInfo.AvgTime, _ = strconv.ParseFloat(value["avgTime"], 64)
		paInfo.SumTime, _ = strconv.ParseInt(value["sumTime"], 10, 64)
		paInfo.TxCount, _ = strconv.ParseInt(value["count"], 10, 64)
	}
	paInfo.MaxInfo = max
	paInfo.MinInfo = min
	return paInfo
}
