package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
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
		fmt.Sprintf("--asset-proxy-url=%s", proxyURL))
	c.Dir = root
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Start()
}

func isAvailable(
	ctx context.Context,
	url string,
) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodHead,
		url,
		nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func waitForAvailability(
	ctx context.Context,
	urls ...string,
) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, url := range urls {
		g.Go(func() error {
			for {
				if err := isAvailable(ctx, url); err == nil {
					return nil
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Millisecond * 100):

				}
			}
		})
	}
	return g.Wait()
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

	viteURL := fmt.Sprintf("http://localhost:%d", flags.Vite.Port)
	if err := startGoServer(
		ctx,
		flags.Root,
		viteURL,
	); err != nil {
		log.Panic(err)
	}

	{
		goURL := "http://localhost:8067/"
		ctx, done := context.WithTimeout(ctx, time.Second*5)
		defer done()

		if err := waitForAvailability(
			ctx,
			viteURL,
			goURL,
		); err != nil {
			log.Panic(err)
		}

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Println()
		fmt.Printf("development server running...\n%s\n", green(goURL))
	}

	<-ctx.Done()
}
