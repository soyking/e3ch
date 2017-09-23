e3ch
===

etcd v3 client with hierarchy

There are directory and key-value in etcd v2, which is convenient to manage the key-value store. But etcd v3 only supports flat key-value space (see [#633](https://github.com/coreos/etcd/issues/633#issuecomment-152768632)).Though you could use `prefix` to adjust the new API, it is not easy to manage key-value store or make the structure clearly. e3ch is built for making etcd v3 'look like' a key-value store supporting hierarchy.

## Design

e3ch will set a special value for a directory, which can be defined by user. For a key `/database/mysql/setting.json`, the directory key such as `/database` and `/database/mysql` will be set as the special value.

When creating/putting key-value, e3ch will confirm the parent key should be the directory special value. For example, when putting `/database/fake-mysql/setting.json`, if the value of `/database/fake-mysql` is not the special value, the action will be rejected.

When listing/deleting a directory, e3ch will get key-value with the prefix. For example, when listing `/database`, e3ch will get all key-value with prefix `/database/`. But it will only return the keys under the directory, and the keys in deeper directory will be ignored. So it may cost too long if there are too many keys under the directory, escpecially when listing the root `/`

## Example

```
package main

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/soyking/e3ch"
)

func main() {
	// initial etcd v3 client
	e3Clt, err := clientv3.New(clientv3.Config{Endpoints: []string{"127.0.0.1:2379"}})
	if err != nil {
		panic(err)
	}

	// new e3ch client with namespace(rootKey)
	clt, err := client.New(e3Clt, "my_space")
	if err != nil {
		panic(err)
	}

	// set the rootKey as directory
	err = clt.FormatRootKey()
	if err != nil {
		panic(err)
	}

	clt.CreateDir("/dir1")
	clt.Create("/dir1/key1", "")
	clt.Create("/dir", "")
	clt.Put("/dir1/key1", "value1")
	clt.Get("/dir1/key1")
	clt.List("/dir1")
	clt.Delete("/dir")
}
```

more examples in `*_test.go`

## Run Test

`go test -check.vv`
