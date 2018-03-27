package main

import (
	"flag"
	"github.com/HALtheWise/o-links/context"
	"github.com/HALtheWise/o-links/web"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var version string

func getVersion() string {
	if version == "" {
		return "none"
	}
	return version
}

func main() {
	flagData := flag.String("data", "data",
		"The location to use for the data store")
	flagAddr := flag.String("addr", ":"+os.Getenv("PORT"), //I hope this works, used to be "8067" - I made a similar change in cmd\dump-loader
		"The address that the HTTP server will bind")
	flagAdmin := flag.Bool("admin", false,
		"If allowing admin level requests")
	flag.Parse()

	ctx, err := context.Open(*flagData) //leveldb.Open() is called on flagData to create a db in that directory
	if err != nil {
		log.Panic(err)
	}
	defer ctx.Close()

	log.Printf("Serving on port %s", *flagAddr)
	log.Panic(web.ListenAndServe(*flagAddr, *flagAdmin, getVersion(), ctx))
}
