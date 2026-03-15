package cli

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ne006/shepherd/supervisor"
)

type CommandListener struct {
	socketPath string
	socket     net.Listener

	Supervisor *supervisor.Supervisor
}

func (cl *CommandListener) Init() error {
	cl.setDefaults()

	if err := cl.initSocket(); err != nil {
		return err
	}

	return nil
}

func (cl *CommandListener) Listen() error {
	if err := cl.listenSocket(); err != nil {
		return err
	}

	return nil
}

func (cl *CommandListener) setDefaults() {
	if cl.socketPath == "" {
		cl.socketPath = "/tmp/shepherd.sock"
	}
}

func (cl *CommandListener) initSocket() error {
	socket, err := net.Listen("unix", cl.socketPath)
	if err != nil {
		return err
	}

	cl.socket = socket

	// Cleanup the sockfile.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(cl.socketPath)
	}()

	return nil
}

func (cl *CommandListener) listenSocket() error {
	for {
		// Accept an incoming connection.
		conn, err := cl.socket.Accept()
		if err != nil {
			fmt.Println(err)
		}

		// Handle the connection in a separate goroutine.
		go func(conn net.Conn) {
			defer conn.Close()
			// Create a buffer for incoming data.
			buf := make([]byte, 4096)

			// Read data from the connection.
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println(err)
			}

			if resp, cmdErr := cl.processCommand(string(buf[:n])); cmdErr == nil {
				// Echo the data back to the connection.
				_, err = conn.Write([]byte(resp))
				if err != nil {
					fmt.Println(err)
				}
			} else {
				// Echo the data back to the connection.
				_, err = conn.Write([]byte(fmt.Sprintf("%v\n", cmdErr)))
				if err != nil {
					fmt.Println(err)
				}
			}

		}(conn)
	}
}

func (cl *CommandListener) processCommand(input string) (string, error) {
	switch input {
	case "start\n":
		if err := cl.Supervisor.Start(); err != nil {
			return "", err
		} else {
			return "ok\n", nil
		}
	case "stop\n":
		if err := cl.Supervisor.Stop(); err != nil {
			return "", err
		} else {
			return "ok\n", nil
		}
	case "list\n":
		return fmt.Sprintf("%+v\n", cl.Supervisor), nil
	default:
		return "", fmt.Errorf("input not recognized")
	}
}
