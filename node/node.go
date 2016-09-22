package main


import (
	"fmt"
	"net/http"
//	"net"
	"github.com/gorilla/mux"
	"io/ioutil"
	"encoding/json"
	"log"
	"strings"
	"os"
	"github.com/joonnna/ds_chord/nodeRpc"
	"github.com/joonnna/ds_chord/node_communication"
)


type Node struct {
	storage map[string]string
	Ip string
	NameServer string
	Test *shared.Comm
}

func getKey(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["key"]
}


func (n *Node) NodeHttpHandler() {
	r := mux.NewRouter()
	r.HandleFunc("/{key}", n.getHandler).Methods("GET")
	r.HandleFunc("/{key}", n.putHandler).Methods("PUT")

	fmt.Println("Server listening...")

	http.ListenAndServe(n.Ip, r)
}


func (n *Node) getHandler(w http.ResponseWriter, r *http.Request) {
	key := getKey(r)

	err := json.NewEncoder(w).Encode(key)
	if err != nil {
		log.Fatal(err)
	}
}

func (n *Node) putHandler(w http.ResponseWriter, r *http.Request) {
	key := getKey(r)

	body, _ := ioutil.ReadAll(r.Body)
	value := string(body)

	n.storage[key] = value
}

func (n *Node) PutIp() {
	req, err := http.NewRequest("PUT", n.NameServer+"/", strings.NewReader(n.Ip))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("sending put request..")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else {
		resp.Body.Close()
	}
}

func (n *Node) GetNodeList() []string  {
	fmt.Println("sending GET request..")
	var nodeIps []string

	r, err := http.Get(n.NameServer)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(r.Body)

	err = json.Unmarshal(body, &nodeIps)
	if err != nil {
		log.Fatal(err)
	}
	return nodeIps
}

func (n *Node) FindSuccessor(id int, test *string) error {
	fmt.Println("FindSuccessor on id " + string(id) + "on node " +  n.Ip)

	return nil
}



func main() {
	hostName, _ := os.Hostname()
	hostName = strings.Split(hostName, ".")[0]
	fmt.Println("Started node on " + hostName)

	args := os.Args[1:]
	nameServer := strings.Join(args, "")

	var str string
	n := new(Node)
	n.storage = make(map[string]string)
	n.Ip = hostName
	n.NameServer = "http://" + nameServer + ":8080"

	nodeRpc.InitRpcServer(n.Ip)
	n.PutIp()

	list := n.GetNodeList()
	fmt.Println(list)

	client := nodeRpc.DialNeighbour(list[0])
	n.Test = &shared.Comm{Client: client}

	err := n.Test.FindSuccessor(3, &str)
	fmt.Println(err)

	n.NodeHttpHandler()



}
