package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/wlwanpan/orbit-drive/common"
)

type Sources map[string]string

func GetSources() (Sources, error) {
	store := make(Sources)
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
		return Sources{}, err
	}
	return store, nil
}

func (s Sources) Dump() error {
	b := new(leveldb.Batch)
	for k, _ := range s {
		b.Delete(common.ToByte(k))
	}
	return Db.Write(b, nil)
}

func (s Sources) Save() error {
	b := new(leveldb.Batch)
	for k, v := range s {
		b.Put(common.ToByte(k), common.ToByte(v))
	}
	return Db.Write(b, nil)
}
