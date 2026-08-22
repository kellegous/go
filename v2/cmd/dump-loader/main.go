package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	apiPath = "%s://%s:%s/api/url/%s"
)

type config struct {
	proto    string
	host     string
	port     string
	dumpFile string
}

type RoutesDump struct {
	OK     bool    `json:"ok"`
	Routes []Route `json:"routes"`
	Next   string  `json:"next"`
}

type Route struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	c := config{}
	pflag.StringVar(&c.proto, "protocol", "https", "protocol to use. Only HTTP or HTTPS supported")
	pflag.StringVar(&c.host, "host", "localhost", "host to post data to")
	pflag.StringVar(&c.port, "port", "8067", "port on host to talk to")
	pflag.StringVar(&c.dumpFile, "file", "", "dump file to load from")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Panic(err)
	}

	// allow env vars to set pflags
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if c.dumpFile == "" {
		log.Fatal("dump file must be specified with --file argument")
	}

	var d RoutesDump

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

	for _, v := range d.Routes {
		req := fmt.Sprintf(apiPath, c.proto, c.host, c.port, v.Name)
		p, err := json.Marshal(&v)
		if err != nil {
			log.Printf("error marshalling data for link : %s\n", v.Name)
			log.Println(err)
			continue
		}
		resp, err := http.Post(req, "application/json", bytes.NewReader(p))
		if err != nil {
			log.Printf("error POSTing link : %s %s\n", v.Name, err.Error())
		} else {
			log.Printf("POSTed short link (%s) : %s\n", resp.Status, v.Name)
			resp.Body.Close()
		}
	}
}
