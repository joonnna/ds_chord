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
		n.table.fingers[i].start = calcStart((i-1), lenOfId, n.id)
	}
	n.logger.Debug("Inited finger table")
}

func (n *Node) fixFingers() {
	for {

		if n.table.fingers[1].node.Ip == n.Ip {
			time.Sleep(time.Second * 1)
			continue
		}

		index := rand.Int() % lenOfId
		if index == 1 || index == 0 {
			index = 2
		}

		node, err := n.findPreDecessor(n.table.fingers[index].start)
		if err != nil {
			n.logger.Error("FAILED TO SET SUCCESOR")
			n.logger.Error("Succ : " + n.table.fingers[1].node.Ip)
			n.logger.Error(err.Error())
			time.Sleep(time.Second*1)
			continue
		}

		succ, err := n.getSucc(node.Ip)
		if err != nil {
			time.Sleep(time.Second*1)
			continue
		}
		//n.logger.Info("SET SUCCESSOR")
		n.table.fingers[index].node = succ
		time.Sleep(time.Second*1)
	}
}

func(n *Node) stabilize() {
	var succ string
//	n.logger.Info("STABILIZING")
	arg := 0

	succ = n.table.fingers[1].node.Ip

	r := &shared.Test{}
	if succ == n.Ip {
		return
	}
	err := shared.SingleCall("Node.GetPreDecessor", (succ + n.RpcPort), arg, r)
	if err != nil {
		n.logger.Error(err.Error())
		n.logger.Error("Succ : " + succ)
		return
	}

//	n.logger.Debug("Found predecessor and want to stabilize, pre : " + r.Ip)
	tmp := new(big.Int)
	tmp.SetBytes(r.Id)

	if util.BetweenNodes(n.id, n.table.fingers[1].node.Id, *tmp) {
//		n.logger.Debug("New succ is :" + r.Ip)
		n.table.fingers[1].node.Ip = r.Ip
		n.table.fingers[1].node.Id = *tmp
		/*
		if n.table.fingers[1].node.Ip == n.Ip {
			n.logger.Debug("Found myself as succ, sleeping")
			return
		}
		*/
	}

	args := shared.Test {
		Ip: n.Ip,
		Id: n.id.Bytes()}

	reply := &shared.Reply{}
//	n.logger.Info("SENDING NOTIFY")
//	n.logger.Debug("SUCC : " + n.table.fingers[1].node.Ip)
	err = shared.SingleCall("Node.Notify", (n.table.fingers[1].node.Ip + n.RpcPort), args, reply)
	if err != nil {
		n.logger.Error(err.Error())
	}
}

func (n *Node) join() {
	n.putIp()

	randomNode, err := util.GetNode(n.Ip, n.NameServer)
	if err != nil {
		n.logger.Error(err.Error())
		return
	}

	if randomNode == "" {
		self := shared.NodeInfo {
			Id: n.id,
			Ip: n.Ip }
		n.table.fingers[1].node = self
		n.prev = self
		return
	}
	node := shared.Test {
		Ip: n.Ip,
		Id: n.id.Bytes() }

	r := &shared.Reply{}
	n.logger.Debug(randomNode)
	err = shared.SingleCall("Node.FindSuccessor", (randomNode + n.RpcPort), node, r)
	if err != nil {
		n.logger.Error(err.Error())
		return
	}

	n.prev.Id = n.id
	n.prev.Ip = n.Ip
	n.table.fingers[1].node = r.Next
}

func (n *Node) closestFinger(id big.Int) shared.NodeInfo {
	for i := (lenOfId-1); i >= 1; i-- {
		entry := n.table.fingers[i].node.Id
		if entry.BitLen() != 0 && util.BetweenNodes(n.id, id, entry) {
			return n.table.fingers[i].node
		}
	}
	self := shared.NodeInfo {
		Ip: n.Ip,
		Id: n.id }

	return self
}

func (n *Node) getSucc(ip string) (shared.NodeInfo, error) {
	var err error
	args := 0
	r := &shared.Test{}

	if (ip == n.Ip) {
		return n.table.fingers[1].node, nil
	} else {
		err = shared.SingleCall("Node.GetSuccessor", (ip + n.RpcPort), args, r)
	}

	tmp := new(big.Int)
	tmp.SetBytes(r.Id)
	retVal := shared.NodeInfo {
		Ip: r.Ip,
		Id: *tmp }
	return retVal, err
}

func (n *Node) findPreDecessor(id big.Int) (shared.NodeInfo, error) {
	var err error
	self := shared.NodeInfo {
		Ip: n.Ip,
		Id: n.id }
	currNode := self
	succ := n.table.fingers[1].node

	r := &shared.Reply{}
	args := shared.Test {
		Id: id.Bytes() }

	for {
		if util.InKeySpace(currNode.Id, succ.Id, id) {
			break
		}

		if currNode.Ip == n.Ip {
			currNode = n.closestFinger(id)
		} else {
			err = shared.SingleCall("Node.ClosestPrecedingFinger", (currNode.Ip + n.RpcPort), args, r)
			if err != nil {
				return currNode, err
			}
			currNode = r.Next
		}

		if currNode.Ip == n.Ip {
			succ = n.table.fingers[1].node
		} else {
			succ, err = n.getSucc(currNode.Ip)
			if err != nil {
				return currNode, err
			}
		}
	}

	return currNode, nil
}

