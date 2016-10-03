package node

import (
	//"fmt"
	"time"
	"net"
	"strings"
	"os"
	"github.com/joonnna/ds_chord/node_communication"
	"github.com/joonnna/ds_chord/util"
	"github.com/joonnna/ds_chord/logger"
	"sync"
)

type Node struct {
	data map[string]string
	ip string
	id string
	NameServer string
	RpcPort string
	logger *logger.Logger

	listener net.Listener
	succList []shared.NodeInfo

	Next shared.NodeInfo
	prev shared.NodeInfo

	update sync.Mutex
}

func InitNode(nameServer, httpPort, rpcPort string) *Node {
	hostName, _ := os.Hostname()
	hostName = strings.Split(hostName, ".")[0]

	hashId := util.HashKey(hostName)

	log := new(logger.Logger)
	log.Init((os.Stdout), hostName, 0)

	n := &Node {
		id: hashId,
		ip: hostName,
		logger: log,
		NameServer: "http://" + nameServer + httpPort,
		RpcPort: rpcPort,
		data: make(map[string]string) }



	l, err := shared.InitRpcServer(hostName + rpcPort, n)
	if err != nil {
		n.logger.Error(err.Error())
		os.Exit(1)
	}

	n.listener = l
	n.Next.Id = ""
	n.prev.Id = ""


	n.logger.Info("Started node")

	return n
}


func (n *Node) joinNetwork() {
	n.putIp()
	list, err := GetNodeList(n.NameServer)
	if err != nil {
		n.logger.Error(err.Error())
		return
	}

	if len(list) == 1 && list[0] == n.ip {
		n.logger.Info("No other nodes in network")
	} else {
		randomNode := util.GetNode(list, n.ip)
		n.logger.Info("Contacting " + randomNode)


		r := &shared.Reply{}
		args := util.CreateArgs(n.ip, n.id)
		err := shared.SingleCall("Node.FindSuccessor", (randomNode + n.RpcPort), args, r)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}
		n.Next = r.Next
		n.prev = r.Prev

		n.logger.Info("My successor is " + n.Next.Ip)
		n.logger.Info("My Pre-descessor is " + n.prev.Ip)

		/*
		arguments := updateArgs(n.id, n.prev.Id)
		r, err = shared.SingleCall("Node.SplitKeys", (n.next.Ip + n.rpcPort), arguments)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}
		//n.store = r.store

		*/
		args = util.CreateArgs(n.ip, n.id)
		err = shared.SingleCall("Node.UpdateSuccessor", (n.Next.Ip + n.RpcPort), args, r)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}

		args = util.CreateArgs(n.ip, n.id)
		err = shared.SingleCall("Node.UpdatePreDecessor", (n.prev.Ip + n.RpcPort), args, r)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}
	}
}


func Run(n *Node) {
	defer n.listener.Close()

	n.joinNetwork()
	for {
		time.Sleep(time.Second * 5)
	}
}
