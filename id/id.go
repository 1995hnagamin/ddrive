package id

import (
	"math/big"

	"./mod"
)

type ID struct {
	Id *mod.Mod
}

func NewID(hash []byte) *ID {
	id := big.NewInt(0)
	id.SetString(string(hash), 16)
	return &ID{mod.NewMod(id)}
}

func IsBetween(obj, x, y *ID) bool {
	xToObj := mod.Sub(obj.Id, x.Id)
	xToY := mod.Sub(y.Id, x.Id)
	return mod.Less(xToObj, xToY)
}
