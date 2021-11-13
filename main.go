package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tupyy/vwap/internal/compute"
	"github.com/tupyy/vwap/internal/conf"
	"github.com/tupyy/vwap/internal/log"
	"github.com/tupyy/vwap/internal/manager"
	"github.com/tupyy/vwap/internal/repo/output"
	"github.com/tupyy/vwap/internal/repo/ws"
)

// CommitID contains the SHA1 Git commit of the build.
// It's evaluated during compilation.
var CommitID string

func main() {
	config := conf.Get()

	logger := log.GetLogger()

	logger.Info("Git commit: %s", CommitID)
	logger.Info("Conf used: %+v", config)

	// setup output
	var out *output.OutputWriter
	if len(config.OutputFile) == 0 {
		out = output.NewStdOutputWriter()
	} else {
		// try to create the output file
		outputFile, err := os.OpenFile(config.OutputFile, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}

		out = output.NewFileOutputWriter(outputFile)
	}

	// create message channel
	msgCh := make(chan interface{})

	// setup calculators
	avgManager := manager.NewAvgManager(out)
	for _, p := range config.TradingPairs {
		c := compute.NewAvgCalculator(int(config.MaxDataPoints))
		avgManager.AddAvgCalculator(p, c)
	}

	// define our context
	ctx, cancel := context.WithCancel(context.Background())

	avgManager.Start(ctx, msgCh)

	// dial the connection
	wsClient := ws.NewClient(config.TradingPairs)

	connectCtx, connectCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer connectCancel()

	// connect to cointbase
	if err := wsClient.Connect(connectCtx, config.Endpoint); err != nil {
		logger.Error("error connecting to ws: %v", err)
		os.Exit(1)
	}
	defer func() {
		err := wsClient.Disconnect()
		if err != nil {
			logger.Error("error disconnecting: %v", err)
		}
	}()

	// subscribe
	subscribeCtx, subscribeCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer subscribeCancel()
	if err := wsClient.Subscribe(subscribeCtx); err != nil {
		logger.Error("error subscribing: %v", err)
		os.Exit(1)
	}

	// start reading
	errCh := make(chan error)
	wsClient.Receive(ctx, msgCh, errCh)
	go func() {
		for e := range errCh {
			logger.Error("error reading ws: %+v", e)
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh

	logger.Info("shutting down")

	// shutdown usecase
	avgManager.Shutdown()

	// cancel context
	cancel()
}
