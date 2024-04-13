package client

import (
	"go.etcd.io/bbolt"
	"strconv"
	"testing"
)

func TestTransactionOp_Get(t *testing.T) {
	db, err := bbolt.Open("test.db", 0600, nil)
	if err != nil || db == nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()
	op, err := newClientOp(db)
	if err != nil {
		t.Fatalf("failed to create client op: %v", err)
	}
	// debian magnet
	test_data := "magnet:?xt=urn:btih:FNTJQAETXQIYA35LKDFTZNAYGW4VUA3C&dn=debian-12.5.0-amd64-netinst.iso&xl=659554304&tr=http%3A%2F%2Fbttracker.debian.org%3A6969%2Fannounce"
	err = op.Put([]byte("debian"), []byte(test_data))
	if err != nil {
		t.Fatalf("failed to put test data: %v", err)
	}
	newtest_data := op.Get([]byte("debian"))
	if string(newtest_data) != test_data {
		t.Fatalf("failed to put test data: got %s, want %s", string(newtest_data), test_data)
	}

	// delete test data from database
	err = op.Delete([]byte("debian"))
	if err != nil {
		t.Fatalf("failed to delete test data: %v", err)
	}
	return
}

func TestTransactionOp_ForEach(t *testing.T) {
	db, err := bbolt.Open("test.db", 0600, nil)
	if err != nil || db == nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()
	op, err := newClientOp(db)
	t.Logf("created client op")
	if err != nil {
		t.Fatalf("failed to create client op: %v", err)
	}
	testData := []string{"debian", "ubuntu", "centos"}
	for i, datum := range testData {
		err = op.Put([]byte(strconv.Itoa(i)), []byte(datum))
		t.Logf("inserted test data: %v", datum)
		if err != nil {
			t.Fatalf("failed to put test data: %v", err)
		}
	}
	err = op.ForEach(func(k, v []byte) error {
		t.Logf("key: %s, value: %s", string(k), string(v))
		key, err := strconv.Atoi(string(k))
		if err != nil {
			t.Fatalf("failed to convert key: %v", err)
			return err
		}
		if testData[key] != string(v) {
			t.Fatalf("failed to put test data: got %s, want %s", string(k), string(v))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to for test data")
	}
	return
}
