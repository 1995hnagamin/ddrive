package mod

import (
	"math/big"
)

type Mod struct {
	Num *big.Int
}

func NewMod(n *big.Int) *Mod {
	mod := big.NewInt(0)
	mod.Exp(big.NewInt(2), big.NewInt(256), nil)
	n.Mod(n, mod)
	return &Mod{n}
}

func Add(x, y *Mod) *Mod {
	tmp := big.NewInt(0)
	tmp.Add(x.Num, y.Num)
	return NewMod(tmp)
}

func Sub(x, y *Mod) *Mod {
	tmp := big.NewInt(0)
	tmp.Sub(x.Num, y.Num)
	return NewMod(tmp)
}

func Mul(x, y *Mod) *Mod {
	tmp := big.NewInt(0)
	tmp.Mul(x.Num, y.Num)
	return NewMod(tmp)
}

func Less(x, y *Mod) bool {
	return x.Num.Cmp(y.Num) == -1
}

func Equal(x, y *Mod) bool {
	return x.Num.Cmp(y.Num) == 0
}
