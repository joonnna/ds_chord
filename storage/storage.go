package storage

import (
	"github.com/joonnna/ds_chord/util"
	"github.com/joonnna/ds_chord/chord"
	"github.com/joonnna/ds_chord/logger"
	"runtime"
	"net"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os"
	"github.com/gorilla/mux"
	"time"
)


type Storage struct {
	chord *chord.Chord
	log *logger.Logger
}

/* Handles put requests
   Finds and stores the given key/value on the appropriate node. */
func (s *Storage) putHandler(w http.ResponseWriter, r *http.Request) {
	key := util.GetKey(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.log.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	value := string(body)

	successor, err := s.chord.FindSuccessor(key)
	if err != nil {
		s.log.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = s.chord.PutKey(successor, key, value)
	if err != nil {
		s.log.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
}



/* Handles put requests
   Finds and retrieves the given key on the appropriate node. */
func (s *Storage) getHandler(w http.ResponseWriter, r *http.Request) {
	key := util.GetKey(r)

	successor, err := s.chord.FindSuccessor(key)
	if err != nil {
		s.log.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	value, err := s.chord.GetKey(successor, key)
	if err != nil {
		s.log.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(value)
	if err != nil {
		s.log.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
	}
}

/* Responsible for handling http requests*/
func (s *Storage) httpHandler(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/{key}", s.getHandler).Methods("GET")
	r.HandleFunc("/{key}", s.putHandler).Methods("PUT")

	l, err := net.Listen("tcp4", port)
	if err != nil {
		s.log.Error(err.Error())
		os.Exit(1)
	}
	defer l.Close()

	err = http.Serve(l, r)
	if err != nil {
		s.log.Error(err.Error())
		os.Exit(1)
	}
}
/* Inits and runs the chord implementation, responsible for handling requests*/
func Run(nameServer, httpPort, rpcPort string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.DefaultTransport.(*http.Transport).IdleConnTimeout = time.Second * 1
	http.DefaultTransport.(*http.Transport).MaxIdleConns = 10000

	l := new(logger.Logger)
	l.Init((os.Stdout), "Storage", 0)

	storage := &Storage{
		chord: chord.Init(nameServer, httpPort, rpcPort),
		log: l}

	go storage.chord.Run()

	storage.httpHandler(httpPort)
}
