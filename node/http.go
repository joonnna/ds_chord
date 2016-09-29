package node

import (
	"net"
	"os"
	"io/ioutil"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"log"
	"strings"
	"time"
	"github.com/joonnna/ds_chord/util"
	"github.com/joonnna/ds_chord/node_communication"
)

func rpcArgs(ip string, key string, value string) shared.Args {
	args := shared.Args {
		Key: key,
		Value: value,
		Address: ip }

	return args
}

func (n *Node) nodeHttpHandler(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/{key}", n.getHandler).Methods("GET")
	r.HandleFunc("/{key}", n.putHandler).Methods("PUT")
	/*
	err := http.ListenAndServe(port, r)
	if err != nil {
		n.logger.Error(err.Error())
		os.Exit(1)
	}
	*/

	l, err := net.Listen("tcp4", port)
	if err != nil {
		n.logger.Error(err.Error())
		os.Exit(1)
	}
	defer l.Close()

	err = http.Serve(l, r)
	if err != nil {
		n.logger.Error(err.Error())
		os.Exit(1)
	}
}


func (n *Node) getHandler(w http.ResponseWriter, r *http.Request) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	key := util.GetKey(r)

	n.logger.Info("GET key " + key)

	var value string

	hashKey := util.HashKey(key)

	if util.InKeySpace(n.id, hashKey, n.prev.Id) {
		value = n.storage[hashKey]
	} else {
		args := createArgs((n.next.Ip + n.rpcPort), n.ip, hashKey)
		reply, err := shared.SingleCall("Node.FindSuccessor", args)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}

		args = rpcArgs((reply.Next.Ip + n.rpcPort), hashKey, value)
		reply, err = shared.SingleCall("Node.GetKey", args)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}
		value = reply.Value
	}

	err := json.NewEncoder(w).Encode(value)
	if err != nil {
		n.logger.Error(err.Error())
	}

}

func (n *Node) putHandler(w http.ResponseWriter, r *http.Request) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	key := util.GetKey(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		n.logger.Error(err.Error())
		return
	}
	value := string(body)
	n.logger.Info("PUT key/value " + key + "/" + value)

	hashKey := util.HashKey(key)
	if util.InKeySpace(n.id, hashKey, n.prev.Id) {
		n.storage[hashKey] = value
		return
	}

	args := createArgs((n.next.Ip + n.rpcPort), n.ip, hashKey)
	reply, err := shared.SingleCall("Node.FindSuccessor", args)
	if err != nil {
		n.logger.Error(err.Error())
		return
	}

	args = rpcArgs((reply.Next.Ip + n.rpcPort), hashKey, value)
	_, err = shared.SingleCall("Node.PutKey", args)
	if err != nil {
		n.logger.Error(err.Error())
		return
	}
}

func (n *Node) putIp() {
	req, err := http.NewRequest("PUT", n.nameServer+"/", strings.NewReader(n.ip))
	if err != nil {
		n.logger.Error(err.Error())
	}

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{Timeout : timeout}
	resp, err := client.Do(req)
	if err != nil {
		n.logger.Error(err.Error())
	} else {
		resp.Body.Close()
	}
}

func (n *Node) getNodeList() []string  {
	var nodeIps []string

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{Timeout : timeout}

	r, err := client.Get(n.nameServer)
	if err != nil {
		n.logger.Error(err.Error())
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		n.logger.Error(err.Error())
	}

	err = json.Unmarshal(body, &nodeIps)
	if err != nil {
		n.logger.Error(err.Error())
	}
	return nodeIps
}
