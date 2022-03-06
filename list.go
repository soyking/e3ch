package client

import (
	"strings"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// list a directory
func (clt *EtcdHRCHYClient) List(key string) ([]*Node, error) {
	key, _, err := clt.ensureKey(key)
	if err != nil {
		return nil, err
	}
	// directory start with /
	dir := key + "/"

	txn := clt.client.Txn(clt.ctx)
	// make sure the list key is a directory
	txn.If(
		clientv3.Compare(
			clientv3.Value(key),
			"=",
			clt.dirValue,
		),
	).Then(
		clientv3.OpGet(dir, clientv3.WithPrefix()),
	)

	txnResp, err := txn.Commit()
	if err != nil {
		return nil, err
	}

	if !txnResp.Succeeded {
		return nil, ErrorListKey
	} else {
		if len(txnResp.Responses) > 0 {
			rangeResp := txnResp.Responses[0].GetResponseRange()
			return clt.list(dir, rangeResp.Kvs)
		} else {
			// empty directory
			return []*Node{}, nil
		}
	}
}

// pick key/value under the dir
func (clt *EtcdHRCHYClient) list(dir string, kvs []*mvccpb.KeyValue) ([]*Node, error) {
	nodes := []*Node{}
	for _, kv := range kvs {
		name := strings.TrimPrefix(string(kv.Key), dir)
		if strings.Contains(name, "/") {
			// secondary directory
			continue
		}
		nodes = append(nodes, clt.createNode(kv))
	}
	return nodes, nil
}
