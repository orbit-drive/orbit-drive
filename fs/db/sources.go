package db

import (
	"encoding/json"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/wlwanpan/orbit-drive/common"
)

// Source represents the meta data of a file stored locally.
type Source struct {
	// ipfs hash
	Src string `json:src`

	// file size
	Size int `json:size`

	// file md5 checksum
	Checksum string `json:checksum`
}

type Sources map[string]*Source

func (s *Source) SetSrc(src string) {
	s.Src = src
}

func (s Source) GetSrc() string {
	return s.Src
}

func (s *Source) IsUploaded() bool {
	return s.GetSrc() == ""
}

func (s *Source) Save(k []byte) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return Put(k, data)
}

// IsSame check if the 2 sources are the same.
func (s *Source) IsSame(c *Source) bool {
	return s.Size == c.Size && s.Checksum == c.Checksum
}

func GetSources() (Sources, error) {
	store := make(Sources)
	iter := Db.NewIterator(nil, nil)
	for iter.Next() {
		k := common.ToStr(iter.Key())
		switch k {
		case common.ROOT_KEY, common.CONFIG_KEY:
		default:
			s := &Source{}
			err := json.Unmarshal(iter.Value(), s)
			if err != nil {
				log.Println(err)
				continue
			}
			store[k] = s
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
	for k, source := range s {
		data, err := json.Marshal(source)
		if err != nil {
			log.Println(err)
			continue
		}
		b.Put(common.ToByte(k), data)
	}
	return Db.Write(b, nil)
}
