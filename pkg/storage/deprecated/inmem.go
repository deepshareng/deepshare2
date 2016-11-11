package storage

import "sync"

// inmem is an in memory implementation of KV interface.
type inmem struct {
	mu        sync.RWMutex // protect kvs, writable
	writeable bool
	kvs       map[string][]byte // kvs stores key-value paris in memory
}

type inmemTx struct {
	inmem *inmem
}

func NewInMem() KV {
	return &inmem{kvs: make(map[string][]byte)}
}

func (im *inmem) Begin(writeable bool) (Tx, error) {
	tx := &inmemTx{inmem: im}

	if writeable {
		im.mu.Lock()
		im.writeable = true
		return tx, nil
	}
	im.mu.RLock()
	return tx, nil
}

func (im *inmem) Close() error {
	// future access will simply panic...
	im.kvs = nil
	return nil
}

func (it *inmemTx) Get(k []byte) []byte {
	if k == nil {
		panic("nil key isn't supported")
	}
	_, ok := it.inmem.kvs[string(k)]
	if !ok {
		return nil
	}
	return it.inmem.kvs[string(k)]
}

func (it *inmemTx) Set(k []byte, v []byte) {
	if k == nil {
		panic("nil key isn't supported")
	}
	if v == nil {
		panic("nil value isn't supported")
	}
	if !it.inmem.writeable {
		panic("not allowed")
	}
	it.inmem.kvs[string(k)] = v
}

func (it *inmemTx) Delete(k []byte) {
	if k == nil {
		panic("nil key isn't supported")
	}
	if !it.inmem.writeable {
		panic("not allowed")
	}
	delete(it.inmem.kvs, string(k))
}

func (it *inmemTx) Commit() error {
	if it.inmem.writeable {
		it.inmem.writeable = false
		it.inmem.mu.Unlock()
	} else {
		it.inmem.mu.RUnlock()
	}
	return nil
}

func (it *inmemTx) Rollback() error {
	if it.inmem.writeable {
		panic("do not support rolling back writable tx")
	}
	it.inmem.mu.RUnlock()
	return nil
}
