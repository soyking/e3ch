package client

import (
	"os"
	"testing"

	"go.etcd.io/etcd/clientv3"
	. "gopkg.in/check.v1"
)

var (
	TEST_ETCD_ADDR = "127.0.0.1:2379"
	TEST_ROOT_KEY  = "e3ch_test"
)

func Test(t *testing.T) { TestingT(t) }

var (
	client *EtcdHRCHYClient
)

func init() {
	// FIXME
	if addr, ok := os.LookupEnv("TEST_ETCD_ADDR"); ok {
		TEST_ETCD_ADDR = addr
	}

	clt, err := clientv3.New(clientv3.Config{Endpoints: []string{TEST_ETCD_ADDR}})
	if err != nil {
		panic(err)
	}

	client, err = New(clt, TEST_ROOT_KEY)
	if err != nil {
		panic(err)
	}

	err = client.FormatRootKey()
	if err != nil {
		panic(err)
	}

	Suite(&ClientSuite{})
	Suite(&PutSuite{})
	Suite(&GetSuite{})
	Suite(&ListSuite{})
	Suite(&DeleteSuite{})
	Suite(&CreateSuite{})
	Suite(&AuthSuite{})
}

type ClientSuite struct{}

func (s *ClientSuite) TestNewClient(c *C) {
	clt, err := clientv3.New(clientv3.Config{Endpoints: []string{TEST_ETCD_ADDR}})
	if err != nil {
		c.Error(err)
	}

	_, err = New(clt, "")
	c.Assert(err, Equals, ErrorInvalidRootKey)

	_, err = New(clt, "/")
	c.Assert(err, Equals, ErrorInvalidRootKey)

	_, err = New(clt, "abc/")
	c.Assert(err, Equals, ErrorInvalidRootKey)

	_, err = New(clt, "abc")
	c.Assert(err, Equals, nil)
}

func (s *ClientSuite) TestEnsureKey(c *C) {
	clt, err := clientv3.New(clientv3.Config{Endpoints: []string{TEST_ETCD_ADDR}})
	if err != nil {
		c.Error(err)
	}

	client, err := New(clt, "abc")
	c.Assert(err, Equals, nil)

	key, parentKey, err := client.ensureKey("abd")
	c.Assert(err, Equals, ErrorInvalidKey)

	key, parentKey, _ = client.ensureKey("/")
	c.Assert(key, Equals, "abc")
	c.Assert(parentKey, Equals, "abc")

	key, parentKey, _ = client.ensureKey("/def")
	c.Assert(key, Equals, "abc/def")
	c.Assert(parentKey, Equals, "abc")

	key, parentKey, _ = client.ensureKey("/def/ghi")
	c.Assert(key, Equals, "abc/def/ghi")
	c.Assert(parentKey, Equals, "abc/def")
}
