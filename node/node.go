package node

import (
//	"fmt"
	"log"
	"strings"
	"os"
	"strconv"
	"github.com/joonnna/ds_chord/node_communication"
	"github.com/joonnna/ds_chord/util"
	"github.com/joonnna/ds_chord/logger"
)

type Node struct {
	storage map[string]string
	ip string
	id int
	nameServer string
	logger *logger.Logger
	next *shared.Comm
	prev *shared.Comm
}

func getNode(list []string, curNode string) string {
	for _, ip := range list {
		if ip != curNode {
			return ip
		}
	}
	return ""
}

func (n *Node) joinNetwork() {

	n.putIp()
	list := n.getNodeList()

	if len(list) == 1 && list[0] == n.ip {
		n.logger.Info("No other nodes in network")
	} else {
		randomNode := getNode(list, n.ip)
		n.logger.Info("Contacting " + randomNode)

		client, err := shared.DialNeighbour(randomNode)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}
		temp := &shared.Comm{Client: client}

		r, err := temp.FindSuccessor(n.id, 10)
		if err != nil {
			n.logger.Error(err.Error())
		}
		n.logger.Info("My successor is " + r)
	}

}

func (n *Node) initNode(ip string, nameServer string, id int) {
	n.storage = make(map[string]string)
	n.ip = ip
	n.nameServer = "http://" + nameServer + ":7551"
	n.id = id

	l := new(logger.Logger)
	l.Init((os.Stdout), ip, 0)
	n.logger = l

	n.logger.Info("Started node")

	shared.InitRpcServer(n.ip, n)

	n.joinNetwork()
}


func (n *Node) FindSuccessor(id int, reply *shared.ReplyType) error {
	n.logger.Debug("FindSuccessor on id " + strconv.Itoa(id) + " on node " +  n.ip)
	if id < n.id {
		temp := *reply
		temp.Ip = n.ip
	} else {
		n.logger.Info("Found no successor")
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
	n.initNode(hostName, nameServer, nodeId)

	n.nodeHttpHandler()

}
