package client

import "github.com/coreos/etcd/clientv3"

const (
	COMPARE_BIGGER = ">"
	COMPARE_EQUAL  = "="
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

	// check the key-value to ensure it's not empty or not dir
	// it is not safe but etcd v3 doesn't support != in compare
	// and conditions in if can only be combined by &&
	resp, err := clt.client.Get(clt.ctx, key)
	if err != nil {
		return err
	}

	if len(resp.Kvs) > 0 {
		if mustEmpty {
			return ErrorKeyExist
		}

		kv := resp.Kvs[0]
		if string(kv.Value) == clt.dirValue {
			return ErrorPutDir
		}
	}

	txn := clt.client.Txn(clt.ctx)
	// make sure the parentKey is a directory
	txn.If(
		clientv3.Compare(
			clientv3.Value(parentKey),
			"=",
			clt.dirValue,
		),
	).Then(
		clientv3.OpPut(key, value),
	)

	txnResp, err := txn.Commit()
	if err != nil {
		return err
	}

	if !txnResp.Succeeded {
		return ErrorKeyParent
	}
	return nil
}
