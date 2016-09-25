package nameserver

import (
	"github.com/joonnna/ds_chord/util"
	"github.com/joonnna/ds_chord/logger"
	"strings"
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"sync"
	"encoding/json"
	"io/ioutil"
	"errors"
)

type nameServer struct {
	nodeIps []string
	mutex sync.RWMutex
	logger *logger.Logger
}

var (
	ErrEncode = errors.New("Unable to encode body")
	ErrRead = errors.New("Unable to read body")
)


func (n *nameServer) Init(ip string) {
	l := new(logger.Logger)
	l.Init((os.Stdout), "Nameserver", 0)

	n.logger = l
}


func (n *nameServer) httpServer() {
	r := mux.NewRouter()

	r.Methods("GET").Path("/").HandlerFunc(n.getHandler)
	r.Methods("PUT").Path("/").HandlerFunc(n.putHandler)

	port := ":7551"

	n.logger.Info("Listening on " + port)

	http.ListenAndServe(port, r)
}



func (n *nameServer) getHandler(w http.ResponseWriter, r *http.Request) {
	n.logger.Info("Received get in nameserver")
	n.mutex.RLock()

	err := json.NewEncoder(w).Encode(n.nodeIps)

	if err != nil {
		n.logger.Error(ErrEncode.Error())
	}

	n.mutex.RUnlock()
}

func (n *nameServer) putHandler(w http.ResponseWriter, r *http.Request) {
	n.logger.Info("Received get in nameserver")
	n.mutex.Lock()

	newIp, err:= ioutil.ReadAll(r.Body)
	if err != nil {
		n.mutex.Unlock()
		n.logger.Error(ErrRead.Error())
		return
	}
	n.nodeIps = append(n.nodeIps, string(newIp))

	n.mutex.Unlock()
}

func Run() {
	go util.CheckInterrupt()

	hostName, _ := os.Hostname()
	hostName = strings.Split(hostName, ".")[0]

	n := new(nameServer)

	n.Init(hostName)

	n.httpServer()
}
