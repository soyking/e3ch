package client

func (clt *EtcdHRCHYClient) CreateDir(key string) error {
	return clt.Create(key, clt.dirValue)
}

// set kv or directory
func (clt *EtcdHRCHYClient) Create(key string, value string) error {
	return clt.put(key, value, true)
}
