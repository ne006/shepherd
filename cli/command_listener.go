package cli

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/ne006/shepherd/config_loader"
	"github.com/ne006/shepherd/supervisor"
	"github.com/ne006/shepherd/utils"
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
		if homeDir, err := os.UserHomeDir(); err == nil {
			cl.socketPath = filepath.Join(homeDir, ".shepherd", "shepherd.sock")
		} else {
			panic(fmt.Sprintf("Could not obtain user home directory %s", err))
		}
	}
}

func (cl *CommandListener) initSocket() error {
	if socketDir := filepath.Dir(cl.socketPath); !utils.FileExists(socketDir) {
		os.MkdirAll(socketDir, 0766)
	}

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

// Command handling
var commands = map[string]func(*CommandListener, []string) (string, error){
	"start": Start,
	"stop":  Stop,
	"list":  List,
	"load":  Load,
}

func (cl *CommandListener) processCommand(input string) (string, error) {
	cmdName, args := parseCommand(input)

	if cmd := commands[cmdName]; cmd == nil {
		return "", fmt.Errorf("input not recognized")
	} else {
		return cmd(cl, args)
	}
}

func parseCommand(input string) (string, []string) {
	splitInput := strings.Split(strings.TrimRight(input, "\r\n"), " ")

	cmdName := splitInput[0]
	args := splitInput[1:]

	return cmdName, args
}

// Commands
func Start(cl *CommandListener, _ []string) (string, error) {
	if err := cl.Supervisor.Start(); err != nil {
		return "", err
	} else {
		return "ok\n", nil
	}
}

func Stop(cl *CommandListener, _ []string) (string, error) {
	if err := cl.Supervisor.Stop(); err != nil {
		return "", err
	} else {
		return "ok\n", nil
	}
}

func List(cl *CommandListener, _ []string) (string, error) {
	return fmt.Sprintf("%+v\n", cl.Supervisor), nil
}

func Load(cl *CommandListener, args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("Should pass a path to supervision config")
	}

	configPath := args[0]

	if configPath == "" {
		return "", fmt.Errorf("Should pass a path to supervision config")
	}

	if !utils.FileExists(configPath) {
		return "", fmt.Errorf("%s does not exist\n", configPath)
	}

	if app, err := config_loader.LoadConfig(configPath); err != nil {
		return "", fmt.Errorf("Error loading %s: %s\n", configPath, err)
	} else {
		cl.Supervisor.LoadApp(*app)
		cl.Supervisor.Start()

		return "ok\n", nil
	}
}
