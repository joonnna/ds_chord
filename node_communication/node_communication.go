package main

import (
	"shared"
)

type Comm struct {
	client *rpc.Client
}


func (c *Comm) FindSuccessor(id int, test *string) error {
	fmt.Println("yoyoyoyo")

	c.client.Call("Node.FinndSuccessor", id)

	return nil
}



