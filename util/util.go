package util

import (
	"math/big"
//	"strings"
	"fmt"
	"os"
	"io"
	"crypto/sha1"
	"io/ioutil"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/joonnna/ds_chord/node_communication"
)
/* With  Upper Inlcude */
func InKeySpace(start, end, newId big.Int) (bool, string) {
	var ret bool

	startEndCmp := start.Cmp(&end)

	if startEndCmp == -1 {
		if start.Cmp(&newId) == -1 && end.Cmp(&newId) >= 0 {
			ret = true
		} else {
			ret = false
		}
	} else {
		if start.Cmp(&newId) == -1 || end.Cmp(&newId) >=0 {
			ret = true
		} else {
			ret = false
		}
	}
	if !ret {
		str := fmt.Sprintf("Start : %s, End : %s, newId : %s\n", start.String(), end.String(), newId.String())
		return ret, str
	} else {
		return ret, ""
	}
}

func BetweenNodes(start, end, newId big.Int) (bool, string) {
	var ret bool

	startEndCmp := start.Cmp(&end)

	if startEndCmp == -1 {
		if start.Cmp(&newId) == -1 && end.Cmp(&newId) == 1 {
			ret = true
		} else {
			ret = false
		}
	} else {
		if start.Cmp(&newId) == -1 || end.Cmp(&newId) == 1 {
			ret = true
		} else {
			ret = false
		}
	}
	if !ret {
		str := fmt.Sprintf("Start : %s, End : %s, newId : %s\n", start.String(), end.String(), newId.String())
		return ret, str
	} else {
		return ret, ""
	}
}



/*
func InKeySpace(currId, newId, prevId big.Int) bool {
	cmp := currId.Cmp(&newId)
	if cmp == 0 {
		return true
	}

	prevCmp := currId.Cmp(&prevId)

	idPrevCmp := prevId.Cmp(&newId)

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
*/
func NoInclude(currId, newId, prevId big.Int) bool {
	cmp := currId.Cmp(&newId)

	prevCmp := currId.Cmp(&prevId)

	idPrevCmp := prevId.Cmp(&newId)

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

func ConvertKey(key string) big.Int {
	h := sha1.New()
	io.WriteString(h, key)

	ret := new(big.Int)

	ret.SetBytes(h.Sum(nil))
	return *ret
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

		if string(msg) == "kill" {
			os.Exit(1)
		}
	}
}

func RpcArgs(key big.Int, value string) shared.Args {
	args := shared.Args {
		Key: key,
		Value: value }

	return args
}


func CreateArgs(nodeAddr string, nodeId big.Int) shared.Args {
	n := shared.NodeInfo{
		Ip: nodeAddr,
		Id: nodeId }

	args := shared.Args{
		Node: n }

	return args
}

func UpdateArgs(id big.Int, prevId string) shared.UpdateArgs {
	args := shared.UpdateArgs {
		Id: id,
		PrevId: prevId }

	return args
}
