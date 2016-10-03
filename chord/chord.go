package chord

import (
	"github.com/joonnna/ds_chord/node"
	"github.com/joonnna/ds_chord/node_communication"
	"github.com/joonnna/ds_chord/util"
	"errors"
)


type Chord struct {
	node *node.Node
}

var (
	ErrFind = errors.New("Couldn't find successor")
	ErrPut = errors.New("Couldn't put key")
	ErrGet = errors.New("Couldn't get key")
)

func Init(nameServer, httpPort, rpcPort string) *Chord {
	c := &Chord {
		node: node.InitNode(nameServer, httpPort, rpcPort) }

	return c
}

func (c *Chord) FindSuccessor(id string) (string, error) {
	hashKey := util.HashKey(id)

	r := &shared.Reply{}
	args := &shared.Args{
		Key: hashKey }
	err := shared.SingleCall("Node.FindSuccessor", (c.node.Next.Ip + c.node.RpcPort), args, r)
	if err != nil {
		return "", ErrFind
	}
	return r.Next.Ip, nil
}


func (c *Chord) PutKey(address, key, value string) error {
	hashKey := util.HashKey(key)

	r := &shared.Reply{}
	args := util.RpcArgs(hashKey, value)
	err := shared.SingleCall("Node.PutKey", (address + c.node.RpcPort), args, r)
	if err != nil {
		return ErrPut
	}
	return nil
}


func (c *Chord) GetKey(address, key string) (string, error) {
	hashKey := util.HashKey(key)

	r := &shared.Reply{}
	args := &shared.Args{
		Key: hashKey}

	err := shared.SingleCall("Node.GetKey", (address + c.node.RpcPort), args, r)
	if err != nil {
		return "", ErrGet
	}
	return r.Value, nil
}

func (c *Chord) Run() {
	node.Run(c.node)
}

