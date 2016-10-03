package node

import (
	//"github.com/joonnna/ds_chord/logger"
	"github.com/joonnna/ds_chord/node_communication"
//	"github.com/joonnna/ds_chord/storage"
	"github.com/joonnna/ds_chord/util"
	"errors"
)

var (
	ErrGet = errors.New("Unable to get key")
	ErrPut = errors.New("Unable to put key")

)


func (n *Node) PutKey(args shared.Args, reply *shared.Reply) error {
	if util.InKeySpace(n.id, args.Key, n.prev.Id) {
		n.data[args.Key] = args.Value
		return nil
	} else {
		n.logger.Error("PUT Wrong node " + args.Key)
		return ErrPut
	}
}

func (n *Node) GetKey(args shared.Args, reply *shared.Reply) error {
	if util.InKeySpace(n.id, args.Key, n.prev.Id) {
		reply.Value = n.data[args.Key]
		return nil
	} else {
		n.logger.Error("GET Wrong node " + args.Key)
		return ErrGet
	}
}

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
/*
func (n *Node) SplitKeys(args shared.UpdateArgs, reply *shared.UpdateReply) error {
	reply.Values = n.store.SplitStorage(args.Id, args.PrevId)
	return nil
}
*/
