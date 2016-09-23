package shared

import(
	"fmt"
	"net/rpc"
	"net"
	"log"
	//"time"
)

type Comm struct {
	Client *rpc.Client
}


func (c *Comm) FindSuccessor(id int, reply int) string {
	r := new(ReplyType)
	err := c.Client.Call("Node.FindSuccessor", id, r)
	if err != nil {
		fmt.Println("Javell ja...")
	}
	fmt.Println(r.Ip)
	return r.Ip
}

func InitRpcServer(ip string, api RPC) {
	server := rpc.NewServer()

	server.RegisterName("Node", api)

	l, err := net.Listen("tcp", ":8245")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Initing RPC on node : " + ip)

	go server.Accept(l)
}

func DialNeighbour(ip string) *rpc.Client {
	//timeout := time.Duration(5 *time.Second)
	connection, err := net.Dial("tcp", ip + ":8245")
	if err != nil {
		log.Fatal(err)
	}
	return rpc.NewClient(connection)
}


