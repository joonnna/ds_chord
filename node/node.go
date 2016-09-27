package node

import (
//	"fmt"
	"net"
	"log"
	"strings"
	"os"
	"strconv"
	"github.com/joonnna/ds_chord/node_communication"
	"github.com/joonnna/ds_chord/util"
	"github.com/joonnna/ds_chord/logger"
	"sync"
	"errors"
)

type Node struct {
	storage map[int]string
	ip string
	id int
	nameServer string
	logger *logger.Logger
	next shared.NodeInfo
	prev shared.NodeInfo
	mutex sync.RWMutex
}
var (
	ErrGet = errors.New("Unable to get key")
	ErrPut = errors.New("Unalbe to put key")
)

func (n *Node) initNode(ip string, nameServer string, id int) net.Listener {
	n.storage = make(map[int]string)
	n.ip = ip
	n.nameServer = "http://" + nameServer + ":7551"
	n.id = id
	n.next.Id = -1
	n.prev.Id = -1

	l := new(logger.Logger)
	l.Init((os.Stdout), n.ip, 0)
	n.logger = l

	n.logger.Info("Started node")

	listener, err := shared.InitRpcServer(n.ip, n)
	if err != nil {
		n.logger.Error(err.Error())
		os.Exit(1)
	}

	return listener
}

func createArgs(address string, nodeAddr string, nodeId int) shared.Args {
	n := shared.NodeInfo{
		Ip: nodeAddr,
		Id: nodeId }

	args := shared.Args{
		Address: address,
		Node: n }

	return args
}


func (n *Node) joinNetwork() {
	n.putIp()
	list := n.getNodeList()

	if len(list) == 1 && list[0] == n.ip {
		n.logger.Info("No other nodes in network")
	} else {
		randomNode := util.GetNode(list, n.ip)
		n.logger.Info("Contacting " + randomNode)

		args := createArgs(randomNode, n.ip, n.id)
		r, err := shared.SingleCall("Node.FindSuccessor", args)
		if err != nil {
			n.logger.Debug("Args " + randomNode)
			n.logger.Error(err.Error())
			return
		}
		n.next = r.Next
		n.prev = r.Prev

		n.logger.Info("My successor is " + strconv.Itoa(n.next.Id))
		n.logger.Info("My Pre-descessor is " + strconv.Itoa(n.prev.Id))

		args = createArgs(n.next.Ip, n.ip, n.id)
		r, err = shared.SingleCall("Node.UpdateSuccessor", args)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}

		args = createArgs(n.prev.Ip, n.ip, n.id)
		r, err = shared.SingleCall("Node.UpdatePreDecessor", args)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}
	}
}

func (n *Node) PutKey(args shared.Args, reply *shared.Reply) error {
	if util.InKeySpace(args.Key, n.id, n.prev.Id) {
		n.storage[args.Key] = args.Value
		return nil
	} else {
		n.logger.Error("PUT Wrong node, nodeid " + strconv.Itoa(n.id) + " key " + strconv.Itoa(args.Key))
		return ErrPut
	}
}

func (n *Node) GetKey(args shared.Args, reply *shared.Reply) error {
	if util.InKeySpace(args.Key, n.id, n.prev.Id) {
		reply.Value = n.storage[args.Key]
		return nil
	} else {
		n.logger.Error("GET Wrong node, nodeid " + strconv.Itoa(n.id) + " key " + strconv.Itoa(args.Key))
		return ErrGet
	}
}

func (n *Node) UpdateSuccessor(args shared.Args, info *shared.Reply) error {
	n.prev = args.Node
	n.logger.Info("New Pre-descessor " + strconv.Itoa(n.prev.Id))
	return nil
}

func (n *Node) UpdatePreDecessor(args shared.Args, info *shared.Reply) error {
	n.next = args.Node
	n.logger.Info("New successor " + strconv.Itoa(n.next.Id))
	return nil
}

func (n *Node) FindSuccessor(args shared.Args, reply *shared.Reply) error {
	n.logger.Debug("FindSuccessor on id " + strconv.Itoa(args.Node.Id))
	if n.next.Id == -1 && n.prev.Id == -1 {
		reply.Next.Ip = n.ip
		reply.Next.Id = n.id
		reply.Prev.Ip = n.ip
		reply.Prev.Id = n.id
		n.logger.Debug("Second node case, linking circle")
	} else if util.InKeySpace(args.Node.Id, n.id, n.prev.Id) {
		n.logger.Debug("In my keyspace")
		reply.Next.Ip = n.ip
		reply.Next.Id = n.id
		reply.Prev = n.prev
	} else {
		n.logger.Debug("Not in my keyspace")

		args := createArgs(n.next.Ip, args.Node.Ip, args.Node.Id)
		r, err := shared.SingleCall("Node.FindSuccessor", args)
		if err != nil {
			n.logger.Error(err.Error())
			return err
		}

		//reply = r
		reply.Next = r.Next
		reply.Prev = r.Prev
		reply.Value = r.Value
	}
	return nil
}


func Run(nameServer string, id string) {
	go util.CheckInterrupt()

	hostName, _ := os.Hostname()
	hostName = strings.Split(hostName, ".")[0]

	n := new(Node)
	nodeId, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}
	l := n.initNode(hostName, nameServer, nodeId)
	defer l.Close()

	n.joinNetwork()
	n.nodeHttpHandler()
}
