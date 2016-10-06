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
	start big.Int
}

type fingerTable struct {
	fingers []fingerEntry
}


var (
	ErrFingerNotFound = errors.New("Cant find closesetpreceding finger")
)


func calcStart(exponent int, modExp int, id big.Int) big.Int {
	base2 := big.NewInt(int64(2))
	k := big.NewInt(int64(exponent))

	tmp := big.NewInt(int64(0))
	tmp.Exp(base2, k, nil)

	sum := big.NewInt(int64(0))
	sum.Add(tmp, &id)

	modExponent := big.NewInt(int64(modExp))
	mod := big.NewInt(int64(0))
	mod.Exp(base2, modExponent, nil)

	ret := big.NewInt(int64(0))

	ret.Mod(sum, mod)

	return *ret
}


func (n *Node) initFingerTable() {
	n.table.fingers = make([]fingerEntry, lenOfId)

	for i := 1; i < (lenOfId-1); i++ {
		n.table.fingers[i].start = calcStart(i, lenOfId, n.id)
	}
	n.logger.Debug("Inited finger table")
}

func (n *Node) fixFingers() {
	for {
		index := rand.Int() % lenOfId
		if index == 1 {
			index = 2
		}
		r := &shared.Reply{}
		args := shared.Test {
			Ip: n.table.fingers[1].node.Ip,
			Id: n.table.fingers[index].start.Bytes() }
		err := shared.SingleCall("Node.FindSuccessor", (n.table.fingers[1].node.Ip + n.RpcPort), args, r)
		if err != nil {
			n.logger.Error("FAILED TO SET SUCCESOR")
			n.logger.Error(err.Error())
		}

		n.logger.Info("SET SUCCESSOR")
		n.table.fingers[index].node = r.Next
		time.Sleep(time.Second*5)
	}
}

func(n *Node) stabilize() {
	var succ string
	for {

		n.logger.Info("STABILIZING")
		arg := 0

		r := &shared.Test{}
		succ = n.table.fingers[1].node.Ip
	/*
		if succ == n.ip {
			r.Id = n.id.Bytes()
			r.Ip = n.ip
		} else {*/
			err := shared.SingleCall("Node.GetPreDecessor", (succ + n.RpcPort), arg, r)
			if err != nil {
				n.logger.Error(err.Error())
				time.Sleep(time.Second*5)
				continue
			}
		//}
		n.logger.Debug("Found predecessor and want to stabilize, pre : " + r.Ip)
		tmp := new(big.Int)
		tmp.SetBytes(r.Id)

		ok, _ := util.BetweenNodes(n.id, n.table.fingers[1].node.Id, *tmp)
		if succ == n.ip || ok {
			n.logger.Debug("New succ is :" + r.Ip)
			n.table.fingers[1].node.Ip = r.Ip
			n.table.fingers[1].node.Id = *tmp
			if n.table.fingers[1].node.Ip == n.ip {
				n.logger.Debug("Found myself as succ, sleeping")
				time.Sleep(time.Second*5)
				continue
			}
		}

		args := shared.Test {
			Ip: n.ip,
			Id: n.id.Bytes()}

		reply := &shared.Reply{}
		n.logger.Info("SENDING NOTIFY")
		err = shared.SingleCall("Node.Notify", (n.table.fingers[1].node.Ip + n.RpcPort), args, reply)
		if err != nil {
			n.logger.Error(err.Error())
		}
		time.Sleep(time.Second*5)
	}
}

func (n *Node) join() {
	n.putIp()

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
		n.prev = self
		n.logger.Debug("ONLY MYSELF CASE")
		return
	}
	n.logger.Debug("NOT MYSELF CASE")
	node := shared.Test {
		Ip: n.ip,
		Id: n.id.Bytes() }

	r := &shared.Reply{}
	n.logger.Debug(randomNode)
	err = shared.SingleCall("Node.FindSuccessor", (randomNode + n.RpcPort), node, r)
	if err != nil {
		n.logger.Error(err.Error())
		return
	}

	n.prev.Id = n.id
	n.prev.Ip = n.ip
	n.table.fingers[1].node = r.Next
	n.logger.Debug("Finished joining, succ : " + r.Next.Ip)
}

func (n *Node) closestFinger(id big.Int) (shared.NodeInfo, error) {
	r := &shared.Reply{}
	args := shared.Test {
		Id: id.Bytes() }

	err := n.ClosestPrecedingFinger(args, r)
	if err != nil {
		err = shared.SingleCall("Node.ClosestPrecedingFinger", (n.table.fingers[1].node.Ip + n.RpcPort), args, r)
	}

	return r.Next, err
}

func (n *Node) getSucc(ip string) (shared.NodeInfo, error) {
	args := 0
	r := &shared.Test{}

	err := shared.SingleCall("Node.GetSuccessor", (ip + n.RpcPort), args, r)

	tmp := new(big.Int)
	tmp.SetBytes(r.Id)
	retVal := shared.NodeInfo {
		Ip: r.Ip,
		Id: *tmp }
	return retVal, err
}

func (n *Node) findPreDecessor(id big.Int) (*shared.Reply, error) {
	var err error

	self := shared.NodeInfo {
		Ip: n.ip,
		Id: n.id }
	currNode := self
	succ := n.table.fingers[1].node

	r := &shared.Reply{}

	for {
		ok, str := util.InKeySpace(currNode.Id, succ.Id, id);
		if !ok {
			break
		}

		succ, err = n.getSucc(currNode.Ip)
		if err != nil {
			n.logger.Error(err.Error())
			return nil, err
		}

		currNode, err = n.closestFinger(id)
		if err != nil {
			n.logger.Debug("STR: " + str)
			n.logger.Error(err.Error())
			return nil, err
		}
	}

	r.Next = succ
	r.Prev = currNode

	return r, nil
}

