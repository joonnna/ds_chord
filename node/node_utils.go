package node

import (
	"net/http"
	"time"
	"strings"
)

/* Sends a put request to the nameserver containing the ip of the node*/
func (n *Node) putIp() {
	req, err := http.NewRequest("PUT", n.NameServer+"/", strings.NewReader(n.Ip + n.httpPort))
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

