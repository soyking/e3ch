package client

import (
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
	"strings"
)

const (
	DEFAULT_DIR_VALUE = "etcdv3_dir_$2H#%gRe3*t"
)

var (
	ErrorInvalidRootKey = errors.New("root key should not be empty or end with /")
	ErrorInvalidKey     = errors.New("key should start with /")
	ErrorPutKey         = errors.New("key is not under a directory or key is a directory or key is not empty")
	ErrorKeyNotFound    = errors.New("key has not been set")
	ErrorListKey        = errors.New("can only list a directory")
)

// etcd v3 client with Hierarchy
type EtcdHRCHYClient struct {
	client *clientv3.Client
	ctx    context.Context
	// root key as root directory
	rootKey string
	// special value for directory
	dirValue string
}

func New(clt *clientv3.Client, rootKey string, dirValue ...string) (*EtcdHRCHYClient, error) {
	if !checkRootKey(rootKey) {
		return nil, ErrorInvalidRootKey
	}

	d := DEFAULT_DIR_VALUE
	if len(dirValue) > 0 && dirValue[0] != "" {
		d = dirValue[0]
	}

	return &EtcdHRCHYClient{
		client:   clt,
		rootKey:  rootKey,
		dirValue: d,
		ctx:      context.TODO(),
	}, nil
}

func (clt *EtcdHRCHYClient) EtcdClient() *clientv3.Client {
	return clt.client
}

func (clt *EtcdHRCHYClient) RootKey() string {
	return clt.rootKey
}

func (clt *EtcdHRCHYClient) DirValue() string {
	return clt.dirValue
}

// clone client with new etcdClt, for changing some config of etcdClt
func (clt *EtcdHRCHYClient) Clone(etcdClt *clientv3.Client) *EtcdHRCHYClient {
	return &EtcdHRCHYClient{
		client:   etcdClt,
		rootKey:  clt.rootKey,
		dirValue: clt.dirValue,
		ctx:      context.TODO(),
	}
}

// make sure the rootKey is a directory
func (clt *EtcdHRCHYClient) FormatRootKey() error {
	_, err := clt.client.Put(clt.ctx, clt.rootKey, clt.dirValue)
	return err
}

type Node struct {
	*mvccpb.KeyValue
	IsDir bool `json:"is_dir"`
}

func (clt *EtcdHRCHYClient) isDir(value []byte) bool {
	return string(value) == clt.dirValue
}

func (clt *EtcdHRCHYClient) trimRootKey(key string) string {
	return strings.TrimPrefix(key, clt.rootKey)
}

func (clt *EtcdHRCHYClient) createNode(kv *mvccpb.KeyValue) *Node {
	// remove rootKey prefix
	kv.Key = []byte(clt.trimRootKey(string(kv.Key)))
	return &Node{
		KeyValue: kv,
		IsDir:    clt.isDir(kv.Value),
	}
}
