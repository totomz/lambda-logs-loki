package main

import (
	"context"
	"golang.org/x/exp/slog"
	"os"
	"os/signal"
	"path"
	"syscall"
)

var stdout = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: true,
}))

func main() {
	extensionName := path.Base(os.Args[0])
	logger := stdout.With("extname", extensionName, "runtime", os.Getenv("AWS_LAMBDA_RUNTIME_API"))

	extensionClient := NewExtensionClient(os.Getenv("AWS_LAMBDA_RUNTIME_API"))
	ctx, cancel := context.WithCancel(context.Background())

	// handle termination signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sigs
		cancel()
		logger.Info("received exit signal", "signal", s)
	}()

	// Register extension as soon as possible
	_, err := extensionClient.Register(ctx, extensionName)
	if err != nil {
		panic(err)
	}

	// TODO Greppa i log

	// Block until a quit signal
	// or when the lambda runtime tell us to close
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// This is a blocking call
			res, err := extensionClient.NextEvent(ctx)
			if err != nil {
				logger.Info("lambda extension event error, exiting", "error", err)
				return
			}
			logger.Info("got a lambda event", "event", res.EventType)

			// TODO Flush log queue in here after waking up

			// Exit if we receive a SHUTDOWN event
			if res.EventType == Shutdown {
				logger.Info("Received SHUTDOWN event")
				// TODO Close and flush everything
				return
			}
		}
	}
}
