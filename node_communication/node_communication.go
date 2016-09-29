package shared

import(
//	"fmt"
	"net/rpc"
	"net"
	//"log"
	//"errors"
	//"time"
)

type Comm struct {
	Client *rpc.Client
}
func InitRpcServer(address string, api RPC) (net.Listener, error) {
	server := rpc.NewServer()
	err := server.RegisterName("Node", api)
	if err != nil {
		return nil, err
	}

	t, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}

	l, err := net.ListenTCP("tcp4", t)
	if err != nil {
		return nil, err
	}

	go server.Accept(l)

	return l, nil
}

func setupConn(address string) (*Comm, error) {
	client, err := dialNode(address)
	if err != nil {
		return nil, err
	}

	return &Comm{Client: client}, nil
}

func dialNode(address string) (*rpc.Client, error) {
	t, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}

	connection, err := net.DialTCP("tcp4", nil,  t)
	if err != nil {
		return nil, err
	}

	return rpc.NewClient(connection), nil
}


func SingleCall(method string, args Args) (*Reply, error) {
	reply := &Reply{}

	c, err := setupConn(args.Address)
	if err != nil {
		return nil, err
	}

	err = c.Client.Call(method, args, reply)
	if err != nil {
		return nil, err
	}

	c.Client.Close()

	return reply, nil
}
