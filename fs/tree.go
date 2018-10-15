package fs

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	ROOT_KEY string = "ROOT_TREE"

	FileCode = iota
	DirCode  = iota
)

var (
	Tree *Object
)

type Object struct {
	Id     []byte   `json:_id`
	Path   string   `json:path`
	Type   int      `json:'type'`
	Links  [][]byte `json:links`
	Source string   `json:source`
}

func (obj *Object) Save() error {
	if obj.Source == "" {
		s, err := UploadFile(obj.Path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		obj.Source = s
	}
	data, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return Db.Put([]byte(obj.Id), data, nil)
}

func (obj *Object) Load() error {
	if len(obj.Id) == 0 {
		return ErrInvalidKey
	}
	data, err := Db.Get(obj.Id, nil)
	if err != nil {
		return err
	}
	log.Println(string(data[:]))

	json.Unmarshal(data, obj)
	return nil
}

func InitTree() error {
	data, err := Db.Get([]byte(ROOT_KEY), nil)
	if err != nil {
		if err == ErrNotFound {
			Tree = &Object{}
			return nil
		}
		return err
	}
	err = json.Unmarshal(data, Tree)
	if err != nil {
		return err
	}
	return nil
}

func GenTreeFromPath(c *Config) (*Object, error) {
	files, err := ioutil.ReadDir(c.Root)
	if err != nil {
		return &Object{}, err
	}

	t := &Object{
		Id:    []byte(ROOT_KEY),
		Path:  "/",
		Type:  DirCode,
		Links: [][]byte{},
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		abspath := c.Root + "/" + f.Name()
		k := GenKeyFromPath(abspath)

		has, err := Db.Has(k, nil)
		if err != nil {
			return &Object{}, err
		}
		if !has {
			n := &Object{
				Id:     k,
				Path:   abspath,
				Type:   FileCode,
				Links:  [][]byte{},
				Source: "",
			}
			go n.Save()
		}

		t.Links = append(t.Links, k)
	}

	return t, nil
}

func GenKeyFromPath(p string) []byte {
	hash := sha256.Sum256([]byte(p))
	return hash[:]
}
