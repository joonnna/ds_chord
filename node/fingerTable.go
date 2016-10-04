package node

import (
	"github.com/joonnna/ds_chord/node_communication"
	"github.com/joonnna/ds_chord/util"
	"math/rand"
	"math/big"
	"time"
	"errors"
)

const (
	lenOfId = 160
)


type fingerEntry struct {
	node shared.NodeInfo
	start *big.Int
}

type fingerTable struct {
	fingers []fingerEntry
}


var (
	ErrFingerNotFound = errors.New("Cant find closesetpreceding finger")
)


func calcStart(exponent int, modExp int, id *big.Int) *big.Int {
	base2 := big.NewInt(int64(2))
	k := big.NewInt(int64(exponent))

	tmp := big.NewInt(int64(0))
	tmp.Exp(base2, k, nil)

	sum := big.NewInt(int64(0))
	sum.Add(tmp, id)

	modExponent := big.NewInt(int64(modExp))
	mod := big.NewInt(int64(0))
	mod.Exp(base2, modExponent, nil)

	ret := big.NewInt(int64(0))

	ret.Mod(sum, mod)

	return ret
}


func (n *Node) initFingerTable() {
	n.table.fingers = make([]fingerEntry, lenOfId)

	for i := 1; i < (lenOfId-1); i++ {
		n.table.fingers[i].start = calcStart(i, lenOfId, n.id)
		n.logger.Debug(n.table.fingers[i].start.String())
	}
	n.logger.Debug("Inited finger table")
}



func (n *Node) FindSuccessor(args shared.Args, reply *shared.Reply) error {
	r, err := n.findPreDecessor(args.Node)
	if err != nil {
		return err
	}
	reply.Next = r.Next
	return nil
}



func (n *Node) fixFingers() {
	for {
		index := rand.Int() % lenOfId

		r := &shared.Reply{}
		args := shared.NodeInfo {
			Ip: n.table.fingers[1].node.Ip,
			Id: n.table.fingers[index].start }
		err := shared.SingleCall("Node.Findsuccessor", (n.table.fingers[1].node.Ip + n.RpcPort), args, r)
		if err != nil {
			n.logger.Error(err.Error())
		}

		n.table.fingers[index].node = r.Next
		time.Sleep(time.Second*5)
	}
}

func(n *Node) stabilize() {
	for {
		currPre, err := n.findPreDecessor(n.Next)
		if err != nil {
			n.logger.Error(err.Error())
			time.Sleep(time.Second*5)
			continue
		}
		if util.InKeySpace(n.table.fingers[1].node.Id, currPre.Prev.Id, n.id) {
			n.table.fingers[1].node = currPre.Prev
		}

		r := &shared.Reply{}
		args := shared.NodeInfo {
			Ip: n.table.fingers[1].node.Ip,
			Id: n.table.fingers[1].node.Id }
		err = shared.SingleCall("Node.Notify", (n.table.fingers[1].node.Ip + n.RpcPort), args, r)
		if err != nil {
			n.logger.Error(err.Error())
		}
		time.Sleep(time.Second*5)
	}
}

func (n *Node) join() {
	list, err := GetNodeList(n.NameServer)
	if err != nil {
		n.logger.Error(err.Error())
		return
	}
	randomNode := util.GetNode(list, n.ip)
	if randomNode == "" {
		self := shared.NodeInfo {
			Id: n.id,
			Ip: n.ip }
		n.table.fingers[1].node = self
		return
	}

	node := shared.NodeInfo{
		Ip: n.ip,
		Id: n.id }
	args := shared.Args{
		Node: node }

	r := &shared.Reply{}

	err = shared.SingleCall("Node.FindSuccessor", (randomNode +n.RpcPort), args, r)
	if err != nil {
		n.logger.Error(err.Error())
		return
	}

	n.table.fingers[1].node = r.Next
}


func (n *Node) Notify(args shared.Args, reply *shared.Reply) error {
	if n.prev.Id == nil || util.InKeySpace(n.id, args.Node.Id, n.prev.Id) {
		n.prev = args.Node
	}
	return nil
}

func (n *Node) findPreDecessor(node shared.NodeInfo) (*shared.Reply, error) {
	currNode := node
	r := &shared.Reply{}
	for {
		if util.InKeySpace(n.Next.Id, currNode.Id, n.id) {
			return r, nil
		}
		args := shared.Args{
			Key: node.Id }

		err := shared.SingleCall("Node.ClosestPrecedingFinger", (node.Ip + n.RpcPort), args, r)
		if err != nil {
			n.logger.Error(err.Error())
			return nil, err
		}
		currNode = r.Next
	}
}

func (n *Node) ClosestPrecedingFinger(args shared.Args, reply *shared.Reply) error{
	for i := lenOfId; i >= 1; i-- {
		entry := n.table.fingers[i].node.Id
		if util.InKeySpace(args.Key, entry, n.id){
			reply.Next = n.table.fingers[i].node
			reply.Prev = n.prev
			return nil
		}
	}

	return ErrFingerNotFound
}
