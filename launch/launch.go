package main

import (
	"os/exec"
	"log"
	"fmt"
	"os"
)



func parseArgs() {


}

func getNodeList() {
	fmt.Println(os.Getenv("PATH"))
	nodeListCmd := "/test.sh"
	cmd := exec.Command(nodeListCmd)

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	result, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	//return result
}


func initClient() {


}

func initNode(ip string) {

}


func initNameserver() {


}


func main () {
	getNodeList()
	//fmt.Println(nodeList)

}
