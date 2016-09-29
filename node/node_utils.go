package node

import (
	"strings"
	"github.com/joonnna/ds_chord/logger"
)


func (n *Node) assertSuccessor(newSucc string) {
	cmp := strings.Compare(n.next.Id, newSucc)
	if cmp == 0 {
		n.logger.Error("Invalid successor")
	}

	c := strings.Compare(n.next.Id, n.id)
	if c == 1 && cmp == -1 {
		n.logger.Error("Invalid successor")
	}
}

func (n *Node) assertPreDecessor(newPre string) {
	cmp := strings.Compare(n.prev.Id, newPre)
	if cmp == 0 {
		n.logger.Error("Invalid predecessor")
	}

	c := strings.Compare(n.prev.Id, n.id)
	if c == -1 && cmp == 1 {
		n.logger.Error("Invalid predecessor")
	}
}

func createArgs(address string, nodeAddr string, nodeId string) shared.Args {
	n := shared.NodeInfo{
		Ip: nodeAddr,
		Id: nodeId }

	args := shared.Args{
		Address: address,
		Node: n }

	return args
}

func updateArgs(address, id, prevId string) shared.UpdateArgs {
	args := shared.UpdateArgs {
		Address: address,
		Id: id,
		PrevId: prevId }

	return args
}
