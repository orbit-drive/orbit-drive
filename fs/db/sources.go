package db

import (
	"encoding/json"
	"log"
	"os"

	"github.com/orbit-drive/orbit-drive/fsutil"
	"github.com/syndtr/goleveldb/leveldb"
)

// Source represents the meta data of a file stored locally.
type Source struct {
	// Src represents the ipfs hash of the file.
	Src string `json:"src"`

	// Size represents the size of the file.
	Size int64 `json:"size"`

	// Checksum represents the md5 checksum hash of file.
	Checksum string `json:"checksum"`
}

// Sources represents the store of the locally saved files.
type Sources map[string]*Source

// NewSource generates a new source instance froma given path
// and validates the path, computes the file checksum and size.
func NewSource(path string) *Source {
	fi, err := os.Stat(path)
	if err != nil || fi.IsDir() {
		return &Source{}
	}
	checksum, err := fsutil.Md5Checksum(path)
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

// SetSrc is a setter for Source src.
func (s *Source) SetSrc(src string) {
	s.Src = src
}

// GetSrc is a getter for Source src.
func (s Source) GetSrc() string {
	return s.Src
}

// DeepCopy deep copies a source instance and return a new instance.
func (s *Source) DeepCopy() *Source {
	return &Source{
		Src:      s.GetSrc(),
		Size:     s.Size,
		Checksum: s.Checksum,
	}
}

// IsUploaded check is the Source src is a non zero value.
func (s *Source) IsNew() bool {
	return s.GetSrc() != ""
}

// Save write the Source instance to the db.
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
		k := fsutil.ToStr(iter.Key())
		switch k {
		case fsutil.ROOTKEY, fsutil.CONFIGKEY:
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
		return source.DeepCopy()
	}
	return &Source{}
}

// Dump batch deletes all the entries in the mapping.
func (s Sources) Dump() error {
	b := new(leveldb.Batch)
	for k := range s {
		b.Delete(fsutil.ToByte(k))
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
		b.Put(fsutil.ToByte(k), data)
	}
	return Db.Write(b, nil)
}
