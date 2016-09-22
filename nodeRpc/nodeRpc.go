package nodeRpc

import(
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func InitRpcServer(ip string) {
	server := rpc.NewServer()

	server.RegisterName("Node", server)

	l, err := net.Listen("tcp", ip + ":8005")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Initing RPC on node : " + ip)

	go server.Accept(l)
}

func DialNeighbour(ip string) *rpc.Client {
	connection, err := net.Dial("tcp", ip)
	if err != nil {
		return nil
	}

	return rpc.NewClient(connection)
}

