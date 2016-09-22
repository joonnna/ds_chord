package shared

type RPC interface {
	FindSuccessor(id int, test *int) error
}
