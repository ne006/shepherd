package main

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/ne006/shepherd/cli"
)

func main() {
	cmd := os.Args[1]

	if cmd == "" {
		fmt.Println(usage())
		os.Exit(1)
	}

	client := cli.Client{}

	client.Init()

	if resp, err := client.SendCommand(cmd); err != nil {
		fmt.Printf("An error occured: %s\n", err)
	} else {
		fmt.Printf("%s", resp)
	}
}

func usage() string {
	return heredoc.Doc(`
	shepherd - a process supervisor

	Usage:
		- start
		- stop
		- list
	`)
}
