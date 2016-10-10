package client

import (
	"time"
	"runtime"
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
	"log"
)


type Client struct {
	nameServer string
	nodeIps []string
	keys []string
	log *logger.Logger
	port string
	client *http.Client
}

type result struct {
	errors int
	meanLatency float64
}

var (
	numReq = 10000
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
	req, err := http.NewRequest("GET", "http://" + ip + "/" + key, nil)
	if err != nil {
		c.log.Error(err.Error())
	}
	req.Close = true
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
	req, err := http.NewRequest("PUT", "http://" + ip + "/" + key, strings.NewReader(value))
	if err != nil {
		c.log.Error(err.Error())
	}
	req.Close = true
	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error(err.Error())
	} else {
		resp.Body.Close()
	}
}


func (c *Client) assertKeys(numKeys int, start int, wg *sync.WaitGroup, ch chan result)  {
	wg.Done()
	wg.Wait()
	var errors int = 0
	var totalLatency float64 = 0.00

	for i := 0; i < numKeys; i++ {
		key := c.keys[start*numKeys + i]
		value := genValue()

		idx := rand.Int() % len(c.nodeIps)

		start := time.Now()
		c.putValue(c.nodeIps[idx], key, value)
		totalLatency += (time.Since(start).Seconds())

		start = time.Now()
		getVal := c.getValue(c.nodeIps[0], key)
		totalLatency += float64(time.Since(start).Seconds())

		tmp := strings.Split(getVal, "\"")
		if len(tmp) < 1 {
			c.log.Error("Failed to PUT/GET value " + value)
			errors += 1
			continue
		}
		getVal = tmp[1]

		if !(value == getVal) {
			c.log.Error("Failed to PUT/GET key " + key)
			errors += 1
		}
	}
	mean := (totalLatency/float64(numKeys*2))
	retVal := result {
		meanLatency: mean,
		errors: errors }
	ch <- retVal
}


func logResult(d1 string, d2 string, d3 string, d4 string) {
	f, err := os.OpenFile("./result.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	str := fmt.Sprintf("%s\t%s\t%s\t%s\n", d1, d2, d3, d4)

	_, err = f.WriteString(str)
	if err != nil {
		log.Fatal(err)
	}

}

func Run(nameServer string, port string) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	numErrors := 0
	var meanLatency float64 = 0.00

	l := new(logger.Logger)
	l.Init((os.Stdout), "Client", 0)

	c := &Client{
		nameServer: "http://" + nameServer + port,
		port: port,
		log: l,
		client: &http.Client{} }

	totalKeys := numReq/2
	for i := 0; i < totalKeys; i++ {
		c.keys = append(c.keys, strconv.Itoa(rand.Int()))
	}

	c.log.Testing("STARTED TESTING")

	var wg sync.WaitGroup
	numWorkers := 300
	numKeys := totalKeys/numWorkers

	ch := make(chan result, numWorkers)
	info := fmt.Sprintf("Workers : %s", strconv.Itoa(numWorkers))
	c.log.Testing(info)

	c.getNodeList()
	for k := 0; k < numWorkers; k++ {
		wg.Add(1)
		go c.assertKeys(numKeys, k, &wg, ch)
	}
	wg.Wait()

	start := time.Now()
	for i := 0; i < numWorkers; i++ {
		res := <-ch
		numErrors += res.errors
		meanLatency += res.meanLatency
	}
	end := time.Since(start)

	reqPerSec := (float64(numReq))/(float64(end.Seconds()))

	meanLatency = (meanLatency/float64(numWorkers))

	str := fmt.Sprint(meanLatency)
	tmp := fmt.Sprint(reqPerSec)
	total := fmt.Sprint(end.Seconds())

	c.log.Testing("Total Time : " + total)
	c.log.Testing("Mean Latency(s) : " + str)
	c.log.Testing("Request per second(s) : " + tmp)

	c.log.Testing("ERRORS : " + strconv.Itoa(numErrors))
	//logResult(tmp, str, strconv.Itoa(numWorkers), "16")
}

