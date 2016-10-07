package client

import (
	"time"
	//"runtime"
	"fmt"
	"sync"
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
	client *http.Client
}

type result struct {
	errors int
	meanLatency float64
}

const (
	numWorkers = 100
)

func genValue() string {
	randValue := strconv.Itoa(rand.Int())

	return randValue
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
	r.Body.Close()
}

func (c *Client) getValue(ip string, key string) string {
	req, err := http.NewRequest("GET", "http://" + ip + c.port  + "/" + key, nil)
	if err != nil {
		c.log.Error(err.Error())
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error(err.Error())
		return ""
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.log.Error(err.Error())
		}
		resp.Body.Close()
		return string(body)
	}
}


func (c *Client) putValue(ip string, key string, value string) {
	req, err := http.NewRequest("PUT", "http://" + ip + c.port + "/" + key, strings.NewReader(value))
	if err != nil {
		c.log.Error(err.Error())
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error(err.Error())
	} else {
		resp.Body.Close()
	}
}


func (c *Client) assertKeys(start int, wg *sync.WaitGroup, ch chan result)  {
	defer wg.Done()
	errors := 0
	numKeys := 100
	var totalLatency float64
	totalLatency = 0.00

	for i := 0; i < numKeys; i++ {
		key := strconv.Itoa(start + i)
		value := genValue()

		idx := rand.Int() % len(c.nodeIps)

		start := time.Now()
		c.putValue(c.nodeIps[idx], key, value)
		totalLatency += float64((time.Since(start)*time.Second))


		start = time.Now()
		getVal := c.getValue(c.nodeIps[0], key)
		totalLatency += float64((time.Since(start)*time.Second))

		c.log.Debug(getVal)
		tmp := strings.Split(getVal, "\"")
		if len(tmp) <= 1 {
			errors += 1
			continue
		}
		getVal = tmp[1]
		//c.log.Debug("NEW VAL : " + tmp)

		if !(value == getVal) {
			errors += 1
		}
	}

	mean := (totalLatency/float64(numKeys*2))
	retVal := result {
		errors: errors,
		meanLatency: mean }
	ch <- retVal
}

func Run(nameServer string, port string) {
	go util.CheckInterrupt()

	numErrors := 0
	var meanLatency float64
	meanLatency = 0.00

	var wg sync.WaitGroup

	ch := make(chan result, numWorkers)

	l := new(logger.Logger)
	l.Init((os.Stdout), "Client", 0)

	c := &Client{
		nameServer: "http://" + nameServer + port,
		port: port,
		log: l,
		client: &http.Client{} }

	c.getNodeList()

	c.log.Testing("STARTED TESTING BOOOOYS")
	for i := 0; i < numWorkers; i++ {
		go c.assertKeys((i*numWorkers), &wg, ch)
		wg.Add(1)
	}


	for i := 0; i < numWorkers; i++ {
		res := <-ch
		numErrors += res.errors
		meanLatency += res.meanLatency
	}


	meanLatency = (meanLatency/numWorkers)

	str := fmt.Sprint(meanLatency)
	c.log.Testing("Number of PUT/GET errors " + strconv.Itoa(numErrors))
	c.log.Testing("Mean Latency(s) : " + str)
}

