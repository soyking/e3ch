package client

import (
	"go.etcd.io/etcd/clientv3"
	. "gopkg.in/check.v1"
)

const (
	TEST_PUT_KEY = "/put_dir"
)

type PutSuite struct{}

func (s *PutSuite) TearDownTest(c *C) {
	_, err := client.client.Delete(client.ctx, TEST_ROOT_KEY+TEST_PUT_KEY, clientv3.WithPrefix())
	if err != nil {
		c.Error(err)
	}
}

func (s *PutSuite) TestPut1(c *C) {
	_, err := client.client.Put(client.ctx, TEST_ROOT_KEY+TEST_PUT_KEY, client.dirValue)
	if err != nil {
		c.Error(err)
	}

	// key is directory
	c.Assert(
		client.Put(TEST_PUT_KEY, ""),
		Equals,
		ErrorPutKey,
	)
}

func (s *PutSuite) TestPut2(c *C) {
	_, err := client.client.Put(client.ctx, TEST_ROOT_KEY+TEST_PUT_KEY, "")
	if err != nil {
		c.Error(err)
	}

	// parentKey is not a directory
	c.Assert(
		client.Put(TEST_PUT_KEY+"/abc", ""),
		Equals,
		ErrorPutKey,
	)
}

func (s *PutSuite) TestPut3(c *C) {
	// has been set
	_, err := client.client.Put(client.ctx, TEST_ROOT_KEY+TEST_PUT_KEY, "")
	if err != nil {
		c.Error(err)
	}

	// success
	c.Assert(
		client.Put(TEST_PUT_KEY, ""),
		Equals,
		nil,
	)
}
