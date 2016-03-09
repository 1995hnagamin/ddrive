package pref

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"./../node"
)

type Preference struct {
	Self        *node.Node
	Successor   *node.Node
	Predecessor *node.Node
}

func NewPreference(filename string) (*preference, error) {
	var p preference
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(text, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
