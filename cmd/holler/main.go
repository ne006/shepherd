package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/ne006/shepherd/cli"
)

func main() {
	cmd := strings.Join(os.Args[1:], " ")

	if cmd == "" {
		fmt.Println(usage())
		os.Exit(1)
	}

	client := cli.Client{}

	if err := initClient(&client, true); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

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
		- restart
		- stop
		- list
		- load <configPath>
	`)
}

func initClient(client *cli.Client, initBackend bool) error {
	if err := client.Init(); err != nil {
		if client.SocketExists() {
			return fmt.Errorf("An error occured: %s\n", err)
		} else {
			if initBackend {
				if err := runBackend(); err != nil {
					return fmt.Errorf("An error occured when starting shepherd: %s\n", err)
				} else {
					return initClient(client, false)
				}
			} else {
				return initClient(client, false)
			}
		}
	}

	return nil
}

func runBackend() error {
	cmd := exec.Command("shepherd")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	return cmd.Start()
}
