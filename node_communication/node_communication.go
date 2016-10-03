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

func SingleCall(method string, address string, args interface{}, reply interface{})  error {
	//var reply interface{}
	c, err := setupConn(address)
	if err != nil {
		return err
	}

	err = c.Client.Call(method, args, reply)
	if err != nil {
		return err
	}

	c.Client.Close()

	return nil
}
