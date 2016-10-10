package shared

import (
	"math/big"
)

type NodeInfo struct {
	Ip string
	Id big.Int
}

type Search struct {
	Ip string
	Id []byte
	Value string
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


/* All rpc functions */
type RPC interface {
	FindSuccessor(args Search, reply *Reply) error
	Notify(args Search, reply *Reply) error
	ClosestPrecedingFinger(args Search, reply *Reply) error
	PutKey(args Search, reply *Reply) error
	GetKey(args Search, reply *Reply) error
	GetPreDecessor(args int, reply *Search) error
	GetSuccessor(args int, reply *Search) error
}
