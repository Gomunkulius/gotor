package bbolt

import "go.etcd.io/bbolt"

type Client interface {
	Get(key []byte) []byte
	Put(key []byte, value []byte) error
	Delete(key []byte) error
}

var bucketName = []byte(".gotor")

func NewClient(path string) Client {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil
	}
	return newClientOp(db)
}
