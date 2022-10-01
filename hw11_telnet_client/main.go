package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var timeout time.Duration

func main() {
	if flag.NArg() < 2 {
		log.Fatal("not enough arguments")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	telnet := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := telnet.Connect(); err != nil {
		log.Fatal(err)
	}
	defer telnet.Close()

	ctx, ctxCancelF := context.WithCancel(context.Background())
	go func() {
		defer ctxCancelF()

		err := telnet.Send()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	go func() {
		defer ctxCancelF()

		err := telnet.Receive()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
	case <-ctx.Done():
		close(sigCh)
	}
}
