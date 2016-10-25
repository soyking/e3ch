package client

import (
	"github.com/coreos/etcd/clientv3"
	. "gopkg.in/check.v1"
)

const (
	TEST_CREATE_KEY = "/create_dir"
)

type CreateSuite struct{}

func (s *CreateSuite) TearDownTest(c *C) {
	_, err := client.client.Delete(client.ctx, TEST_ROOT_KEY+TEST_CREATE_KEY, clientv3.WithPrefix())
	if err != nil {
		c.Error(err)
	}
}

func (s *CreateSuite) TestCreate1(c *C) {
	// not a dir
	_, err := client.client.Put(client.ctx, TEST_ROOT_KEY+TEST_CREATE_KEY, "")
	if err != nil {
		c.Error(err)
	}

	// parentKey is not a directory
	c.Assert(
		client.Create(TEST_CREATE_KEY+"/def", ""),
		Equals,
		ErrorCreateKey,
	)
}

func (s *CreateSuite) TestCreate2(c *C) {
	_, err := client.client.Put(client.ctx, TEST_ROOT_KEY+TEST_CREATE_KEY, "")
	if err != nil {
		c.Error(err)
	}

	// key has been set
	c.Assert(
		client.Create(TEST_CREATE_KEY, ""),
		Equals,
		ErrorCreateKey,
	)
}

func (s *CreateSuite) TestCreate3(c *C) {
	// success
	c.Assert(
		client.Create(TEST_CREATE_KEY, ""),
		Equals,
		nil,
	)
}
