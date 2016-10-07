package node

import (
	//"github.com/joonnna/ds_chord/logger"
	"github.com/joonnna/ds_chord/node_communication"
//	"github.com/joonnna/ds_chord/storage"
	"github.com/joonnna/ds_chord/util"
	"errors"
	"math/big"
)

var (
	ErrGet = errors.New("Unable to get key")
	ErrPut = errors.New("Unable to put key")

)


func (n *Node) PutKey(args shared.Test, reply *shared.Reply) error {
	key := util.ConvertToBigInt(args.Id)
	n.logger.Debug("PUT key : " + key.String())

	n.update.Lock()
	defer n.update.Unlock()
	if util.InKeySpace(n.prev.Id, n.id, key){
		n.data[key.String()] = args.Value
		return nil
	} else {
		n.logger.Error("PUT Wrong node " + key.String())
		n.logger.Error("Succ : " + n.table.fingers[1].node.Id.String())
		n.logger.Error("Self : " + n.id.String())
		return ErrPut
	}
}

func (n *Node) GetKey(args shared.Test, reply *shared.Reply) error {
	key := util.ConvertToBigInt(args.Id)
	n.logger.Debug("Get key : " + key.String())

	n.update.RLock()
	defer n.update.RUnlock()
	if util.InKeySpace(n.prev.Id, n.id, key){
		reply.Value = n.data[key.String()]
		return nil
	} else {
		n.logger.Error("GET Wrong node " + key.String())
		return ErrGet
	}
}


func (n *Node) FindSuccessor(args shared.Test, reply *shared.Reply) error {
	tmp := new(big.Int)
	tmp.SetBytes(args.Id)
	test := shared.NodeInfo {
		Ip: args.Ip,
		Id: *tmp }

	if n.table.fingers[1].node.Ip == n.Ip {
		reply.Next.Ip = n.Ip
		reply.Next.Id = n.id
		reply.Prev = n.prev
		return nil
	}
	node, err := n.findPreDecessor(test.Id)
	if err != nil {
		n.logger.Error(err.Error())
		return err
	}

	succ, err := n.getSucc(node.Ip)
	if err != nil {
		return err
	}

	reply.Next = succ
	reply.Prev = node

	return nil
}


func (n *Node) ClosestPrecedingFinger(args shared.Test, reply *shared.Reply) error{
	cmpId := new(big.Int)
	cmpId.SetBytes(args.Id)
	for i := (lenOfId-1); i >= 1; i-- {
		entry := n.table.fingers[i].node.Id
		if entry.BitLen() != 0 && util.BetweenNodes(n.id, *cmpId, entry) {
			reply.Next = n.table.fingers[i].node
			return nil
		}
	}

	reply.Next.Id = n.id
	reply.Next.Ip = n.Ip

	return nil
}
func (n *Node) GetPreDecessor(args int, reply *shared.Test) error {
	id := n.prev.Id.Bytes()
	reply.Id = id
	reply.Ip = n.prev.Ip
	return nil
}

func (n *Node) GetSuccessor(args int, reply *shared.Test) error {
	id := n.table.fingers[1].node.Id.Bytes()
	reply.Id = id
	reply.Ip = n.table.fingers[1].node.Ip
	return nil
}

func (n *Node) Notify(args shared.Test, reply *shared.Reply) error {
	tmp := new(big.Int)
	tmp.SetBytes(args.Id)
	//n.logger.Info("RECEIVED NOTIFY")

	node := shared.NodeInfo {
		Id: *tmp,
		Ip: args.Ip	}

	if n.prev.Ip == n.Ip || util.BetweenNodes(n.prev.Id, n.id, *tmp) {
	//	n.logger.Info("UPDATING PRE")
		n.prev = node
		n.logger.Info(n.prev.Ip)
	}

	if n.table.fingers[1].node.Ip == n.Ip {
		n.table.fingers[1].node = node
	}

	return nil
}
