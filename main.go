package main

import (
	"github.com/joonnna/ds_chord/storage"
	"github.com/joonnna/ds_chord/nameserver"
	"github.com/joonnna/ds_chord/client"
	"os"
	"strings"
)



func main() {
	temp := strings.Join(os.Args[1:], "")
	args := strings.Split(temp, ",")

	progType := args[(len(args)-1)]

	switch progType {

		case "nameserver":
			nameserver.Run(args[0])

		case "node":
			storage.Run(args[0], args[1], args[2])

		case "client":
			client.Run(args[0], args[1])
	}

}
