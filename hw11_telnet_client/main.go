package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var timeout time.Duration

func main() {
	host, port := os.Args[0], os.Args[1]
	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	ctx, ctxCancelF := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	defer func(client TelnetClient) {
		if err := client.Close(); err != nil {
			return
		}
	}(client)

	send(client, ctxCancelF)
	receive(client, ctxCancelF)

	select {
	case <-sigCh:
	case <-ctx.Done():
		close(sigCh)
	}
}

func send(client TelnetClient, ctxCancelF context.CancelFunc) {
	err := client.Send()
	if err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			return
		}
	}
	ctxCancelF()
}

func receive(client TelnetClient, ctxCancelF context.CancelFunc) {
	err := client.Receive()
	if err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			return
		}
	}
	ctxCancelF()
}
