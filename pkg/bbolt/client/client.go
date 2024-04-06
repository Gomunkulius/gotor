package client

import "go.etcd.io/bbolt"

type Client interface {
	Get(key []byte) []byte
	Put(key []byte, value []byte) error
	ForEach(fn func(k, v []byte) error) error
	Delete(key []byte) error
}

var bucketName = []byte(".gotor")

func NewClient(path string) (Client, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	return newClientOp(db)
}
