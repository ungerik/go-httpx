package httpx

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutdownServerOnSignal(server *http.Server, signalLog, errorLog Logger, timeout time.Duration, signals ...os.Signal) {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM}
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, signals...)
	go func() {
		sig := <-shutdown
		if signalLog != nil {
			signalLog.Printf("received signal %s", sig)
		}

		ctx := context.Background()
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		err := server.Shutdown(ctx)
		if err != nil && errorLog != nil {
			errorLog.Printf("http.Server shutdown error: %s", err)
		}
	}()
}
