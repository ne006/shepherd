package supervisor

import (
	"fmt"
	"os/exec"
	"slices"
	"syscall"
)

// ProcessState
type ProcessState int

const (
	StateInitialized ProcessState = iota
	StateStarting
	StateRunning
	StateStopping
	StateStopped
)

var processStateName = map[ProcessState]string{
	StateInitialized: "initialized",
	StateStarting:    "starting",
	StateRunning:     "running",
	StateStopping:    "stopping",
	StateStopped:     "stopped",
}

var processStateTransitions = map[ProcessState][]ProcessState{
	StateInitialized: {StateStarting},
	StateStarting:    {StateRunning, StateStopped},
	StateRunning:     {StateStopping, StateStopped},
	StateStopping:    {StateStopped},
	StateStopped:     {StateStarting},
}

func (ps ProcessState) String() string {
	return processStateName[ps]
}

func (from ProcessState) TransitionToIsValid(to ProcessState) bool {
	return slices.Contains(processStateTransitions[from], to)
}

// ProcessStateReason
type ProcessStateReason int

const (
	StateReasonNone ProcessStateReason = iota
	StateReasonExited
	StateReasonUser
	StateReasonError
)

var processStateReasonName = map[ProcessStateReason]string{
	StateReasonNone:   "",
	StateReasonExited: "exited",
	StateReasonUser:   "user",
	StateReasonError:  "error",
}

func (psr ProcessStateReason) String() string {
	return processStateReasonName[psr]
}

// Process
type Process struct {
	exec.Cmd

	Name        string
	state       ProcessState
	stateReason ProcessStateReason
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

func (process *Process) GetUIString() string {
	if process.Process == nil {
		return fmt.Sprintf("%s <not started>", process.Name)
	} else {
		return fmt.Sprintf("%s %v\n", process.Name, process.Process.Pid)
	}
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

func (process *Process) getState() (ProcessState, ProcessStateReason) {
	return process.state, process.stateReason
}

func (process *Process) setState(to ProcessState, reason ProcessStateReason) error {
	if process.state.TransitionToIsValid(to) {
		process.state = to
		process.stateReason = reason

		return nil
	} else {
		return fmt.Errorf("Can't transition from %s to %s", process.state, to)
	}
}
