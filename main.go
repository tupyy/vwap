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

	logger.Infof("Git commit: %s", CommitID)
	logger.Infof("Conf used: %+v", config)

	// setup output
	var out *output.Writer
	if len(config.OutputFile) == 0 {
		out = output.NewStdOutputWriter()
	} else {
		// try to create the output file
		outputFile, err := os.OpenFile(config.OutputFile, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}

		out = output.NewFileWriter(outputFile)
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
		logger.Errorf("error connecting to ws: %v", err)
		os.Exit(1)
	}
	defer func() {
		err := wsClient.Shutdown()
		if err != nil {
			logger.Errorf("error disconnecting: %v", err)
		}
	}()

	// subscribe
	subscribeCtx, subscribeCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer subscribeCancel()
	if err := wsClient.Subscribe(subscribeCtx); err != nil {
		logger.Errorf("error subscribing: %v", err)
		os.Exit(1)
	}

	// start reading
	errCh := make(chan error)
	wsClient.Receive(ctx, msgCh, errCh)
	go func() {
		for e := range errCh {
			logger.Errorf("error reading ws: %+v", e)
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	// handle int & term signals
	sigCh := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		done <- true
	}()

	<-done

	logger.Infof("shutting down")

	// close the client
	err := wsClient.Shutdown()
	if err != nil {
		logger.Errorf("ws client closed with error: %+v", err)
	}

	logger.Infof("ws connection closed")

	// shutdown usecase
	avgManager.Shutdown()

	// cancel context
	cancel()
}
