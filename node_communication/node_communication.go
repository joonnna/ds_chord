package shared

import(
	"fmt"
	"net/rpc"
)

type Comm struct {
	Client *rpc.Client
}


func (c *Comm) FindSuccessor(id int, test *string) error {
	fmt.Println("yoyoyoyo")

	c.Client.Call("Node.FindSuccessor", id, test)

	return nil
}



