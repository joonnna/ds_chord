package node

import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"time"
	"strings"
	//"math/big"
//	"github.com/joonnna/ds_chord/logger"
//	"github.com/joonnna/ds_chord/node_communication"
)

/*
func (n *Node) assertSuccessor(newSucc *big.Int) {
	cmp := n.Next.Id.Cmp(newSucc)
	if cmp == 0 {
		n.logger.Error("Invalid successor")
	}

	c := n.Next.Id.Cmp(&n.id)
	if c == 1 && cmp == -1 {
		n.logger.Error("Invalid successor")
	}
}

func (n *Node) assertPreDecessor(newPre *big.Int) {
	cmp := n.prev.Id.Cmp(newPre)
	if cmp == 0 {
		n.logger.Error("Invalid predecessor")
	}

	c := n.prev.Id.Cmp(n.id)
	if c == -1 && cmp == 1 {
		n.logger.Error("Invalid predecessor")
	}
}
*/
func (n *Node) putIp() {
	req, err := http.NewRequest("PUT", n.NameServer+"/", strings.NewReader(n.ip))
	if err != nil {
		n.logger.Error(err.Error())
	}

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{Timeout : timeout}
	resp, err := client.Do(req)
	if err != nil {
		n.logger.Error(err.Error())
	} else {
		resp.Body.Close()
	}

}

func GetNodeList(nameServer string) ([]string, error)  {
	var nodeIps []string

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{Timeout : timeout}

	r, err := client.Get(nameServer)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &nodeIps)
	if err != nil {
		return nil, err
	}
	return nodeIps, nil
}
