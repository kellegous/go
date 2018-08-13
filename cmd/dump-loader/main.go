package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	apiPath = "https://%s:%s/api/url/%s"
)

type config struct {
	host     string
	port     string
	dumpFile string
}

type goData struct {
	url string
	ts  time.Time
}

func main() {
	c := config{}
	flag.StringVar(&c.host, "host", "localhost", "host to post data to")
	flag.StringVar(&c.port, "port", "8067", "port on host to talk to")
	flag.StringVar(&c.dumpFile, "file", "", "dump file to load from")
	flag.Parse()

	if c.dumpFile == "" {
		log.Fatal("dump file must be specified with -file argument")
	}

	var d interface{}

	f, err := ioutil.ReadFile(c.dumpFile)

	if err != nil {
		log.Printf("error reading dump file : %s\n", c.dumpFile)
		log.Fatal(err)
	}

	err = json.Unmarshal(f, &d)

	if err != nil {
		log.Printf("error parsing dump file : %s\n", c.dumpFile)
		log.Fatal(err)
	}

	links := d.(map[string]interface{})

	for k, v := range links {
		req := fmt.Sprintf(apiPath, c.host, c.port, k)
		p, err := json.Marshal(&v)
		if err != nil {
			log.Printf("error marshalling data for link : %s\n", k)
			log.Println(err)
			continue
		}
		resp, err := http.Post(req, "application/json", bytes.NewReader(p))
		if err != nil {
			log.Printf("error POSTing link : %s : %s\n", k, err)
		} else {
			log.Printf("POSTed short link (%s) : %s\n", resp.Status, k)
		}
	}
}
