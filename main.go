package main

import (
	"flag"
	"log"

	"github.com/kellegous/go/context"
	"github.com/kellegous/go/web"
)

func main() {
	flagData := flag.String("data", "data", "data")
	flagAddr := flag.String("addr", ":8067", "addr")
	flag.Parse()

	ctx, err := context.Open(*flagData)
	if err != nil {
		log.Panic(err)
	}

	log.Panic(web.ListenAndServe(*flagAddr, ctx))
}
