package chord

import (
	"github.com/joonnna/ds_chord/node"
	"github.com/joonnna/ds_chord/node_communication"
	"github.com/joonnna/ds_chord/util"
	"github.com/joonnna/ds_chord/logger"
	"os"
	"errors"
)


type Chord struct {
	node *node.Node
	log *logger.Logger
}

var (
	ErrFind = errors.New("Couldn't find successor")
	ErrPut = errors.New("Couldn't put key")
	ErrGet = errors.New("Couldn't get key")
)

func Init(nameServer, httpPort, rpcPort string) *Chord {
	l := new(logger.Logger)
	l.Init((os.Stdout), "Chord", 0)
	c := &Chord {
		node: node.InitNode(nameServer, httpPort, rpcPort),
		log: l}

	return c
}

func (c *Chord) FindSuccessor(id string) (string, error) {
	key := util.ConvertKey(id)

	r := &shared.Reply{}
	args := &shared.Test{
		Id: key }

	node, err := util.GetNode(c.node.Ip, c.node.NameServer)
	if err != nil {
		c.log.Error(err.Error())
		return "", ErrFind
	}
	err = shared.SingleCall("Node.FindSuccessor", (node + c.node.RpcPort), args, r)
	if err != nil {
		c.log.Error(err.Error())
		return "", ErrFind
	}
	return r.Next.Ip, nil
}


func (c *Chord) PutKey(address, key, value string) error {
	id := util.ConvertKey(key)

	r := &shared.Reply{}
	args := shared.Test {
		Id : id,
		Value: value }
	err := shared.SingleCall("Node.PutKey", (address + c.node.RpcPort), args, r)
	if err != nil {
		c.log.Error(err.Error())
		return ErrPut
	}
	return nil
}


func (c *Chord) GetKey(address, key string) (string, error) {
	id := util.ConvertKey(key)

	r := &shared.Reply{}
	args := &shared.Test{
		Id: id}

	err := shared.SingleCall("Node.GetKey", (address + c.node.RpcPort), args, r)
	if err != nil {
		c.log.Error(err.Error())
		return "", ErrGet
	}
	return r.Value, nil
}

func (c *Chord) Run() {
	node.Run(c.node)
}

