package storage

import (
	"testing"
)

func TestInMemSimpleKVSet(t *testing.T) {
	imskv := NewInMemSimpleKV()
	testSimpleKVSet(imskv, t)

}

func TestInMemSimpleKVDelete(t *testing.T) {
	imskv := NewInMemSimpleKV()
	testSimpleKVDelete(imskv, t)
}

func TestInMemSimpleKVHSet(t *testing.T) {
	imskv := NewInMemSimpleKV()
	testSimpleKVHSet(imskv, t)
}

func TestInMemSimpleKVHDel(t *testing.T) {
	imskv := NewInMemSimpleKV()
	testSimpleKVHDel(imskv, t)
}

func TestInMemSimpleKVIncrBy(t *testing.T) {
	imskv := NewInMemSimpleKV()
	testSimpleKVHIncrBy(imskv, t)
}

func TestInMemSimpleKVSetEx(t *testing.T) {
	imskv := NewInMemSimpleKV()
	testSimpleKVSetEx(imskv, t)
}

func TestInMemSimpleKVSAdd(t *testing.T) {
	imskv := NewInMemSimpleKV()
	testSimpleKVSAdd(imskv, t)
}

func TestInMemSimpleKVSMembers(t *testing.T) {
	imskv := NewInMemSimpleKV()
	testSimpleKVSMembers(imskv, t)
}

func TestInMemSimpleKVSRem(t *testing.T) {
	imskv := NewInMemSimpleKV()
	testSimpleKVSRem(imskv, t)
}

func TestInMemSimpleKVSCard(t *testing.T) {
	imskv := NewInMemSimpleKV()
	testSimpleKVSCard(imskv, t)
}
