package util

import (
	"io/ioutil"
	"os"
	"fmt"
	"log"
	"github.com/gorilla/mux"
	"net/http"
)


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
