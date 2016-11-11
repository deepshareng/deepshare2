package storage

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/MISingularity/deepshare2/pkg/log"
)

type inMemSimpleKV struct {
	mu       sync.Mutex
	kvs      map[string][]byte
	expireAt map[string]time.Time
}

func NewInMemSimpleKV() SimpleKV {
	return &inMemSimpleKV{
		kvs:      make(map[string][]byte),
		expireAt: make(map[string]time.Time),
	}
}

func (imskv *inMemSimpleKV) Get(k []byte) ([]byte, error) {
	if k == nil {
		panic("nil key isn't supported")
	}
	if v, ok := imskv.expireAt[string(k)]; ok && v.Before(time.Now()) {
		log.Debug("Inmem data expired, k:", string(k))
		imskv.mu.Lock()
		delete(imskv.kvs, string(k))
		delete(imskv.expireAt, string(k))
		imskv.mu.Unlock()
		return nil, nil
	}
	b := imskv.kvs[string(k)]
	log.Debugf("Inmem Get Simple KV: k = %s; v = %s", string(k), string(b))
	return b, nil
}

func (imskv *inMemSimpleKV) Delete(k []byte) error {
	if k == nil {
		panic("nil key isn't supported")
	}
	imskv.mu.Lock()
	delete(imskv.kvs, string(k))
	imskv.mu.Unlock()
	return nil
}
func (imskv *inMemSimpleKV) Set(k []byte, v []byte) error {
	if k == nil {
		panic("nil key isn't supported")
	}
	imskv.mu.Lock()
	imskv.kvs[string(k)] = v
	imskv.mu.Unlock()
	log.Debugf("Inmem Set Simple KV: k = %s; v = %s", string(k), string(v))
	return nil
}

func (imskv *inMemSimpleKV) SetEx(k []byte, v []byte, expiration time.Duration) error {
	if k == nil {
		panic("nil key isn't supported")
	}
	imskv.mu.Lock()
	imskv.kvs[string(k)] = v
	imskv.expireAt[string(k)] = time.Now().Add(expiration)
	imskv.mu.Unlock()
	log.Debugf("Inmem SetEx Simple KV: k = %s; v = %s, expiration = %v\n", string(k), string(v), expiration)
	return nil
}

func (imskv *inMemSimpleKV) HSet(k []byte, hk string, v []byte) error {
	if k == nil {
		panic("nil key isn't supported")
	}
	if hk == "" {
		panic("nil sub key isn't supported")
	}
	imskv.mu.Lock()
	defer imskv.mu.Unlock()
	if data, ok := imskv.kvs[string(k)]; !ok {
		m := map[string][]byte{hk: v}
		b, err := json.Marshal(&m)
		if err != nil {
			panic(err)
		}
		imskv.kvs[string(k)] = b
	} else {
		var m map[string][]byte
		if err := json.Unmarshal(data, &m); err != nil {
			panic(err)
		}
		m[hk] = v
		b, err := json.Marshal(&m)
		if err != nil {
			panic(err)
		}
		imskv.kvs[string(k)] = b
	}
	log.Debugf("Inmem HSet Simple KV: k = %s; hk = %s, v = %s", string(k), hk, string(v))
	return nil
}

func (imskv *inMemSimpleKV) HIncrBy(k []byte, hk string, n int) error {
	if k == nil {
		panic("nil key isn't supported")
	}
	if hk == "" {
		panic("nil sub key isn't supported")
	}
	imskv.mu.Lock()
	defer imskv.mu.Unlock()
	if data, ok := imskv.kvs[string(k)]; !ok {
		s := strconv.Itoa(n)
		m := map[string][]byte{hk: []byte(s)}
		b, err := json.Marshal(&m)
		if err != nil {
			panic(err)
		}
		imskv.kvs[string(k)] = b
	} else {
		var m map[string][]byte
		if err := json.Unmarshal(data, &m); err != nil {
			panic(err)
		}
		if m[hk] == nil {
			m[hk] = []byte(strconv.Itoa(n))
		} else {
			s := string(m[hk])
			if base, err := strconv.Atoi(s); err != nil {
				panic(err)
			} else {
				m[hk] = []byte(strconv.Itoa(base + n))
			}
		}
		b, err := json.Marshal(&m)
		if err != nil {
			panic(err)
		}
		imskv.kvs[string(k)] = b
	}
	log.Debugf("Inmem HIncrBy Simple KV: k = %s; n = %d", string(k), n)
	return nil
}

func (imskv *inMemSimpleKV) HGet(k []byte, hk string) ([]byte, error) {
	if k == nil {
		panic("nil key isn't supported")
	}
	if hk == "" {
		panic("nil sub key isn't supported")
	}
	b, ok := imskv.kvs[string(k)]
	if !ok {
		return nil, nil
	}
	var m map[string][]byte
	if err := json.Unmarshal(b, &m); err != nil {
		panic(err)
	}
	v, _ := m[hk]
	return v, nil
}

func (imskv *inMemSimpleKV) HGetAll(k []byte) (map[string]string, error) {
	if k == nil {
		panic("nil key isn't supported")
	}
	b, ok := imskv.kvs[string(k)]
	if !ok {
		return nil, nil
	}
	var m map[string][]byte
	if err := json.Unmarshal(b, &m); err != nil {
		panic(err)
	}
	mRet := make(map[string]string)
	for k, v := range m {
		mRet[k] = string(v)
	}

	return mRet, nil
}

func (imskv *inMemSimpleKV) HDel(k []byte, hk string) error {
	if k == nil {
		panic("nil key isn't supported")
	}
	if hk == "" {
		panic("nil sub key isn't supported")
	}
	imskv.mu.Lock()
	defer imskv.mu.Unlock()
	b, ok := imskv.kvs[string(k)]
	if !ok {
		return nil
	}
	var m map[string][]byte
	if err := json.Unmarshal(b, &m); err != nil {
		panic(err)
	}
	delete(m, hk)
	if len(m) == 0 {
		delete(imskv.kvs, string(k))
		return nil
	}

	b, err := json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	imskv.kvs[string(k)] = b
	return nil
}

func (imskv *inMemSimpleKV) Exists(k []byte) bool {
	_, ok := imskv.kvs[string(k)]
	return ok
}

func (imskv *inMemSimpleKV) SAdd(k []byte, v string) error {
	if k == nil {
		panic("nil key isn't supported")
	}
	if v == "" {
		panic("nil value isn't supported")
	}
	imskv.mu.Lock()
	defer imskv.mu.Unlock()
	b, ok := imskv.kvs[string(k)]
	var m map[string][]byte
	if ok {
		if err := json.Unmarshal(b, &m); err != nil {
			panic(err)
		}
	} else {
		m = make(map[string][]byte)
	}
	m[v] = []byte("1")

	b, err := json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	imskv.kvs[string(k)] = b
	log.Debugf("Inmem SAdd Simple KV: k = %s; v = %s", string(k), string(v))
	return nil
}

func (imskv *inMemSimpleKV) SRem(k []byte, v string) error {
	if k == nil {
		panic("nil key isn't supported")
	}
	if v == "" {
		panic("nil value isn't supported")
	}
	imskv.mu.Lock()
	defer imskv.mu.Unlock()
	b, ok := imskv.kvs[string(k)]
	if !ok {
		return nil
	}
	var m map[string][]byte
	if err := json.Unmarshal(b, &m); err != nil {
		panic(err)
	}
	delete(m, v)

	b, err := json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	imskv.kvs[string(k)] = b
	log.Debugf("Inmem SRem Simple KV: k = %s; v = %s", string(k), string(v))
	return nil
}

func (imskv *inMemSimpleKV) SCard(k []byte) (int64, error) {
	if k == nil {
		panic("nil key isn't supported")
	}
	b, ok := imskv.kvs[string(k)]
	if !ok {
		return 0, nil
	}
	var m map[string][]byte
	if err := json.Unmarshal(b, &m); err != nil {
		panic(err)
	}
	return int64(len(m)), nil
}

func (imskv *inMemSimpleKV) SMembers(k []byte) ([]string, error) {
	s := []string{}
	b := imskv.kvs[string(k)]
	m := make(map[string]string)
	if len(b) == 0 {
		return nil, nil
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	for k, _ := range m {
		s = append(s, k)
	}
	return s, nil
}
