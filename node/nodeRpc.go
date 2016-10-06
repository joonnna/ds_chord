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


func (n *Node) PutKey(args shared.Args, reply *shared.Reply) error {
	ok, _ := util.InKeySpace(n.id, args.Key, n.prev.Id)
	if ok {
		n.data[args.Key.String()] = args.Value
		return nil
	} else {
		n.logger.Error("PUT Wrong node " + args.Key.String())
		return ErrPut
	}
}

func (n *Node) GetKey(args shared.Args, reply *shared.Reply) error {
	ok, _ := util.InKeySpace(n.id, args.Key, n.prev.Id)
	if ok {
		reply.Value = n.data[args.Key.String()]
		return nil
	} else {
		n.logger.Error("GET Wrong node " + args.Key.String())
		return ErrGet
	}
}


func (n *Node) FindSuccessor(args shared.Test, reply *shared.Reply) error {
	tmp := new(big.Int)
	tmp.SetBytes(args.Id)
	test := shared.NodeInfo {
		Ip: args.Ip,
		Id: *tmp }

	if n.table.fingers[1].node.Ip == n.ip {
		reply.Next.Ip = n.ip
		reply.Next.Id = n.id
		reply.Prev = n.prev
		return nil
	}
	r, err := n.findPreDecessor(test.Id)
	if err != nil {
		return err
	}
	reply.Next = r.Next
	reply.Prev = r.Prev
	return nil
}


func (n *Node) ClosestPrecedingFinger(args shared.Test, reply *shared.Reply) error{
	cmpId := new(big.Int)
	cmpId.SetBytes(args.Id)
	for i := (lenOfId-1); i >= 1; i-- {
		entry := n.table.fingers[i].node.Id
		ok, _ := util.BetweenNodes(n.id, *cmpId, entry)
		if entry.BitLen() != 0 && ok {
			reply.Next = n.table.fingers[i].node
			return nil
		}
	}

	return ErrFingerNotFound
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
	n.logger.Info("RECEIVED NOTIFY")
	ok, _ := util.BetweenNodes(n.prev.Id, n.id, *tmp)
	if n.prev.Ip == n.ip || ok {
		node := shared.NodeInfo {
			Id: *tmp,
			Ip: args.Ip	}
		n.logger.Info("UPDATING PRE")
		n.prev = node
		n.logger.Info(n.prev.Ip)
	}
	return nil
}
/*
func (n *Node) UpdateSuccessor(args shared.Args, info *shared.Reply) error {
	n.update.Lock()

	n.assertPreDecessor(args.Node.Id)
	n.prev = args.Node
	n.logger.Info("New Pre-descessor " + n.prev.Ip)

	n.update.Unlock()
	return nil
}

func (n *Node) UpdatePreDecessor(args shared.Args, info *shared.Reply) error {
	n.update.Lock()

	n.assertSuccessor(args.Node.Id)
	n.Next = args.Node
	n.logger.Info("New successor " + n.Next.Ip)

	n.update.Unlock()
	return nil
}

func (n *Node) FindSuccessor(args shared.Args, reply *shared.Reply) error {
	//n.logger.Debug("FindSuccessor on id " + args.Node.Id)
	if n.Next.Id == "" && n.prev.Id == "" {
		reply.Next.Ip = n.ip
		reply.Next.Id = n.id
		reply.Prev.Ip = n.ip
		reply.Prev.Id = n.id
		//n.logger.Debug("Second node case, linking circle")
	} else if util.InKeySpace(n.id, args.Key, n.prev.Id) {
		n.logger.Debug("In my keyspace")
		reply.Next.Ip = n.ip
		reply.Next.Id = n.id
		reply.Prev = n.prev
	} else {
		n.logger.Debug("Not in my keyspace")

		r := &shared.Reply{}
		//args := util.CreateArgs(args.Node.Ip, args.Node.Id)
		err := shared.SingleCall("Node.FindSuccessor", (n.Next.Ip + n.RpcPort), args, r)
		if err != nil {
			n.logger.Error(err.Error())
			return err
		}

		reply.Next = r.Next
		reply.Prev = r.Prev
		reply.Value = r.Value
	}
	return nil
}

*/

/*
func (n *Node) SplitKeys(args shared.UpdateArgs, reply *shared.UpdateReply) error {
	reply.Values = n.store.SplitStorage(args.Id, args.PrevId)
	return nil
}
*/
