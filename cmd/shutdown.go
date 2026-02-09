package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func shutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		println("\ninterrupting...")
		cancel()
	}()
}
