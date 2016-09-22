package nodeRpc

import(
	"log"
	"net"
	"net/rpc"
)

func InitRpc(ip string) {
	server := rpc.NewServer()

	server.RegisterName("Node", server)

	l, err := net.Listen("tcp", ip + ":8005")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("yoyo")

	go server.Accept(l)
}


