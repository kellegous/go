package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"github.com/spf13/pflag"
)

func startViteServer(
	ctx context.Context,
	root string,
	port int,
) error {
	c := exec.CommandContext(
		ctx,
		"node_modules/.bin/vite",
		"--clearScreen=false",
		fmt.Sprintf("--port=%d", port))
	c.Dir = root
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Start()
}

func startGoServer(
	ctx context.Context,
	root string,
	proxyURL string,
) error {
	c := exec.CommandContext(
		ctx,
		"bin/go",
		"--host=http://go",
		fmt.Sprintf("--asset-proxy-url=%s", proxyURL))
	c.Dir = root
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Start()
}

func main() {
	var flags struct {
		Root string
		Vite struct {
			Port int
		}
	}

	pflag.StringVar(
		&flags.Root,
		"root",
		".",
		"Root directory for the server")
	pflag.IntVar(
		&flags.Vite.Port,
		"vite.port",
		3000,
		"Port for the vite server")

	pflag.Parse()

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	if err := startViteServer(
		ctx,
		flags.Root,
		flags.Vite.Port,
	); err != nil {
		log.Panic(err)
	}

	if err := startGoServer(
		ctx,
		flags.Root,
		fmt.Sprintf("http://localhost:%d", flags.Vite.Port),
	); err != nil {
		log.Panic(err)
	}

	<-ctx.Done()
}
