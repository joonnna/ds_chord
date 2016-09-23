package shared

type ReplyType struct {
	Ip string
}

type RPC interface {
	FindSuccessor(id int, reply *ReplyType) error
}
