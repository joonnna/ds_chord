package shared


type NodeInfo struct {
	Ip string
	Id int
}

type Reply struct {
	Next NodeInfo
	Prev NodeInfo
	Value string
}

type Args struct {
	Address string
	Node NodeInfo
	Key int
	Value string
}


type RPC interface {
	FindSuccessor(args Args, reply *Reply) error
	UpdateSuccessor(args Args, reply *Reply) error
	UpdatePreDecessor(args Args, reply *Reply) error
	PutKey(args Args, reply *Reply) error
	GetKey(args Args, reply *Reply) error
}
