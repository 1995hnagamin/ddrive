package node

import (
	"crypto/sha256"
	"fmt"
	"net"

	"../id"
)

type Node struct {
	Address net.IP
	Port    int
	Name    string
}

func (node *Node) String() string {
	return fmt.Sprintf("%s:%d/%s", node.Address, node.Port, node.Name)
}

func (node *Node) hash() [32]byte {
	return sha256.Sum256([]byte(node.String()))
}

func (node *Node) Id() *id.ID {
	hash := node.hash()
	return id.NewID(hash[:])
}
