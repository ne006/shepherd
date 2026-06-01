package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ne006/shepherd/cli"
	"github.com/ne006/shepherd/supervisor"
	"github.com/ne006/shepherd/utils"
	"go.uber.org/zap"
)

func main() {
	baseCtx := context.Background()
	cancelCtx, cancel := context.WithCancel(baseCtx)

	if err := trapSignals(cancelCtx, runSupervisor); err != nil {
		fmt.Printf("Got error: %s", err)
		cancel()
		os.Exit(1)
	} else {
		cancel()
		os.Exit(0)
	}
}

func trapSignals(origCtx context.Context, runFn func(context.Context)) error {
	ctx, stop := signal.NotifyContext(origCtx, syscall.SIGINT, syscall.SIGTERM)

	go runFn(ctx)

	<-ctx.Done()

	stop()
	return ctx.Err()
}

func runSupervisor(ctx context.Context) {
	var cl cli.CommandListener
	var sv supervisor.Supervisor

	if logger, err := initLogger(ctx); err != nil {
		return
	} else {
		if err := cl.Init(&sv, logger); err != nil {
			fmt.Printf("Error initializing command listener: %s\n", err)
			return
		} else {
			go cl.Listen()

			<-ctx.Done()

			sv.Stop("")
		}
	}
}

func initLogger(ctx context.Context) (*zap.SugaredLogger, error) {
	loggerType := os.Getenv("LOGGER")

	if l, err := utils.NewLogger(ctx, loggerType); err != nil {
		return nil, fmt.Errorf("Error initializing logger: %s\n", err)
	} else {
		return l.Sugar(), nil
	}
}
