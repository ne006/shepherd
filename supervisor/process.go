package supervisor

import (
	"os/exec"
)

type Process struct {
	exec.Cmd

	Name string
}

func (process Process) GetName() string {
	return process.Name
}

func (process Process) SetName(name string) error {
	process.Name = name

	return nil
}

func (process Process) Start() error {
	return process.Cmd.Start()
}

func (process Process) Stop() error {
	// TODO: implement
	return nil
}
