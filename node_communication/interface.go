package shared

import (
	"math/big"
)

type NodeInfo struct {
	Ip string
	Id big.Int
}

type Test struct {
	Ip string
	Id []byte
}

type Reply struct {
	Next NodeInfo
	Prev NodeInfo
	Value string
}

type Args struct {
	Node NodeInfo
	Key big.Int
	Value string
}

type UpdateReply struct {
	Values map[string]string
}

type UpdateArgs struct {
	Id big.Int
	PrevId string
}


type RPC interface {
	FindSuccessor(args Test, reply *Reply) error
//	UpdateSuccessor(args Args, reply *Reply) error
//	UpdatePreDecessor(args Args, reply *Reply) error
	Notify(args Test, reply *Reply) error
	ClosestPrecedingFinger(args Test, reply *Reply) error
	PutKey(args Args, reply *Reply) error
	GetKey(args Args, reply *Reply) error
	GetPreDecessor(args int, reply *Test) error
	GetSuccessor(args int, reply *Test) error
	//	SplitKeys(args UpdateArgs, reply *UpdateReply) error
}
