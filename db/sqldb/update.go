package sqldb

import dbtypes "github.com/bCoder778/qitmeer-explorer/db/types"

func (d *DB) UpdatePeer(peer *dbtypes.Peer) error {
	oldPeer := &dbtypes.Peer{}
	ok, err := d.engine.Table(new(dbtypes.Peer)).Where("address = ?", peer.Address).Get(oldPeer)
	if ok {
		_, err = d.engine.Table(new(dbtypes.Peer)).Where("address = ?", peer.Address).Update(map[string]interface{}{
			"find_time": peer.FindTime,
			"other":     peer.Other,
		})
	} else {
		_, err = d.engine.Insert(peer)
	}
	return err
}
