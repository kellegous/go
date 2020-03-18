package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	apiPath  = "%s://%s:%s/api/url/%s"
	urlsPath = "%s://%s:%s/api/urls/?include-generated-names=true"
)

type config struct {
	proto    string
	host     string
	port     string
	dumpFile string
	dump     bool
}

type routesDump struct {
	OK     bool    `json:"ok"`
	Routes []route `json:"routes"`
	Next   string  `json:"next"`
}

type route struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Hits string `json:"hits"`
}

func init() {
	// Initialize logging
	formatter := &log.JSONFormatter{
		TimestampFormat: "2006-02-01 15:04:05",
	}
	log.SetFormatter(formatter)
	log.SetLevel(log.DebugLevel)
}

func main() {
	c := config{}
	pflag.StringVar(&c.proto, "protocol", "http", "protocol to use (HTTP or HTTPS)")
	pflag.StringVar(&c.host, "host", "localhost", "host to post data to")
	pflag.StringVar(&c.port, "port", "8067", "port on host to talk to")
	pflag.StringVar(&c.dumpFile, "file", "", "dump file to load from (or save to if --dump given)")
	pflag.BoolVar(&c.dump, "dump", false, "dump links from the api")
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

	if c.dump {
		req := fmt.Sprintf(urlsPath, c.proto, c.host, c.port)
		resp, err := http.Get(req)
		if err != nil {
			log.Printf("error making get request: %s", err.Error())
			return
		}
		body, readErr := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if readErr != nil {
			log.Fatal(readErr)
		}
		err = ioutil.WriteFile(c.dumpFile, body, 0644)
		if err != nil {
			log.Println("error dumping JSON from API")
			return
		}

		return
	}

	var d routesDump

	f, err := ioutil.ReadFile(c.dumpFile)

	if err != nil {
		log.Printf("error reading dump file: %s\n", c.dumpFile)
		log.Fatal(err)
	}

	log.Println("Reading file...")
	err = json.Unmarshal(f, &d)

	if err != nil {
		log.Printf("error parsing dump file: %s\n", c.dumpFile)
		log.Fatal(err)
	}
	log.Printf("Parsed dump: %+v", d)

	for _, v := range d.Routes {
		log.Println("In for loop...")
		req := fmt.Sprintf(apiPath, c.proto, c.host, c.port, v.Name)
		log.Printf("Request: %+v", req)
		p, err := json.Marshal(&v)
		if err != nil {
			log.Printf("error marshalling data for link: %s\n", v.Name)
			log.Println(err)
			continue
		}
		resp, err := http.Post(req, "application/json", bytes.NewReader(p))
		if err != nil {
			log.Printf("error POSTing link: %s %s\n", v.Name, err.Error())
		} else {
			log.Printf("POSTed short link (%s): %s\n", resp.Status, v.Name)
			resp.Body.Close()
		}
	}
}
