package pref

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"./../node"
)

type preference struct {
	self        *node.Node
	successor   *node.Node
	predecessor *node.Node
}

var (
	UnInitializedError = errors.New("preference uninitialized")
)

var (
	sharedPreference, sharedErr = newPreference()
)

func newPreference() (*preference, error) {
	var p preference
	text, err := ioutil.ReadFile("node.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(text, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetSelf() (*node.Node, error) {
	if sharedErr != nil {
		return nil, sharedErr
	}
	return sharedPreference.self, nil
}

func GetSuccessor() (*node.Node, error) {
	if sharedErr != nil {
		return nil, sharedErr
	}
	return sharedPreference.successor, nil
}

func GetPredecessor() (*node.Node, error) {
	if sharedErr != nil {
		return nil, sharedErr
	}
	return sharedPreference.predecessor, nil
}
