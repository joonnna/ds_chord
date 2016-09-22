package main

import (
//	"github.com/joonnna/ds_chord/nameserver"
	"os/exec"
	"log"
	"fmt"
	"os"
	"strings"
	"strconv"
	"time"
)

func getNodeList(numHosts int) []string {
	scriptName := "./rocks_list_random_hosts.sh"
	cmd := exec.Command("sh", scriptName, strconv.Itoa(numHosts))

	result, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	nodeList := strings.Split(string(result), " ")

	return nodeList[:numHosts]
}

func sshToNode(ip string) {
	cmd := exec.Command("ssh", ip)

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func launch(nodeName string, path string, nameServer string) {
	fmt.Println(nodeName)
	var command string
	if nameServer != "" {
		command = "go run " + path + " " + nameServer
	} else {
		command = "go run " + path
	}
	cmd := exec.Command("ssh", "-T", nodeName, command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
}


func main () {
	numHosts := 1
	nodeList := getNodeList(numHosts)

	fmt.Println(nodeList)

	//nameServerPath := "./go/src/github.com/joonnna/ds_chord/nameserver/nameserver.go"
	nodePath := "./go/src/github.com/joonnna/ds_chord/node/node.go"
	//clientPath := "./go/src/github.com/joonnna/ds_chord/client/client.go"

	nameServerIp := nodeList[0]
	//launch(nameServerIp, nameServerPath, "")
	time.Sleep(10)
	for idx, ip := range nodeList {
		if idx == 0 {
			launch(ip, nodePath, nameServerIp)
		}
//		if idx == 2{
//			launch(ip, clientPath, nameServerIp)
//		}
		time.Sleep(1)
	}
}
