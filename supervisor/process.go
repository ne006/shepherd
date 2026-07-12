package supervisor

import (
	"fmt"
	"os/exec"
	"slices"
	"syscall"
	"time"

	"go.uber.org/zap"
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
	stateTime   time.Time

	Logger *zap.SugaredLogger

	envHash map[string]any
}

func (process *Process) GetName() string {
	return process.Name
}

func (process *Process) SetName(name string) error {
	process.Name = name

	return nil
}

func (process *Process) GetEnv() map[string]any {
	return process.envHash
}

func (process *Process) SetEnv(envHash map[string]any) {
	process.envHash = envHash
}

func (process *Process) Start() error {
	if state, _ := process.getState(); state == StateStopped {
		if err := process.recreate(); err != nil {
			process.setState(StateStopped, StateReasonError)
			return err
		}
	}

	process.Env = loadEnv(process)

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

func (process *Process) Restart(reason ProcessStateReason) error {
	if err := process.Stop(reason); err != nil {
		return err
	} else if err := process.Start(); err != nil {
		return err
	} else {
		return nil
	}
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
		if pstate, _ := process.getState(); pstate != StateStopped && pstate != StateStopping {
			process.setState(StateStopped, StateReasonError)
			process.Logger.Errorf("%v wait error: %s", pid, err)
		}
	} else {
		if _, stateReason := process.getState(); stateReason == StateReasonNone {
			if state.Success() {
				process.setState(StateStopped, StateReasonExited)
			} else {
				process.setState(StateStopped, StateReasonError)
			}
		}
		process.Logger.Infof("%v exited with %v: %+v", pid, process.stateReason, state)
	}

	if _, stateReason := process.getState(); stateReason != StateReasonUser {
		err := process.restart()

		process.Logger.Infof("%v restarted: %s", pid, err)

		return err
	} else {
		return nil
	}
}

func (process *Process) restart() error {
	if err := process.recreate(); err != nil {
		process.Logger.Infof("Process recreation failed: %s", err)
		return fmt.Errorf("Process recreation failed: %s", err)
	}

	return process.Start()
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
		process.Logger.Infof("%s: %s (%s) -> %s (%s)", process.Name, process.state, process.stateReason, to, reason)
		process.state = to
		process.stateReason = reason
		process.stateTime = time.Now()

		return nil
	} else {
		return fmt.Errorf("Can't transition from %s to %s", process.state, to)
	}
}
