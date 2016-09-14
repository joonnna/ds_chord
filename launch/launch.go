package launch

import (
	"os/exec"
	"log"
)



func parseArgs() {


}

func getNodeList() {
	nodeListCmd = "sh rocks_list_random_hosts.sh"
	cmd := exec.command(nodeListCmd)

	output, err := cmd.Run().Output()
	if err != nil {
		log.Fatal(err)
	}
}


func initClient() {


}

func initNode(ip string) {

}


func initNameserver() {


}


func main () {
	nodeList := getNodeList()
	fmt.Println(nodeList)

}
