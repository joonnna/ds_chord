
package node
/*
import (
	"net"
	"os"
	"io/ioutil"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
//	"log"
	"github.com/joonnna/ds_chord/util"
	"github.com/joonnna/ds_chord/node_communication"
)


func (n *Node) nodeHttpHandler(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/{key}", n.getHandler).Methods("GET")
	r.HandleFunc("/{key}", n.putHandler).Methods("PUT")

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

	var err error
	var value string

	hashKey := util.HashKey(key)


	if util.InKeySpace(n.id, hashKey, n.prev.Id) {
		value, err = n.store.GetData(hashKey)
		if err != nil {
			n.logger.Error(err.Error())
		}
	} else {
		args := createArgs(n.ip, hashKey)
		reply, err := shared.SingleCall("Node.FindSuccessor", (n.next.Ip + n.rpcPort), args)
		if err != nil {
			n.logger.Error(err.Error())
			return
		}

		args = rpcArgs(hashKey, value)
		reply, err = shared.SingleCall("Node.GetKey", (reply.Next.Ip + n.rpcPort), args)
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
		err = n.store.PutData(hashKey, value)
		if err != nil {
			n.logger.Error(err.Error())
		}
		return
	}

	args := createArgs(n.ip, hashKey)
	reply, err := shared.SingleCall("Node.FindSuccessor", (n.next.Ip + n.rpcPort), args)
	if err != nil {
		n.logger.Error(err.Error())
		return
	}

	args = rpcArgs(hashKey, value)
	_, err = shared.SingleCall("Node.PutKey", (reply.Next.Ip + n.rpcPort), args)
	if err != nil {
		n.logger.Error(err.Error())
		return
	}
}
*/
