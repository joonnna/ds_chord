package node

import (
	"io/ioutil"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
	"github.com/joonnna/ds_chord/util"
	"github.com/joonnna/ds_chord/node_communication"
)

func rpcArgs(ip string, key int, value string) shared.Args {
	args := shared.Args {
		Key: key,
		Value: value,
		Address: ip }

	return args
}

func (n *Node) nodeHttpHandler() {
	r := mux.NewRouter()
	r.HandleFunc("/{key}", n.getHandler).Methods("GET")
	r.HandleFunc("/{key}", n.putHandler).Methods("PUT")

	http.ListenAndServe(":8080", r)
}


func (n *Node) getHandler(w http.ResponseWriter, r *http.Request) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	key, err := strconv.Atoi(util.GetKey(r))
	if err != nil {
		n.logger.Error(err.Error())
		return
	}

	n.logger.Info("GET key " + strconv.Itoa(key) + " on nodeid " + strconv.Itoa(n.id))

	var value string

	if util.InKeySpace(key, n.id, n.prev.Id) {
		value = n.storage[key]
	} else {
		args := createArgs(n.next.Ip, n.ip, key)
		reply, err := shared.SingleCall("Node.Findsuccessor", args)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}

		args = rpcArgs(reply.Next.Ip, key, value)
		reply, err = shared.SingleCall("Node.GetKey", args)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}
		value = reply.Value
	}

	err = json.NewEncoder(w).Encode(value)
	if err != nil {
		n.logger.Error(err.Error())
	}

}

func (n *Node) putHandler(w http.ResponseWriter, r *http.Request) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	key, err := strconv.Atoi(util.GetKey(r))
	if err != nil {
		n.logger.Error(err.Error())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		n.logger.Error(err.Error())
		return
	}
	value := string(body)
	n.logger.Info("PUT key/value " + strconv.Itoa(key) + "/" + value + " on nodeid " + strconv.Itoa(n.id))

	if util.InKeySpace(key, n.id, n.prev.Id) {
		n.storage[key] = value
		return
	}

	args := createArgs(n.next.Ip, n.ip, key)
	reply, err := shared.SingleCall("Node.Findsuccessor", args)
	if err != nil {
		n.logger.Error(err.Error())
		return
	}

	args = rpcArgs(reply.Next.Ip, key, value)
	_, err = shared.SingleCall("Node.PutKey", args)
	if err != nil {
		n.logger.Error(err.Error())
		return
	}
}

func (n *Node) putIp() {
	req, err := http.NewRequest("PUT", n.nameServer+"/", strings.NewReader(n.ip))
	if err != nil {
		log.Fatal(err)
	}

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{Timeout : timeout}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Panic in node")
		log.Fatal(err)
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
		log.Fatal(err)
	}

	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	err = json.Unmarshal(body, &nodeIps)
	if err != nil {
		log.Fatal(err)
	}
	return nodeIps
}
