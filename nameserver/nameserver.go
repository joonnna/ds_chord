package nameserver

import (
	"github.com/joonnna/ds_chord/logger"
	"github.com/joonnna/ds_chord/util"
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
	ip string
	port string
}

var (
	ErrEncode = errors.New("Unable to encode body")
	ErrRead = errors.New("Unable to read body")
)


func (n *nameServer) Init(ip string, port string) {
	l := new(logger.Logger)
	l.Init((os.Stdout), "Nameserver", 0)

	n.ip = ip
	n.logger = l
	n.port = port
}


func (n *nameServer) httpServer() {
	r := mux.NewRouter()

	r.Methods("GET").Path("/").HandlerFunc(n.getHandler)
	r.Methods("PUT").Path("/").HandlerFunc(n.putHandler)

	n.logger.Info("Listening on " + n.port)

	http.ListenAndServe(n.port, r)
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
	n.logger.Info("Received put in nameserver")
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

func Run(port string) {
	go util.CheckInterrupt()

	hostName, _ := os.Hostname()
	hostName = strings.Split(hostName, ".")[0]

	n := new(nameServer)

	n.Init(hostName, port)

	n.httpServer()
}
