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
	process.setState(StateStarting, StateReasonNone)

	if startErr == nil {
		process.setState(StateRunning, StateReasonNone)
		go process.supervise()
	} else {
		process.setState(StateStopped, StateReasonError)
	}

	return startErr
}

func (process *Process) Stop(reason ProcessStateReason) error {
	if process.Process == nil {
		process.setState(StateStopped, StateReasonExited)
		return fmt.Errorf(("Associated process does not exist"))
	}

	pid := process.Process.Pid

	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return err
	}

	process.setState(StateStopping, reason)

	// Wait for the process to exit
	process.Wait()

	process.setState(StateStopped, reason)
	return nil
}

func (process *Process) GetUIString() string {
	var processState string

	if process.stateReason == StateReasonNone {
		processState = fmt.Sprintf("%s", process.state)
	} else {
		processState = fmt.Sprintf("%s (%s)", process.state, process.stateReason)
	}

	if process.Process == nil {
		return fmt.Sprintf("%s %s", process.Name, processState)
	} else {
		return fmt.Sprintf("%s %s %v\n", process.Name, processState, process.Process.Pid)
	}
}

func (process *Process) supervise() error {
	if process.Process == nil {
		return fmt.Errorf(("Associated process does not exist"))
	}

	pid := process.Process.Pid

	if state, err := process.Cmd.Process.Wait(); err != nil {
		process.setState(StateStopped, StateReasonError)
		fmt.Printf("%v wait error: %s\n", pid, err)
	} else {
		process.setState(StateStopped, StateReasonExited)
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
