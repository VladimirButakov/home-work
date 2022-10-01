package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

var timeout time.Duration

func main() {
	host, port := os.Args[0], os.Args[1]
	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	defer func(client TelnetClient) {
		if err := client.Close(); err != nil {
			return
		}
	}(client)

	ctx, ctxCancelF := context.WithCancel(context.Background())

	go func() {
		defer ctxCancelF()

		err := client.Send()
		if err != nil {
			if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
				return
			}
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
