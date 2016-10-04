package shared

import (
	"math/big"
)

type NodeInfo struct {
	Ip string
	Id *big.Int
}

type Reply struct {
	Next NodeInfo
	Prev NodeInfo
	Value string
}

type Args struct {
	Node NodeInfo
	Key *big.Int
	Value string
}

type UpdateReply struct {
	Values map[string]string
}

type UpdateArgs struct {
	Id *big.Int
	PrevId string
}


type RPC interface {
	FindSuccessor(args Args, reply *Reply) error
//	UpdateSuccessor(args Args, reply *Reply) error
//	UpdatePreDecessor(args Args, reply *Reply) error
	Notify(args Args, reply *Reply) error
	ClosestPrecedingFinger(args Args, reply *Reply) error
	PutKey(args Args, reply *Reply) error
	GetKey(args Args, reply *Reply) error
//	SplitKeys(args UpdateArgs, reply *UpdateReply) error
}
