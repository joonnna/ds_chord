package node

import (
	"net"
	"strings"
	"os"
	"github.com/joonnna/ds_chord/node_communication"
	"github.com/joonnna/ds_chord/util"
	"github.com/joonnna/ds_chord/logger"
	"sync"
	"math/big"
	"net/http"
	"time"
)
/* Node defenition */
type Node struct {
	data map[string]string
	Ip string
	id big.Int
	NameServer string
	RpcPort string
	httpPort string
	logger *logger.Logger

	listener net.Listener

	prev shared.NodeInfo

	table fingerTable

	update sync.RWMutex
}
/* Inits and returns the node object
   nameserver: Ip address of the nameserver
   httpPort: Port for http communication
   rpcPort: Port for rpc communication

   Returns node object
*/
func InitNode(nameServer, httpPort, rpcPort string) *Node {
	hostName, _ := os.Hostname()
	hostName = strings.Split(hostName, ".")[0]
	http.DefaultTransport.(*http.Transport).MaxIdleConns = 1000
	tmp := util.ConvertKey(hostName)

	log := new(logger.Logger)
	log.Init((os.Stdout), hostName, 0)

	id := new(big.Int)
	id.SetBytes(tmp)

	n := &Node {
		id: *id,
		Ip: hostName,
		logger: log,
		NameServer: "http://" + nameServer + httpPort,
		httpPort: httpPort,
		RpcPort: rpcPort,
		data: make(map[string]string) }


	l, err := shared.InitRpcServer(hostName + rpcPort, n)
	if err != nil {
		n.logger.Error(err.Error())
		os.Exit(1)
	}

	n.listener = l

	n.initFingerTable()
	return n
}
/* Joins the network and calls stabilize and fix fingers periodically.
	node: the node to run
*/
func Run(n *Node) {
	defer n.listener.Close()

	n.join()
//	n.add()
//	go n.updateState()

	go n.fixFingers()
	for {
		n.stabilize()
		time.Sleep(time.Second * 1)
	}
}
