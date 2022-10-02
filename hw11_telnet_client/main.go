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
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout")
	flag.Parse()
	host, port := flag.Arg(0), flag.Arg(1)

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)

	ctx, ctxCancelF := context.WithCancel(context.Background())

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	go send(client, ctxCancelF)
	go receive(client, ctxCancelF)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
	case <-ctx.Done():
		close(sigCh)
	}
}

func send(client TelnetClient, ctxCancelF context.CancelFunc) {
	defer ctxCancelF()
	err := client.Send()
	if err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			return
		}
	}
}

func receive(client TelnetClient, ctxCancelF context.CancelFunc) {
	defer ctxCancelF()
	err := client.Receive()
	if err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			return
		}
	}
}
