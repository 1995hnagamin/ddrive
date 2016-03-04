package pref

import (
	"errors"
	"io/ioutil"
	"json"
	"os"
)

type preference struct {
	self        Node
	successor   Node
	predecessor Node
}

const (
	UnInitializedError = errors.new("preference uninitialized")
)

var (
	sharedPreference *preference = nil
	sharedErr        error       = UnInitializedError
)

func newPreference() (*preference, os.Error) {
	var p preference
	json, err := ioutil.ReadFile("node.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(json, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func Init() error {
	sharedPreference, sharedErr := newPreference()
	if sharedErr != nil {
		return sharedErr
	}
	return nil
}

func getSelf() (*Node, error) {
	if sharedErr != nil {
		return nil, sharedErr
	}
	return sharedPreference.self, nil
}

func getSuccessor() (*Node, err) {
	if sharedErr != nil {
		return nil, sharedErr
	}
	return sharedPreference.successor, nil
}

func getPredecessor() (*Node, error) {
	if sharedErr != nil {
		return nil, sharedErr
	}
	return sharedPreference.predecessor, nil
}
