package main

import (
	"github.com/joonnna/ds_chord/nameserver"
	//"github.com/joonnna/ds_chord/client"
	"github.com/joonnna/ds_chord/node"
	//"reflect"
	//"fmt"
	//"io/ioutil"
)





func main () {

	go nameserver.HttpServer()

	nodeIp := "http://127.0.0.1:8080"

	n := &node.Node{Ip : nodeIp, NameServer: nodeIp}
	//fmt.Println(n.Ip)
	n.PutIp()
	n.PutIp()
	n.PutIp()
}
