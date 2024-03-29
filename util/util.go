package util

import (
	"strings"
	"os"
	"io"
	"fmt"
	"crypto/sha1"
	"io/ioutil"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/joonnna/ds_chord/node_communication"
)

func InKeySpace(currId, newId, prevId string) bool {
	cmp := strings.Compare(currId, newId)
	if cmp == 0 {
		return true
	}
	prevCmp := strings.Compare(currId, prevId)

	idPrevCmp := strings.Compare(prevId, newId)

	if prevCmp == -1 {
		if (cmp == -1 && idPrevCmp == -1) || (cmp == 1 && idPrevCmp == 1) {
			return true
		} else {
			return false
		}
	} else {
		if cmp == 1 && idPrevCmp == -1 {
			return true
		} else {
			return false
		}
	}
}
func HashKey(key string) string {
	h := sha1.New()
	io.WriteString(h, key)
	return string(h.Sum(nil))
}

func GetNode(list []string, curNode string) string {
	for _, ip := range list {
		if ip != curNode {
			return ip
		}
	}
	return ""
}


func GetKey(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["key"]
}

func CheckInterrupt() {
	for {
		msg, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("YOOYOYOYOYOYO")
		if string(msg) == "kill" {
			os.Exit(1)
		}
	}
}

func RpcArgs(key string, value string) shared.Args {
	args := shared.Args {
		Key: key,
		Value: value }

	return args
}


func CreateArgs(nodeAddr string, nodeId string) shared.Args {
	n := shared.NodeInfo{
		Ip: nodeAddr,
		Id: nodeId }

	args := shared.Args{
		Node: n }

	return args
}

func UpdateArgs(id, prevId string) shared.UpdateArgs {
	args := shared.UpdateArgs {
		Id: id,
		PrevId: prevId }

	return args
}
