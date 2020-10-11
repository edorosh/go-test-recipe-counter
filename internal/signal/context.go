package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// ContextWithSIGTERM creates a Context with cancellation when SIGINT or SIGTERM are sent
func ContextWithSIGTERM(ctx context.Context) context.Context {
	newCtx, cancel := context.WithCancel(ctx)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-signals:
			// outputs to stderr
			println("[signal] SIGTERM caught. Cancelling the context")
			cancel()
			close(signals)
		}
	}()

	return newCtx
}
