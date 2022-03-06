package client

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

// list a directory
func (clt *EtcdHRCHYClient) Delete(key string) error {
	key, _, err := clt.ensureKey(key)
	if err != nil {
		return err
	}
	// directory start with /
	dir := key + "/"

	txn := clt.client.Txn(clt.ctx)
	// delete the whole dir if it's a directory
	txn.If(
		clientv3.Compare(
			clientv3.Value(key),
			"=",
			clt.dirValue,
		),
	).Then(
		clientv3.OpDelete(key),
		clientv3.OpDelete(dir, clientv3.WithPrefix()),
	).Else(
		clientv3.OpDelete(key),
	)

	_, err = txn.Commit()
	return err
}
