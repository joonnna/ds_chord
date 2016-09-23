package main

import (
	"github.com/joonnna/ds_chord/node"
	"github.com/joonnna/ds_chord/nameserver"
	"os"
	"strings"
//	"fmt"
)



func main() {
	temp := strings.Join(os.Args[1:], "")
	args := strings.Split(temp, ",")

	progType := args[(len(args)-1)]

	switch progType {

		case "nameserver":
			nameserver.Run()

		case "node":
			node.Run(args[0], args[1])
	}

}
