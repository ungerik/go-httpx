package httpx

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// GracefulShutdownServerOnSignal sets up a goroutine that listens for OS signals
// and gracefully shuts down the HTTP server when a signal is received.
//
// The function configures the server to perform a graceful shutdown, which means:
//   - The server stops accepting new connections
//   - All active connections are allowed to complete
//   - The server waits for all handlers to finish (up to the timeout)
//   - Resources are cleaned up properly
//
// Parameters:
//   - server: The http.Server to shut down
//   - signalLog: Optional logger for received signals (can be nil)
//   - errorLog: Optional logger for shutdown errors (can be nil)
//   - timeout: Maximum duration to wait for active connections to complete.
//     A value of zero means no timeout (wait indefinitely).
//   - signals: OS signals to listen for. If empty, defaults to SIGHUP, SIGINT, SIGTERM.
//
// Example:
//
//	server := &http.Server{Addr: ":8080", Handler: mux}
//	logger := log.New(os.Stdout, "", log.LstdFlags)
//	httpx.GracefulShutdownServerOnSignal(server, logger, logger, 30*time.Second)
//	if err := server.ListenAndServe(); err != http.ErrServerClosed {
//	    log.Fatal(err)
//	}
//
// Note: This function must be called before server.ListenAndServe() to ensure
// the signal handler is registered before the server starts.
func GracefulShutdownServerOnSignal(server *http.Server, signalLog, errorLog Logger, timeout time.Duration, signals ...os.Signal) {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM}
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, signals...)
	go func() {
		sig := <-shutdown
		if signalLog != nil {
			signalLog.Printf("Received signal: %s", sig)
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
