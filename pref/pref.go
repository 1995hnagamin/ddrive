package pref

import (
	"net"
	"strconv"

	"./../node"
	"github.com/boltdb/bolt"
)

type Preference struct {
	DB *bolt.DB
}

const (
	bucketName = "bucket"
)

func NewPreference(filename string) (*Preference, error) {
	DB, err := bolt.Open("database", 0644, nil)
	if err != nil {
		return nil, err
	}
	return &Preference{DB}, nil
}

func (pref *Preference) Close() {
	defer pref.DB.Close()
}

func (pref *Preference) performGetTransaction(key string) (string, error) {
	var value []byte
	err := pref.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		value = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (pref *Preference) getFromDB(a, p, n string) (*node.Node, error) {
	address, err := pref.performGetTransaction(a)
	if err != nil {
		return nil, err
	}
	portStr, err := pref.performGetTransaction(p)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	name, err := pref.performGetTransaction(n)
	if err != nil {
		return nil, err
	}
	return &node.Node{
			net.ParseIP(address),
			port,
			name},
		nil
}

func (pref *Preference) GetSuccessor() (*node.Node, error) {
	return pref.getFromDB("SuccessorAddress", "SuccessorPort", "SuccessorName")
}

func (pref *Preference) GetSelf() (*node.Node, error) {
	return pref.getFromDB("SelfAddress", "SelfPort", "SelfName")
}
