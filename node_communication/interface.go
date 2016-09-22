package shared

type Comm interface {
	FindSuccessor(id int, test *string)
}
