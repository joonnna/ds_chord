package shared

import(
	//"fmt"
	"net/rpc"
	"net"
	//"log"
	//"errors"
	//"time"
)

type Comm struct {
	Client *rpc.Client
}


func (c *Comm) FindSuccessor(id int, reply int) (string, error) {
	r := new(ReplyType)
	err := c.Client.Call("Node.FindSuccessor", id, r)
	if err != nil {
		return "", err
	}

	return r.Ip, nil
}

func InitRpcServer(ip string, api RPC) error {
	server := rpc.NewServer()

	server.RegisterName("Node", api)

	l, err := net.Listen("tcp", ":8132")
	if err != nil {
		return err
	}

	go server.Accept(l)

	return nil
}

func DialNeighbour(ip string) (*rpc.Client, error) {
	//timeout := time.Duration(5 *time.Second)
	connection, err := net.Dial("tcp", ip + ":8132")
	if err != nil {
		return nil, err
	}
	return rpc.NewClient(connection), nil
}


