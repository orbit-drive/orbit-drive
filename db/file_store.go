package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/wlwanpan/orbit-drive/common"
)

type FileStore map[string]string

func GetFileStore() (FileStore, error) {
	store := make(FileStore)
	iter := Db.NewIterator(nil, nil)
	for iter.Next() {
		k := common.ToStr(iter.Key())
		switch k {
		case common.ROOT_KEY, common.CONFIG_KEY:
		default:
			v := common.ToStr(iter.Value())
			store[k] = v
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return FileStore{}, err
	}
	return store, nil
}

func (s FileStore) Dump() error {
	b := new(leveldb.Batch)
	for k, _ := range s {
		b.Delete(common.ToByte(k))
	}
	return Db.Write(b, nil)
}

func (s FileStore) Save() error {
	b := new(leveldb.Batch)
	for k, v := range s {
		b.Put(common.ToByte(k), common.ToByte(v))
	}
	return Db.Write(b, nil)
}
