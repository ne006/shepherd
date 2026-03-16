package supervisor

import (
	"fmt"
	"os/exec"
	"syscall"
)

type Process struct {
	exec.Cmd

	Name string
}

func (process *Process) GetName() string {
	return process.Name
}

func (process *Process) SetName(name string) error {
	process.Name = name

	return nil
}

func (process *Process) Start() error {
	return process.Cmd.Start()
}

func (process *Process) Stop() error {
	if process.Process == nil {
		return fmt.Errorf(("Associated process does not exist"))
	}

	pid := process.Process.Pid

	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return err
	}

	// Wait for the process to exit
	process.Wait()
	return nil
}
