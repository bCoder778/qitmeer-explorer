package sqldb

import (
	"fmt"
	dbtype "github.com/bCoder778/qitmeer-explorer/db/types"
	"github.com/bCoder778/qitmeer-sync/storage/types"
	"github.com/bCoder778/qitmeer-sync/verify/stat"
	"strconv"
	"time"
)

func (d *DB) BlocksCount() int64 {
	count, _ := d.engine.Table(new(types.Block)).Where("`order` > ? and stat in (?, ?)", 0, stat.Block_Confirmed, stat.Block_Unconfirmed).Count()
	return count
}

func (d *DB) BlocksDistribution(page, size int) []*dbtype.MinerStatus {
	page -= 1
	start := page * size
	status := []*dbtype.MinerStatus{}
	d.engine.Table(new(types.Block)).Select("address, count(*) as count").Where("`order` > ? and stat in (?, ?)", 0, stat.Block_Confirmed, stat.Block_Unconfirmed).GroupBy("address").Desc("count").Limit(size, start).Find(&status)
	return status
}

func (d *DB) BlocksDistributionCount() int64 {
	sql := fmt.Sprintf("select count(t.address) as count from (select address from block where block.order>%d and stat in (%d, %d)  group by  address) t", 0, stat.Block_Confirmed, stat.Block_Unconfirmed)
	rs, _ := d.engine.QueryString(sql)
	for _, value := range rs {
		sCount := value["count"]
		count, _ := strconv.ParseUint(sCount, 10, 64)
		return int64(count)
	}
	return 0
}

func (d *DB) AvgBlocks(sec uint64) float64 {
	now := uint64(time.Now().Unix())
	avgTimes := float64(0)
	sql := fmt.Sprintf("select %d/count(*) as avgTimes from block where timestamp > %d", sec, now-sec)
	maxRs, err := d.engine.QueryString(sql)
	if err != nil {
		return 0
	}
	for _, value := range maxRs {
		avgTimes, err = strconv.ParseFloat(value["avgTimes"], 64)
	}
	return avgTimes
}

func (d *DB) PackageTime(count int) *dbtype.Package {
	paInfo := &dbtype.Package{}
	sql := ""
	if count != 0 {
		sql = fmt.Sprintf("select block.address, block.peer_id, block.timestamp-tx.timestamp as maxTime, tx_id, block_hash  from (select transaction.tx_id, transaction.block_hash,transaction.timestamp,transaction.is_coinbase from transaction INNER JOIN block on transaction.block_hash = block.hash where transaction.is_coinbase=0  and block.timestamp-transaction.timestamp>0 and transaction.duplicate = 0 order by transaction.timestamp desc limit %d) as tx INNER JOIN block on tx.block_hash = block.hash where tx.is_coinbase=0 and block.timestamp-tx.timestamp > 0 order by block.timestamp-tx.timestamp desc limit 1;", count)
	} else {
		sql = "select block.address, block.peer_id, block.timestamp-transaction.timestamp as maxTime, tx_id, block_hash  from transaction INNER JOIN block on transaction.block_hash = block.hash where transaction.is_coinbase=0 and block.timestamp-transaction.timestamp > 0 and transaction.duplicate = 0 order by block.timestamp-transaction.timestamp desc limit 1;"
	}
	maxRs, err := d.engine.QueryString(sql)
	if err != nil {
		return paInfo
	}
	max := &dbtype.TimeInfo{
		WaitTime:  "",
		BlockHash: "",
		TxId:      "",
	}
	for _, value := range maxRs {
		allsec, _ := strconv.ParseInt(value["maxTime"], 10, 64)
		hour, minute, sec := resolveTime(allsec)
		max.WaitSec = allsec
		max.WaitTime = fmt.Sprintf("%02dh:%02dm:%02ds", hour, minute, sec)
		max.BlockHash, _ = value["block_hash"]
		max.TxId, _ = value["tx_id"]
		max.Address = value["address"]
		max.PeerId = value["peer_id"]
	}

	if count != 0 {
		sql = fmt.Sprintf("select block.address, block.peer_id, block.timestamp-tx.timestamp as minTime, tx_id, block_hash  from (select transaction.tx_id, transaction.block_hash,transaction.timestamp,transaction.is_coinbase from transaction INNER JOIN block on transaction.block_hash = block.hash where transaction.is_coinbase=0  and block.timestamp-transaction.timestamp>0 and transaction.duplicate = 0 order by transaction.timestamp desc limit %d) as tx INNER JOIN block on tx.block_hash = block.hash where tx.is_coinbase=0 and block.timestamp-tx.timestamp > 0 order by block.timestamp-tx.timestamp limit 1;", count)
	} else {
		sql = "select block.address, block.peer_id, block.timestamp-transaction.timestamp as minTime, tx_id, block_hash  from transaction INNER JOIN block on transaction.block_hash = block.hash where transaction.is_coinbase=0 and block.timestamp-transaction.timestamp > 0 and transaction.duplicate = 0 order by block.timestamp-transaction.timestamp limit 1;"
	}
	minRs, err := d.engine.QueryString(sql)
	if err != nil {
		return paInfo
	}
	min := &dbtype.TimeInfo{
		WaitTime:  "",
		BlockHash: "",
		TxId:      "",
	}
	for _, value := range minRs {
		allsec, _ := strconv.ParseInt(value["minTime"], 10, 64)
		hour, minute, sec := resolveTime(allsec)
		min.WaitSec = allsec
		min.WaitTime = fmt.Sprintf("%02dh:%02dm:%02ds", hour, minute, sec)
		min.BlockHash, _ = value["block_hash"]
		min.TxId, _ = value["tx_id"]
		min.Address = value["address"]
		min.PeerId = value["peer_id"]
	}

	if count == 0 {
		sql = "select sum(block.timestamp-transaction.timestamp) as sumTime, count(*) as count,sum(block.timestamp-transaction.timestamp)/count(*) as avgTime  from transaction INNER JOIN block on transaction.block_hash = block.hash where transaction.is_coinbase=0  and block.timestamp-transaction.timestamp>0 and transaction.duplicate = 0;"
	} else {
		sql = fmt.Sprintf("select sum(block.timestamp-tx.timestamp) as sumTime,count(*) as count, sum(block.timestamp-tx.timestamp)/count(*) as avgTime  from (select transaction.tx_id, transaction.block_hash,transaction.timestamp,transaction.is_coinbase from transaction INNER JOIN block on transaction.block_hash = block.hash where transaction.is_coinbase=0  and block.timestamp-transaction.timestamp>0 and transaction.duplicate = 0 order by transaction.timestamp desc limit %d) as tx  INNER JOIN block on tx.block_hash = block.hash where tx.is_coinbase=0  and block.timestamp-tx.timestamp>0", count)
	}
	avgRs, err := d.engine.QueryString(sql)
	if err != nil {
		return paInfo
	}
	for _, value := range avgRs {
		paInfo.AvgSeconds, _ = strconv.ParseFloat(value["avgTime"], 64)
		hour, minute, sec := resolveTime(int64(paInfo.AvgSeconds))
		paInfo.AvgTime = fmt.Sprintf("%02dh:%02dm:%02ds", hour, minute, sec)
		paInfo.SumSec, _ = strconv.ParseInt(value["sumTime"], 10, 64)
		paInfo.TxCount, _ = strconv.ParseInt(value["count"], 10, 64)
	}
	paInfo.MaxInfo = max
	paInfo.MinInfo = min

	if count == 0 {
		sql = "select block.address, block.peer_id, sum(block.timestamp-transaction.timestamp) as sumTime, count(*) as count,sum(block.timestamp-transaction.timestamp)/count(*) as avgTime  from transaction INNER JOIN block on transaction.block_hash = block.hash where transaction.is_coinbase=0  and block.timestamp-transaction.timestamp>0 and transaction.duplicate = 0 group by block.address order by avgTime desc;"
	} else {
		sql = fmt.Sprintf("select block.address, block.peer_id, sum(block.timestamp-tx.timestamp) as sumTime,count(*) as count, sum(block.timestamp-tx.timestamp)/count(*) as avgTime  from (select transaction.tx_id, transaction.block_hash,transaction.timestamp,transaction.is_coinbase from transaction INNER JOIN block on transaction.block_hash = block.hash where transaction.is_coinbase=0  and block.timestamp-transaction.timestamp>0 and transaction.duplicate = 0 order by transaction.timestamp desc limit %d) as tx  INNER JOIN block on tx.block_hash = block.hash where tx.is_coinbase=0  and block.timestamp-tx.timestamp>0 group by block.peer_id order by avgTime desc;", count)
	}
	minerRs, err := d.engine.QueryString(sql)
	if err != nil {
		return paInfo
	}
	miners := []dbtype.MinerInfo{}
	for _, value := range minerRs {
		miner := dbtype.MinerInfo{
			Address:    value["address"],
			Miner:      "",
			TxCount:    0,
			AvgTime:    "",
			AvgSeconds: 0,
		}
		miner.Address = value["address"]
		miner.AvgSeconds, _ = strconv.ParseFloat(value["avgTime"], 64)
		hour, minute, sec := resolveTime(int64(miner.AvgSeconds))
		miner.AvgTime = fmt.Sprintf("%02dh:%02dm:%02ds", hour, minute, sec)
		miner.SumSec, _ = strconv.ParseInt(value["sumTime"], 10, 64)
		miner.TxCount, _ = strconv.ParseInt(value["count"], 10, 64)
		miner.PeerId = value["peer_id"]
		miners = append(miners, miner)
	}
	paInfo.Miners = miners

	return paInfo
}

func resolveTime(seconds int64) (hour, minute, second int64) {
	var day = seconds / (24 * 3600)
	hour = (seconds - day*3600*24) / 3600
	minute = (seconds - day*24*3600 - hour*3600) / 60
	second = seconds - day*24*3600 - hour*3600 - minute*60
	return
}

func (d *DB) GetChainVolume(before int64) int64 {
	var rs int64
	if before == 0 {
		rs, _ = d.engine.Table(new(types.Transfer)).Where("is_coinbase = ? and tx_id <> ? and transfer.change < ? ",
			0, "d67673726970e06d866052188d1da4711570c343ea22ad334cac14d6de958599", 0).SumInt(new(types.Transfer), "change")
	} else {
		end := time.Now().Unix()
		start := end - before
		rs, _ = d.engine.Where("timestamp between ? and ? and is_coinbase = ? and tx_id <> ? and transfer.change < ? ",
			start, end, 0, "d67673726970e06d866052188d1da4711570c343ea22ad334cac14d6de958599", 0).SumInt(new(types.Transfer), "change")
	}
	return rs
}
