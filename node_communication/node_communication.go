package shared

import(
	"fmt"
	"net/rpc"
	"net"
	"log"
)

type Comm struct {
	Client *rpc.Client
}


func (c *Comm) FindSuccessor(id int, test *int) error {
	fmt.Println("yoyoyoyo")
	var smeg int
	c.Client.Call("Node.FindSuccessor", id, &smeg)

	return nil
}

func InitRpcServer(ip string, api RPC) {
	server := rpc.NewServer()

	server.RegisterName("Node", api)

	l, err := net.Listen("tcp", ip + ":8005")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Initing RPC on node : " + ip)

	go server.Accept(l)
}

func DialNeighbour(ip string) *rpc.Client {
	connection, err := net.Dial("tcp", ip + ":8005")
	if err != nil {
		return nil
	}

	return rpc.NewClient(connection)
}


