package cli

import (
	"net"
	"time"
)

type Client struct {
	socketPath string
	socket     net.Conn
}

func (cl *Client) Init() error {
	cl.setDefaults()

	if err := cl.initSocket(); err != nil {
		return err
	}

	return nil
}

func (cl *Client) SendCommand(input string) (string, error) {
	if _, err := cl.socket.Write([]byte(input)); err != nil {
		return "", err
	}

	resp := make([]byte, 4096)

	cl.socket.SetReadDeadline(time.Now().Add(10 * time.Second))

	if n, err := cl.socket.Read(resp); err != nil {
		return "", err
	} else {
		if n == 0 {
			return "", nil
		}

		return string(resp), nil
	}
}

func (cl *Client) setDefaults() {
	if cl.socketPath == "" {
		cl.socketPath = "/tmp/shepherd.sock"
	}
}

func (cl *Client) initSocket() error {
	socket, err := net.Dial("unix", cl.socketPath)
	if err != nil {
		return err
	}

	cl.socket = socket

	return nil
}
