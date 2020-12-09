package sqldb

import (
	dbtype "github.com/bCoder778/qitmeer-explorer/db/types"
	"github.com/bCoder778/qitmeer-sync/storage/types"
	"github.com/bCoder778/qitmeer-sync/verify/stat"
)

func (d *DB) BlocksDistribution() []*dbtype.MinerStatus {
	status := []*dbtype.MinerStatus{}
	d.engine.Table(new(types.Block)).Select("address, count(*) as count").Where("`order` > ? and stat in (?, ?)", 0, stat.Block_Confirmed, stat.Block_Unconfirmed).GroupBy("address").Find(&status)
	return status
}
