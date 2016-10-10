package main

import (
	"os/exec"
	"log"
	"os"
	"syscall"
	"os/signal"
	"strings"
	"strconv"
	"time"
	"fmt"
)
/* Ports to use */
var (
	http = 2345
	rpc = 7453
)
func cleanUp() {
	cmd := exec.Command("sh", "/share/apps/bin/cleanup.sh")
	cmd.Run()
}

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
/* Launches the given application(client, nameserver or node)
   nodeName: address of the node to launch
   path: path to the executable to launch
   nameserver: address of the nameserver, empty if application is the nameserver
   flag: -1 if launching the client
   */
func launch(nodeName string, path string, nameServer string, flag int)  {
	var command string

	httpPort := ":" + strconv.Itoa(http)
	rpcPort := ":" + strconv.Itoa(rpc)

	if flag == -1 {
		command = "go run " + path + " " + nameServer + "," + httpPort + ",client"
	} else if nameServer != "" {
		command = "go run " + path + " " + nameServer + "," + httpPort + "," + rpcPort + ",node"
	} else {
		command = "go run " + path + " " + httpPort + " ,nameserver"
	}
	cmd := exec.Command("ssh", "-T", nodeName, command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Start()
}


func main () {
	numHosts := 18
	nodeList := getNodeList(numHosts)

	path := "./go/src/github.com/joonnna/ds_chord/main.go"

	nameServerIp := nodeList[0]
	fmt.Println(nameServerIp)
	launch(nameServerIp, path, "", 0)

	time.Sleep(3 * time.Second)

	for idx, ip := range nodeList {
		if idx == len(nodeList) - 1  {
			//time.Sleep((20 * time.Second))
			//launch(ip, path, nameServerIp, -1)
		} else if idx != 0 {
			launch(ip, path, nameServerIp, idx)
		}
	}

	/* Wait for CTRL-C then shut all nodes down*/
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	cleanUp()
}
