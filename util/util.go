package util

import (
	"io/ioutil"
	"os"
	"fmt"
	"log"
	"github.com/gorilla/mux"
	"net/http"
)

func InKeySpace(id int, nodeId int, prevId int) bool {
	if nodeId < prevId {
		if (id > nodeId && id > prevId) || (id < nodeId && id < prevId){
			return true
		} else {
			return false
		}
	} else {
		if id < nodeId && id > prevId {
			return true
		} else {
			return false
		}
	}
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
		m, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("RECEIVED KILL SIGNAL")
		if string(m) == "kill" {
			fmt.Println("EXITING")
			os.Exit(1)
		}
	}
}
