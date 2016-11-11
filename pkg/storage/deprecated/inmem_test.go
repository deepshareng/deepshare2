package storage

import "testing"

func TestInmemSet(t *testing.T) {
	testKVSet(NewInMem(), t)
}

func TestInmemDelete(t *testing.T) {
	testKVDelete(NewInMem(), t)
}
