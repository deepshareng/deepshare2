package storage

import (
	"log"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func testSimpleKVSet(skv SimpleKV, t *testing.T) {
	if err := skv.Set([]byte("foo"), []byte("bar")); err != nil {
		t.Fatal(err)
	}

	v, err := skv.Get([]byte("foo"))
	if err != nil {
		t.Fatal(err)
	}
	if string(v) != "bar" {
		t.Errorf("value = %s, want %s", string(v), "bar")
	}
}

func testSimpleKVDelete(skv SimpleKV, t *testing.T) {
	if err := skv.Set([]byte("foo"), []byte("bar")); err != nil {
		t.Fatal(err)
	}

	if err := skv.Delete([]byte("foo")); err != nil {
		t.Fatal(err)
	}

	v, err := skv.Get([]byte("foo"))
	if err != nil {
		t.Fatal("Get against a non-exist key, should return nil with nil err")
	}
	if v != nil {
		t.Errorf("value = %s, want nil", string(v))
	}

	if skv.Exists([]byte("foo")) {
		t.Error("Exists() should return false after delete")
	}
}

func testSimpleKVSetEx(skv SimpleKV, t *testing.T) {
	expiration := time.Millisecond * time.Duration(50)
	sleepDuration := time.Millisecond * time.Duration(60)
	if err := skv.SetEx([]byte("foo"), []byte("bar"), expiration); err != nil {
		t.Fatal(err)
	}

	v, err := skv.Get([]byte("foo"))
	if err != nil {
		t.Fatal(err)
	}
	if string(v) != "bar" {
		t.Errorf("value = %s, want %s", string(v), "bar")
	}

	time.Sleep(sleepDuration)
	v, err = skv.Get([]byte("foo"))
	if err != nil {
		t.Fatal(err)
	}
	if v != nil {
		t.Errorf("should expire, value = %s, want nil", string(v))
	}
}

func testSimpleKVHSet(skv SimpleKV, t *testing.T) {
	if err := skv.HSet([]byte("foo"), "foo1", []byte("bar1")); err != nil {
		t.Fatal(err)
	}
	v, err := skv.HGet([]byte("foo"), "foo1")
	if err != nil {
		t.Fatal(err)
	}
	if string(v) != "bar1" {
		t.Errorf("value = %s, want = %s\n", string(v), "bar1")
	}

	if err := skv.HSet([]byte("foo"), "foo2", []byte("bar2")); err != nil {
		t.Fatal(err)
	}
	v, err = skv.HGet([]byte("foo"), "foo2")
	if err != nil {
		t.Fatal(err)
	}
	if string(v) != "bar2" {
		t.Errorf("value = %s, want = %s\n", string(v), "bar2")
	}

	m, err := skv.HGetAll([]byte("foo"))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"foo1": "bar1", "foo2": "bar2"}
	if !reflect.DeepEqual(m, want) {
		t.Errorf("HGetAll value = %v, want = %v", m, want)
	}
}

func testSimpleKVHIncrBy(skv SimpleKV, t *testing.T) {
	if err := skv.HIncrBy([]byte("foo"), "foo1", 5); err != nil {
		t.Fatal(err)
	}
	v, err := skv.HGet([]byte("foo"), "foo1")
	if err != nil {
		t.Fatal(err)
	}
	n, err := strconv.ParseInt(string(v), 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Errorf("value = %d, want = %d\n", n, 5)
	}

	if err := skv.HIncrBy([]byte("foo"), "foo1", 2); err != nil {
		t.Fatal(err)
	}
	v, err = skv.HGet([]byte("foo"), "foo1")
	if err != nil {
		t.Fatal(err)
	}
	n, err = strconv.ParseInt(string(v), 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	if n != 7 {
		t.Errorf("value = %d, want = %d\n", n, 7)
	}

	skv.Delete([]byte("foo"))
}

func testSimpleKVHDel(skv SimpleKV, t *testing.T) {
	if err := skv.HSet([]byte("foo"), "foo1", []byte("bar1")); err != nil {
		t.Fatal(err)
	}
	if err := skv.HDel([]byte("foo"), "foo1"); err != nil {
		t.Fatal(err)
	}
	v, err := skv.HGet([]byte("foo"), "foo1")
	if err != nil {
		t.Fatal(err)
	}
	if v != nil {
		t.Errorf("value = %s, want = nil\n", string(v))
	}
}

func testSimpleKVSAdd(skv SimpleKV, t *testing.T) {
	if err := skv.SAdd([]byte("foo"), "foo1"); err != nil {
		t.Fatal(err)
	}
	if err := skv.SAdd([]byte("32"), "555"); err != nil {
		t.Fatal(err)
	}
}

func testSimpleKVSMembers(skv SimpleKV, t *testing.T) {
	skv.SAdd([]byte("kk"), "a")
	skv.SAdd([]byte("kk"), "b")
	skv.SAdd([]byte("kk"), "c")
	s, err := skv.SMembers([]byte("kk"))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(s, []string{"a", "b", "c"}) &&
		!reflect.DeepEqual(s, []string{"a", "c", "b"}) &&
		!reflect.DeepEqual(s, []string{"b", "a", "c"}) &&
		!reflect.DeepEqual(s, []string{"b", "c", "a"}) &&
		!reflect.DeepEqual(s, []string{"c", "a", "b"}) &&
		!reflect.DeepEqual(s, []string{"c", "b", "a"}) {
		t.Error("testSimpleKVSMemebrs failed, want =", []string{"a", "b", "c"}, ", got =", s)
	}
}

func testSimpleKVSRem(skv SimpleKV, t *testing.T) {
	if err := skv.SAdd([]byte("foo"), "foo1"); err != nil {
		t.Fatal(err)
	}
	if err := skv.SRem([]byte("foo"), "foo1"); err != nil {
		t.Fatal(err)
	}
	if err := skv.SRem([]byte("foo"), "foo2"); err != nil {
		t.Fatal(err)
	}
}

func testSimpleKVSCard(skv SimpleKV, t *testing.T) {
	testcases := []struct {
		pairs []struct{ op, k, v string }
		k     string
		count int64
	}{
		{
			[]struct{ op, k, v string }{
				{"SAdd", "k1", "v1"},
				{"SAdd", "k1", "v2"},
				{"SAdd", "k1", "v1"},
				{"SAdd", "k1", "v3"},
				{"SAdd", "k1", "v3"},
			},
			"k1",
			3,
		},
		{
			[]struct{ op, k, v string }{
				{"SAdd", "k2", "v1"},
				{"SAdd", "k2", "v2"},
				{"SAdd", "k2", "v3"},
				{"SRem", "k2", "v1"},
				{"SRem", "k2", "v1"},
			},
			"k2",
			2,
		},
	}

	for i, v := range testcases {
		err := skv.Delete([]byte(v.k))
		if err != nil {
			log.Fatal(err)
		}
		for _, item := range v.pairs {
			switch item.op {
			case "SAdd":
				err := skv.SAdd([]byte(item.k), item.v)
				if err != nil {
					t.Fatal(err)
				}
			case "SRem":
				err := skv.SRem([]byte(item.k), item.v)
				if err != nil {
					t.Fatal(err)
				}
			}
		}
		count, err := skv.SCard([]byte(v.k))
		if err != nil {
			t.Fatal(err)
		}
		if count != v.count {
			t.Errorf("#%d testcase: want size of set = %d, get = %d\n", i, v.count, count)
		}
	}
}
