package client

import "github.com/coreos/etcd/clientv3"

const (
	COMPARE_BIGGER = ">"
	COMPARE_EQUAL  = "="
)

// set kv or directory
func (clt *EtcdHRCHYClient) Put(key, value string) error {
	success, err := clt.put(key, value, COMPARE_BIGGER)
	if err != nil {
		return err
	}

	if !success {
		return ErrorPutKey
	}

	return nil
}

// with cond to confirm key has or has not been set, return (isTxnSuccess, error)
func (clt *EtcdHRCHYClient) put(key string, value, cond string) (bool, error) {
	key, parentKey, err := clt.ensureKey(key)
	if err != nil {
		return false, err
	}

	txn := clt.client.Txn(clt.ctx)
	// make sure the parentKey is a directory and key has been created
	txn.If(
		clientv3.Compare(
			clientv3.Value(parentKey),
			"=",
			clt.dirValue,
		),
		clientv3.Compare(
			clientv3.Version(key),
			cond,
			0,
		),
	).Then(
		clientv3.OpPut(key, value),
	)

	txnResp, err := txn.Commit()
	if err != nil {
		return false, err
	}

	return txnResp.Succeeded, nil
}
