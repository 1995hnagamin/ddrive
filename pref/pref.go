package pref

import (
	"net"
	"strconv"

	"./../node"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
)

type Preference struct {
	DB store.Store
}

func NewPreference(filename string) (*Preference, error) {
	DB, err := boltdb.New([]string{"database"}, nil)
	if err != nil {
		return nil, err
	}
	return &Preference{DB}, nil
}

func (pref *Preference) Close() {
	defer pref.DB.Close()
}

func (pref *Preference) getFromDB(a, p, n string) (*node.Node, error) {
	addresskv, err := pref.DB.Get(a)
	if err != nil {
		return nil, err
	}
	address := net.ParseIP(string(addresskv.Value))
	portkv, err := pref.DB.Get(p)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(string(portkv.Value))
	if err != nil {
		return nil, err
	}
	namekv, err := pref.DB.Get(n)
	if err != nil {
		return nil, err
	}
	return &node.Node{address, port, string(namekv.Value)}, nil
}

func (pref *Preference) GetSuccessor() (*node.Node, error) {
	return pref.getFromDB("SuccessorAddress", "SuccessorPort", "SuccessorName")
}

func (pref *Preference) GetSelf() (*node.Node, error) {
	return pref.getFromDB("SelfAddress", "SelfPort", "SelfName")
}
