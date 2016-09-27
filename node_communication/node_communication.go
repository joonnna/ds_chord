package shared

import(
	"fmt"
	"net/rpc"
	"net"
	//"log"
	//"errors"
	//"time"
)

type Comm struct {
	Client *rpc.Client
}


func (c *Comm) UpdateSuccessor(id int, ip string) error {

	n := NodeInfo {
		Ip: ip,
		Id: id }

	r := Reply{
		Prev: n }

	err := c.Client.Call("Node.UpdateSuccessor", r, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Comm) UpdatePreDecessor(id int, ip string) error {
	n := NodeInfo {
		Ip: ip,
		Id: id }

	r := Reply{
		Next: n }

	err := c.Client.Call("Node.UpdatePreDecessor", r, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Comm) FindSuccessor(id int) (*Reply, error) {
	r := new(Reply)
	err := c.Client.Call("Node.FindSuccessor", id, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func InitRpcServer(ip string, api RPC) (net.Listener, error) {
	server := rpc.NewServer()
	fmt.Println("IP : ", ip)
	err := server.RegisterName("Node", api)
	if err != nil {
		return nil, err
	}

	t, err := net.ResolveTCPAddr("tcp4", ip + ":2000")

	fmt.Println("LISTEN TCP ADDR:", t)
	l, err := net.ListenTCP("tcp4", t)
	if err != nil {
		return nil, err
	}
	//conn, err := l.Accept()
	go server.Accept(l)

	return l, nil
}

func setupConn(ip string) (*Comm, error) {
	client, err := dialNode(ip)
	if err != nil {
		return nil, err
	}

	return &Comm{Client: client}, nil
}

func dialNode(ip string) (*rpc.Client, error) {
	//timeout := time.Duration(10 *time.Second)
	t, err := net.ResolveTCPAddr("tcp4", ip + ":2000")

	if err != nil {
		return nil, err
	}
	fmt.Println("TCP ADDR:", t)
	connection, err := net.DialTCP("tcp4", nil,  t)

	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to:", connection.RemoteAddr())
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
