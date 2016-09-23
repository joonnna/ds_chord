package node

import (
	"fmt"
	"log"
	"strings"
	"os"
	"strconv"
	"github.com/joonnna/ds_chord/node_communication"
	"github.com/joonnna/ds_chord/util"
)

type Node struct {
	storage map[string]string
	ip string
	id int
	nameServer string
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

func(n *Node) initNode(ip string, nameServer string, id int) {
	n.storage = make(map[string]string)
	n.ip = ip
	n.nameServer = "http://" + nameServer + ":7551"
	n.id = id

	shared.InitRpcServer(n.ip, n)
	n.putIp()

	list := n.getNodeList()
	if len(list) == 1 && list[0] == n.ip {
		fmt.Println("IM ALONE BITCHES")
	} else {
		randomNode := getNode(list, n.ip)
		fmt.Printf("CONNECTING TO %s\n", randomNode)
		client := shared.DialNeighbour(randomNode)
		temp := &shared.Comm{Client: client}

		r := temp.FindSuccessor(n.id, 10)
		fmt.Println("MY SUCCESSOR IS : " + r)
	}
}


func (n *Node) FindSuccessor(id int, reply *shared.ReplyType) error {
	fmt.Println("FindSuccessor on id " + strconv.Itoa(id) + " on node " +  n.ip)
	if id < n.id {
		temp := *reply
		temp.Ip = n.ip
	} else {
		fmt.Println("NOT SUCCESSOR")
	}
	return nil
}


func Run(nameServer string, id string) {
	go util.CheckInterrupt()

	hostName, _ := os.Hostname()
	hostName = strings.Split(hostName, ".")[0]
	fmt.Println("Started node on " + hostName)

	n := new(Node)
	nodeId, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
	n.initNode(hostName, nameServer, nodeId)

	n.nodeHttpHandler()

}
