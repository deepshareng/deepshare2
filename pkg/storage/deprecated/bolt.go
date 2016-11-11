package storage

import "github.com/boltdb/bolt"

var (
	bucketKV = []byte("kv")
)

// bolt is an boltdb based implementation of KV interface.
type boltkv struct {
	db *bolt.DB
}

type boltTx struct {
	tx *bolt.Tx
}

func NewBolt(path string) (KV, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(bucketKV)
		if err != nil && err != bolt.ErrBucketExists {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &boltkv{db: db}, nil
}

func (bkv *boltkv) Begin(writeable bool) (Tx, error) {
	tx, err := bkv.db.Begin(writeable)
	if err != nil {
		return nil, err
	}
	return &boltTx{tx: tx}, nil
}

func (bkv *boltkv) Close() error {
	return bkv.db.Close()
}

func (btx *boltTx) Get(k []byte) []byte {
	if k == nil {
		panic("nil key isn't supported")
	}
	return btx.tx.Bucket(bucketKV).Get(k)
}

func (btx *boltTx) Set(k []byte, v []byte) {
	if k == nil {
		panic("nil key isn't supported")
	}
	if v == nil {
		panic("nil value isn't supported")
	}
	btx.tx.Bucket(bucketKV).Put(k, v)
}

func (btx *boltTx) Delete(k []byte) {
	if k == nil {
		panic("nil key isn't supported")
	}
	btx.tx.Bucket(bucketKV).Delete(k)
}

func (btx *boltTx) Commit() error {
	return btx.tx.Commit()
}

func (btx *boltTx) Rollback() error {
	return btx.tx.Rollback()
}
