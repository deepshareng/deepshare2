package storage

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

var boltpath string

func init() {
	dir, err := ioutil.TempDir(os.TempDir(), "bolt_test")
	if err != nil {
		log.Fatal(err)
	}
	boltpath = path.Join(dir, "db")
}

func TestBoltSet(t *testing.T) {
	kv, err := NewBolt(boltpath)
	if err != nil {
		t.Fatal(err)
	}

	testKVSet(kv, t)

	os.RemoveAll(boltpath)
}

func TestBoltDelete(t *testing.T) {
	kv, err := NewBolt(boltpath)
	if err != nil {
		t.Fatal(err)
	}

	testKVDelete(kv, t)

	os.RemoveAll(boltpath)
}
