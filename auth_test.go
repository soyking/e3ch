package client

import (
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	. "gopkg.in/check.v1"
)

const (
	TEST_AUTH_ROLE     = "test_role"
	TEST_AUTH_FROM_KEY = "/test_from"
	TEST_AUTH_TO_KEY   = "/test_to"
)

type AuthSuite struct{}

func (s *AuthSuite) SetUpSuite(c *C) {
	_, err := client.client.RoleAdd(client.ctx, TEST_AUTH_ROLE)
	if err != nil {
		c.Error(err)
	}
}

func (s *AuthSuite) TearDownSuite(c *C) {
	_, err := client.client.RoleDelete(client.ctx, TEST_AUTH_ROLE)
	if err != nil {
		c.Error(err)
	}
}

func (s *AuthSuite) TestRole(c *C) {
	c.Assert(
		client.RoleGrantPermission(TEST_AUTH_ROLE, TEST_AUTH_FROM_KEY, TEST_AUTH_TO_KEY, clientv3.PermissionType(clientv3.PermRead)),
		Equals,
		nil,
	)
	perms, err := client.GetRolePerms(TEST_AUTH_ROLE)
	c.Assert(err, Equals, nil)
	b, _ := json.MarshalIndent(perms, "", "    ")
	c.Log(string(b))

	c.Assert(
		client.RoleRevokePermission(TEST_AUTH_ROLE, TEST_AUTH_FROM_KEY, TEST_AUTH_TO_KEY),
		Equals,
		nil,
	)
	perms, err = client.GetRolePerms(TEST_AUTH_ROLE)
	c.Assert(err, Equals, nil)
	b, _ = json.MarshalIndent(perms, "", "    ")
	c.Log(string(b))
}
