package storage

import (
//	"github.com/joonnna/ds_chord/logger"
	"github.com/joonnna/ds_chord/util"
//	"strings"
	"github.com/joonnna/ds_chord/chord"
	"github.com/joonnna/ds_chord/logger"
//	"errors"
	"sync"
	"net"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os"
	"github.com/gorilla/mux"
)

type Storage struct {
	chord *chord.Chord
	mutex sync.RWMutex
	log *logger.Logger
}

/*
func (s *Storage) splitStorage(newId string, prevId string) map[string]string {
	ret := make(map[string]string)

	for key, val := range d.store {
		if util.InKeySpace(newId, key, prevId) {
			ret[key] = val
			delete(d.store, key)
		}
	}
	return ret
}
*/

func (s *Storage) putHandler(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

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



func (s *Storage) getHandler(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

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

func Run(nameServer, httpPort, rpcPort string) {
	go util.CheckInterrupt()

	l := new(logger.Logger)
	l.Init((os.Stdout), "Storage", 0)

	storage := &Storage{
		chord: chord.Init(nameServer, httpPort, rpcPort),
		log: l}

	go storage.chord.Run()

	storage.httpHandler(httpPort)
}
