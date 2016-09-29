package client

import (
	//"time"
	//"runtime"
	//"sync"
	"os"
	"net/http"
	"encoding/json"
	"strings"
	"io/ioutil"
	"math/rand"
	"strconv"
	"github.com/joonnna/ds_chord/logger"
	"github.com/joonnna/ds_chord/util"
)


type Client struct {
	nameServer string
	nodeIps []string
	log *logger.Logger
	port string
}


func (c *Client) genKeyValue() (string, string){
	randKey := strconv.Itoa(rand.Int())
	randValue := strconv.Itoa(rand.Int())

	return randKey, randValue
}

func (c *Client) getNodeList()  {
	r, err := http.Get(c.nameServer)
	if err != nil {
		c.log.Error(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.log.Error(err.Error())
	}

	err = json.Unmarshal(body, &c.nodeIps)
	if err != nil {
		c.log.Error(err.Error())
	}
}

func (c *Client) getValue(ip string, key string) string {
	req, err := http.NewRequest("GET", "http://" + ip + c.port  + "/" + key, nil)
	if err != nil {
		c.log.Error(err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		c.log.Error(err.Error())
		return ""
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.log.Error(err.Error())
		}
		//fmt.Println(string(body))
		return string(body)
	}
}


func (c *Client) putValue(ip string, key string, value string) {
	req, err := http.NewRequest("PUT", "http://" + ip + c.port + "/" + key, strings.NewReader(value))
	if err != nil {
		c.log.Error(err.Error())
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		c.log.Error(err.Error())
	}
}


func (c *Client) assertKeys() int {
	errors := 0
	numKeys := 1000

	for i := 0; i < numKeys; i++ {

		c.getNodeList()

		key, value := c.genKeyValue()

		c.putValue(c.nodeIps[0], key, value)

		getVal := c.getValue(c.nodeIps[0], key)

		tmp := strings.Split(getVal, "\"")
		if len(tmp) < 1 {
			c.log.Error("Failed to PUT/GET key " + key)
			errors += 1
			continue
		}
		getVal = tmp[1]
		c.log.Debug(value)
		//c.log.Debug("NEW VAL : " + tmp)

		if !(value == getVal) {
			c.log.Error("Failed to PUT/GET key " + key)
			errors += 1
		}
	}
	return errors
}

func Run(nameServer string, port string) {
	go util.CheckInterrupt()

	c := new(Client)
	c.nameServer = "http://" + nameServer + port
	c.port = port
	l := new(logger.Logger)
	l.Init((os.Stdout), "Client", 0)
	c.log = l
	c.log.Debug(port)
	c.log.Info("Started Client")
	c.log.Debug(c.nameServer)
	numErrors := c.assertKeys()
	c.log.Info("Number of PUT/GET errors " + strconv.Itoa(numErrors))
}

