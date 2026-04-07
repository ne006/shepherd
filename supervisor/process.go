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
	startErr := process.Cmd.Start()

	if startErr == nil {
		go process.supervise()
	}

	return startErr
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

func (process *Process) supervise() error {
	if process.Process == nil {
		return fmt.Errorf(("Associated process does not exist"))
	}

	pid := process.Process.Pid

	if state, err := process.Cmd.Process.Wait(); err != nil {
		fmt.Printf("%v wait error: %s\n", pid, err)
	} else {
		fmt.Printf("%v exited: %+v\n", pid, state)
	}

	if err := process.recreate(); err != nil {
		return fmt.Errorf("Process recreation failed: %s", err)
	}

	err := process.Start()

	fmt.Printf("%v restarted: %s\n", pid, err)

	return err
}

func (process *Process) recreate() error {
	newCmd := exec.Cmd{
		Path: process.Path,
	}

	if process.Args != nil {
		newCmd.Args = process.Args
	}

	process.Cmd = newCmd

	return nil
}
