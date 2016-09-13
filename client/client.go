package client

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"strings"
	"io/ioutil"
	"math/rand"
	"strconv"

)


type Client struct {
	NameServer string
	NodeIps []string
}


func (c *Client) GenKeyValue() (string, string){
	randKey := strconv.Itoa(rand.Int())
	randValue := strconv.Itoa(rand.Int())

	return randKey, randValue
}

func (c *Client) GetNodeList()  {
	fmt.Println("sending GET request..")

	r, err := http.Get(c.NameServer)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(r.Body)

	err = json.Unmarshal(body, &c.NodeIps)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) GetValue(ip string, key string) {
	req, err := http.NewRequest("GET", ip+"/"+ key, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else {
		body,_ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))

		resp.Body.Close()
	}
}


func (c *Client) PutValue(ip string, key string, value string) {
	
	req, err := http.NewRequest("PUT", ip+"/"+ key, strings.NewReader(value))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("sending put request..")
	
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else {
		resp.Body.Close()
	}


}
