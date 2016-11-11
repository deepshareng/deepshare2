package storage

import "testing"

func testKVSet(kv KV, t *testing.T) {
	tx, err := kv.Begin(true)
	if err != nil {
		t.Fatal(err)
	}

	tx.Set([]byte("foo"), []byte("bar"))
	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	tx, err = kv.Begin(false)
	v := tx.Get([]byte("foo"))
	if string(v) != "bar" {
		t.Errorf("value = %s, want %s", string(v), "bar")
	}

	err = tx.Rollback()
	if err != nil {
		t.Fatal(err)
	}

	kv.Close()
}

func testKVDelete(kv KV, t *testing.T) {
	tx, err := kv.Begin(true)
	if err != nil {
		t.Fatal(err)
	}

	tx.Set([]byte("foo"), []byte("bar"))
	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	// delete foo=bar
	tx, err = kv.Begin(true)
	if err != nil {
		t.Fatal(err)
	}

	tx.Delete([]byte("foo"))
	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	tx, err = kv.Begin(false)
	v := tx.Get([]byte("foo"))
	if v != nil {
		t.Errorf("value = %s, want nil", string(v))
	}

	err = tx.Rollback()
	if err != nil {
		t.Fatal(err)
	}

	kv.Close()
}
