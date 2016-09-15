package main

import (
	"strings"
	"os"
	"log"
	"fmt"
	"net"
	"net/http"
	"github.com/gorilla/mux"
	"sync"
	"encoding/json"
	"io/ioutil"
)

type State struct {
	nodeIps []string
	mutex sync.RWMutex
}

func HttpServer(ip string) {
	r := mux.NewRouter()
	current_state := new(State)

	r.Methods("GET").Path("/").HandlerFunc(current_state.getHandler)
	r.Methods("PUT").Path("/").HandlerFunc(current_state.putHandler)

	fmt.Printf("Server listening on %s...\n", ip)

	http.ListenAndServe(":8080", r)
}



func (s *State) getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received get in nameserver")
	s.mutex.RLock()

	err := json.NewEncoder(w).Encode(s.nodeIps)

	if err != nil {
		log.Fatal(err)
		s.mutex.RUnlock()
	}

	s.mutex.RUnlock()
}

func (s *State) putHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Received put in nameserver")
	s.mutex.Lock()

	newIp, err:= ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	s.nodeIps = append(s.nodeIps, string(newIp))
	fmt.Println(s.nodeIps)
	s.mutex.Unlock()
}

func main() {
	hostName, _ := os.Hostname()
	hostName = strings.Split(hostName, ".")[0]
	fmt.Println("Started nameserver on " + hostName)
	HttpServer(hostName)
}
