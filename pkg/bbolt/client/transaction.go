package client

import "go.etcd.io/bbolt"

type transactionOp struct {
	db *bbolt.DB
}

func (t transactionOp) ForEach(fn func(k []byte, v []byte) error) error {
	tx, err := t.db.Begin(true)
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket(bucketName)
	if b == nil {
		return nil
	}
	key, val := b.Cursor().First()
	print(key, val)
	err = b.ForEach(fn)
	if err != nil {
		return nil
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return nil
	}
	return err
}

func (t transactionOp) Get(key []byte) []byte {
	tx, err := t.db.Begin(true)
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket(bucketName)
	if b == nil {
		return nil
	}
	v := b.Get(key)
	if err != nil {
		return nil
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return nil
	}
	return v
}

func (t transactionOp) Put(key []byte, value []byte) error {
	tx, err := t.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket(bucketName)
	if b == nil {
		return err
	}

	err = b.Put(key, value)
	if err != nil {
		return err
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}
	return err
}

func (t transactionOp) Delete(key []byte) error {
	tx, err := t.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket(bucketName)
	if b == nil {
		return err
	}

	err = b.Delete(key)
	if err != nil {
		return err
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}
	return err
}

func newClientOp(db *bbolt.DB) (Client, error) {
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Use the transaction...
	_, err = tx.CreateBucketIfNotExists(bucketName)
	if err != nil {
		return nil, err
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &transactionOp{db: db}, nil
}
