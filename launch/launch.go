package main

import (
//	"github.com/joonnna/ds_chord/nameserver"
	"os/exec"
	"log"
	"fmt"
	"io"
	"os"
	"syscall"
	"os/signal"
	"strings"
	"strconv"
	"time"
)

func killShit(pipes []io.WriteCloser) {
	fmt.Println("CLEANUP")
	for _, pipe := range pipes {
		pipe.Write([]byte("kill"))
		pipe.Close()
	}
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

func launch(nodeName string, path string, nameServer string, id int) io.WriteCloser {

	var command string
	if nameServer != "" {
		command = "go run " + path + " " + nameServer + "," + strconv.Itoa(id) + ",node"
	} else {
		command = "go run " + path + " ,nameserver"
	}
	cmd := exec.Command("ssh", "-T", nodeName, command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Start()

	return pipe
}


func main () {
	numHosts := 4
	nodeList := getNodeList(numHosts)

	//nameServerPath := "./go/src/github.com/joonnna/ds_chord/nameserver/nameserver.go"
	//nodePath := "./go/src/github.com/joonnna/ds_chord/node/node.go"
	//clientPath := "./go/src/github.com/joonnna/ds_chord/client/client.go"
	path := "./go/src/github.com/joonnna/ds_chord/main.go"

	var pipeSlice []io.WriteCloser

	nameServerIp := nodeList[0]

	pipe := launch(nameServerIp, path, "", 0)

	pipeSlice = append(pipeSlice, pipe)

	time.Sleep((3 * time.Second))


	for idx, ip := range nodeList {
		if idx != 0 {
			if idx == 3 {
				pipe = launch(ip, path, nameServerIp, 1)
			} else {
				pipe = launch(ip, path, nameServerIp, idx+2)
			}
			pipeSlice = append(pipeSlice, pipe)
		}
		time.Sleep((3 * time.Second))

	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	killShit(pipeSlice)
}
