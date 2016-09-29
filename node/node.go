package node

import (
	//"fmt"
	"net"
	"strings"
	"os"
	"github.com/joonnna/ds_chord/node_communication"
	"github.com/joonnna/ds_chord/util"
	"github.com/joonnna/ds_chord/logger"
	"github.com/joonnna/ds_chord/storage"
	"sync"
	"errors"
)

type Node struct {
	storage map[string]string
	ip string
	id string
	nameServer string
	rpcPort string
	logger *logger.Logger

	next shared.NodeInfo
	prev shared.NodeInfo

	mutex sync.RWMutex
	update sync.Mutex
}

var (
	ErrGet = errors.New("Unable to get key")
	ErrPut = errors.New("Unable to put key")
)

func (n *Node) initNode() net.Listener {
	n.next.Id = ""
	n.prev.Id = ""

	l := new(logger.Logger)
	l.Init((os.Stdout), n.ip, 0)
	n.logger = l

	n.logger.Info("Started node")

	listener, err := shared.InitRpcServer(n.ip + n.rpcPort, n)
	if err != nil {
		n.logger.Error(err.Error())
		os.Exit(1)
	}

	return listener
}



func (n *Node) joinNetwork() {
	n.putIp()
	list := n.getNodeList()

	if len(list) == 1 && list[0] == n.ip {
		n.logger.Info("No other nodes in network")
	} else {
		randomNode := util.GetNode(list, n.ip)
		n.logger.Info("Contacting " + randomNode)

		args := createArgs((randomNode + n.rpcPort), n.ip, n.id)
		r, err := shared.SingleCall("Node.FindSuccessor", args)
		if err != nil {
			n.logger.Debug("Args " + randomNode)
			n.logger.Error(err.Error())
			return
		}
		n.next = r.Next
		n.prev = r.Prev

		n.logger.Info("My successor is " + n.next.Ip)
		n.logger.Info("My Pre-descessor is " + n.prev.Ip)

		args = createArgs((n.next.Ip + n.rpcPort), n.ip, n.id)
		r, err = shared.SingleCall("Node.UpdateSuccessor", args)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}

		args = createArgs((n.prev.Ip + n.rpcPort), n.ip, n.id)
		r, err = shared.SingleCall("Node.UpdatePreDecessor", args)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}
		args = updateArgs((n.next.Ip + n.rpcPort), n.id)
		r, err = shared.SingleCall("Node.SplitKeys", args)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}
		n.storage = r.storage
	}
}

func (n *Node) PutKey(args shared.Args, reply *shared.Reply) error {
	if util.InKeySpace(n.id, args.Key, n.prev.Id) {
		n.storage[args.Key] = args.Value
		return nil
	} else {
		n.logger.Error("PUT Wrong node " + args.Key)
		return ErrPut
	}
}

func (n *Node) GetKey(args shared.Args, reply *shared.Reply) error {
	if util.InKeySpace(n.id, args.Key, n.prev.Id) {
		reply.Value = n.storage[args.Key]
		return nil
	} else {
		n.logger.Error("GET Wrong node " + args.Key)
		return ErrGet
	}
}

func (n *Node) UpdateSuccessor(args shared.Args, info *shared.Reply) error {
	n.update.Lock()

	n.assertPreDecessor(args.Node.Id)
	n.prev = args.Node
	n.logger.Info("New Pre-descessor " + n.prev.Ip)

	n.update.Unlock()
	return nil
}

func (n *Node) UpdatePreDecessor(args shared.Args, info *shared.Reply) error {
	n.update.Lock()

	n.assertSuccessor(args.Node.Id)
	n.next = args.Node
	n.logger.Info("New successor " + n.next.Ip)

	n.update.Unlock()
	return nil
}

func (n *Node) FindSuccessor(args shared.Args, reply *shared.Reply) error {
	//n.logger.Debug("FindSuccessor on id " + args.Node.Id)
	if n.next.Id == "" && n.prev.Id == "" {
		reply.Next.Ip = n.ip
		reply.Next.Id = n.id
		reply.Prev.Ip = n.ip
		reply.Prev.Id = n.id
		//n.logger.Debug("Second node case, linking circle")
	} else if util.InKeySpace(n.id, args.Node.Id, n.prev.Id) {
	//	n.logger.Debug("In my keyspace")
		reply.Next.Ip = n.ip
		reply.Next.Id = n.id
		reply.Prev = n.prev
	} else {
	//	n.logger.Debug("Not in my keyspace")

		args := createArgs((n.next.Ip + n.rpcPort), args.Node.Ip, args.Node.Id)
		r, err := shared.SingleCall("Node.FindSuccessor", args)
		if err != nil {
			n.logger.Error(err.Error())
			return err
		}

		reply.Next = r.Next
		reply.Prev = r.Prev
		reply.Value = r.Value
	}
	return nil
}

func (n *Node) SplitKeys(args UpdateArgs, reply *UpdateReply) error {
	reply.Values := storage.SplitStorage(args.Id, args.PrevId, n.storage)
	return nil
}

func Run(nameServer string, httpPort string, rpcPort string) {
	go util.CheckInterrupt()

	hostName, _ := os.Hostname()
	hostName = strings.Split(hostName, ".")[0]

	hashId := util.HashKey(hostName)

	n := &Node {
		id: hashId,
		ip: hostName,
		storage: make(map[string]string),
		nameServer: "http://" + nameServer + httpPort,
		rpcPort: rpcPort }

	l := n.initNode()
	defer l.Close()

	n.joinNetwork()
	n.nodeHttpHandler(httpPort)
	n.logger.Debug("YOYOYOYOYOYOYOYOYYO")
}
