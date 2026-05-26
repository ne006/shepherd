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

	fmt.Println(ctx.Err())
	stop()

	return nil
}

func runSupervisor(ctx context.Context) {
	var cl cli.CommandListener
	var sv supervisor.Supervisor

	cl.Supervisor = &sv

	if err := cl.Init(); err != nil {
		fmt.Printf("Error initializing command listener: %s\n", err)
		return
	} else {
		go cl.Listen()

		<-ctx.Done()

		sv.Stop("")
	}
}
