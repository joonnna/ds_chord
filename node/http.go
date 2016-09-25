package node

import (
	"io/ioutil"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"encoding/json"
	"log"
	"strings"
	"time"
	"github.com/joonnna/ds_chord/util"
)

func (n *Node) nodeHttpHandler() {
	r := mux.NewRouter()
	r.HandleFunc("/{key}", n.getHandler).Methods("GET")
	r.HandleFunc("/{key}", n.putHandler).Methods("PUT")

	http.ListenAndServe(":8080", r)
}


func (n *Node) getHandler(w http.ResponseWriter, r *http.Request) {
	key := util.GetKey(r)

	err := json.NewEncoder(w).Encode(key)
	if err != nil {
		log.Fatal(err)
	}
}

func (n *Node) putHandler(w http.ResponseWriter, r *http.Request) {
	key := util.GetKey(r)

	body, _ := ioutil.ReadAll(r.Body)
	value := string(body)

	n.storage[key] = value
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
