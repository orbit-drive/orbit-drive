package db

import (
	"encoding/json"
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/wlwanpan/orbit-drive/common"
)

// Source represents the meta data of a file stored locally.
type Source struct {
	// ipfs hash
	Src string `json:src`

	// file size
	Size int64 `json:size`

	// file md5 checksum
	Checksum string `json:checksum`
}

type Sources map[string]*Source

func NewSource(path string) *Source {
	fi, err := os.Stat(path)
	if err != nil || fi.IsDir() {
		return &Source{}
	}
	checksum, err := common.Md5Checksum(path)
	if err != nil {
		// CHeck how to deal with error here also
		log.Println(err)
	}
	return &Source{
		Src:      "",
		Size:     fi.Size(),
		Checksum: checksum,
	}
}

func (s *Source) SetSrc(src string) {
	s.Src = src
}

func (s Source) GetSrc() string {
	return s.Src
}

func (s *Source) Copy() *Source {
	return &Source{
		Src:      s.GetSrc(),
		Size:     s.Size,
		Checksum: s.Checksum,
	}
}

func (s *Source) IsUploaded() bool {
	return s.GetSrc() != ""
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

// GetSources iterates through db, populate and return Sources.
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

// ExtractSource look for and return a copy of Source and
// deletes the key from the mapping.
func (s Sources) ExtractSource(k string) *Source {
	source, exist := s[k]
	if exist {
		defer delete(s, k)
		return source.Copy()
	}
	return &Source{}
}

// Dump batch deletes all the entries in the mapping.
func (s Sources) Dump() error {
	b := new(leveldb.Batch)
	for k, _ := range s {
		b.Delete(common.ToByte(k))
	}
	return Db.Write(b, nil)
}

// Save batch put all the entries in the mapping.
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
