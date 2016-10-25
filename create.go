package client

func (clt *EtcdHRCHYClient) CreateDir(key string) error {
	return clt.Create(key, clt.dirValue)
}

// set kv or directory
func (clt *EtcdHRCHYClient) Create(key string, value string) error {
	success, err := clt.put(key, value, COMPARE_EQUAL)
	if err != nil {
		return err
	}

	if !success {
		return ErrorCreateKey
	}

	return nil
}
