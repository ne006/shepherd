package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ne006/shepherd/cli"
	"github.com/ne006/shepherd/supervisor"
)

type supervisorKey string

func main() {
	baseCtx := context.Background()
	cancelCtx, cancel := context.WithCancel(baseCtx)
	valueCtx := context.WithValue(cancelCtx, supervisorKey("supervisor"), supervisor.Supervisor{})

	if err := trapSignals(valueCtx, cancel, runSupervisor, stopSupervisor); err != nil {
		fmt.Printf("Got error: %s", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func trapSignals(origCtx context.Context, cancel context.CancelFunc, runFn func(context.Context, context.CancelFunc), stopFn func(context.Context)) error {
	ctx, stop := signal.NotifyContext(origCtx, syscall.SIGINT, syscall.SIGTERM)

	go runFn(ctx, cancel)

	<-ctx.Done()

	fmt.Println(ctx.Err())
	stopFn(ctx)
	stop()

	return nil
}

func getSupervisor(ctx context.Context) (*supervisor.Supervisor, error) {

	if v := ctx.Value(supervisorKey("supervisor")); v != nil {
		if s, ok := v.(supervisor.Supervisor); ok {
			return &s, nil
		} else {
			return nil, fmt.Errorf("Error initializing supervisor")
		}
	} else {
		return nil, fmt.Errorf("Error initializing supervisor")
	}
}

func runSupervisor(ctx context.Context, cancel context.CancelFunc) {
	var cl cli.CommandListener
	var sv supervisor.Supervisor

	if s, err := getSupervisor(ctx); err != nil {
		fmt.Printf("%s\n", err)
		cancel()
		return
	} else {
		sv = *s
	}

	cl.Supervisor = &sv

	if err := cl.Init(); err != nil {
		fmt.Printf("Error initializing command listener: %s\n", err)
		cancel()
		return
	} else {
		go cl.Listen()
	}
}

func stopSupervisor(ctx context.Context) {
	var sv supervisor.Supervisor

	if s, err := getSupervisor(ctx); err != nil {
		fmt.Printf("%s\n", err)
		return
	} else {
		sv = *s
	}

	sv.Stop()

}
