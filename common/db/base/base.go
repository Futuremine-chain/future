package base

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/btcsuite/goleveldb/leveldb"
	"github.com/btcsuite/goleveldb/leveldb/opt"
	"github.com/btcsuite/goleveldb/leveldb/util"
)

type Base struct {
	Db *leveldb.DB
}

func Open(path string) (*Base, error) {
	var err error
	opts := &opt.Options{
		OpenFilesCacheCapacity: 16,
		Strict:                 opt.DefaultStrict,
		Compression:            opt.NoCompression,
		BlockCacheCapacity:     8 * opt.MiB,
		WriteBuffer:            4 * opt.MiB,
	}
	b := &Base{}
	if b.Db, err = leveldb.OpenFile(path, opts); err != nil {
		if b.Db, err = leveldb.RecoverFile(path, nil); err != nil {
			return nil, errors.New(fmt.Sprintf(`err while recoverfile %s : %s`, path, err.Error()))
		}

	}
	return b, nil
}

func (b *Base) Close() error {
	return b.Db.Close()
}

func (b *Base) Update(key []byte, value []byte) error {
	return b.Db.Put(key, value, nil)
}

func (b *Base) Delete(key []byte) error {
	return b.Db.Delete(key, nil)
}

func (b *Base) Get(key []byte) ([]byte, error) {
	return b.Db.Get(key, nil)
}

func (b *Base) Clear(bucket string) {
	rs := b.Foreach(bucket)
	for key, _ := range rs {
		b.Db.Delete([]byte(key), nil)
	}
}

func (b *Base) Foreach(bucket string) map[string][]byte {
	rs := make(map[string][]byte)
	iter := b.Db.NewIterator(util.BytesPrefix(bytes.Join([][]byte{[]byte(bucket), []byte("-")}, []byte{})), nil)
	defer iter.Release()

	// Iter will affect RLP decoding and reallocate memory to value
	for iter.Next() {
		value := make([]byte, len(iter.Value()))
		copy(value, iter.Value())
		rs[string(iter.Key())] = value
	}
	return rs
}

func Key(bucket string, key []byte) []byte {
	return bytes.Join([][]byte{
		[]byte(bucket + "-"), key}, []byte{})
}
