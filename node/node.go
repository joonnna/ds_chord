package node


import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"encoding/json"
	"log"
	"strings"
)


type Node struct {
	storage map[string]string
	Ip string
	NameServer string
}


func getKey(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["key"]
}


func NodeHttpHandler() {
	r := mux.NewRouter()

	node_state := new(Node)
	node_state.storage = make(map[string]string)

	r.HandleFunc("/{key}", node_state.getHandler).Methods("GET")
	r.HandleFunc("/{key}", node_state.putHandler).Methods("PUT")

	fmt.Println("Server listening...")

	http.ListenAndServe(":8080", r)
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
