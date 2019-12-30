package client

import (
	"go.etcd.io/etcd/clientv3"
)

// set kv or directory
func (clt *EtcdHRCHYClient) Put(key, value string) error {
	return clt.put(key, value, false)
}

// mustEmpty to confirm the key has not been set
func (clt *EtcdHRCHYClient) put(key string, value string, mustEmpty bool) error {
	key, parentKey, err := clt.ensureKey(key)
	if err != nil {
		return err
	}

	// parentKey should be directory
	// key should not be directory
	cmp := []clientv3.Cmp{
		clientv3.Compare(
			clientv3.Value(parentKey),
			"=",
			clt.dirValue,
		),
	}

	if mustEmpty {
		cmp = append(
			cmp,
			clientv3.Compare(
				clientv3.Version(key),
				"=",
				0,
			),
		)
	} else {
		cmp = append(
			cmp,
			clientv3.Compare(
				clientv3.Value(key),
				"!=",
				clt.dirValue,
			),
		)
	}

	txn := clt.client.Txn(clt.ctx)
	// make sure the parentKey is a directory
	txn.If(
		cmp...,
	).Then(
		clientv3.OpPut(key, value),
	)

	txnResp, err := txn.Commit()
	if err != nil {
		return err
	}

	if !txnResp.Succeeded {
		return ErrorPutKey
	}
	return nil
}
