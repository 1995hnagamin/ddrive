package main

import (
	"crypto/sha256"
	"fmt"
	"net"
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
