package lib

import (
	"bytes"
	"errors"
	"os/exec"

	"github.com/google/uuid"
)

type Status string

const (
	ActiveStatus  Status = "Active"
	ExitedStatus         = "Exited"
	KilledStatus         = "Killed"
	StoppedStatus        = "Stopped"
)

type Work struct {
	ID       uuid.UUID
	ClientID uuid.UUID
	Status   Status
	Output   bytes.Buffer
	Command  []string
	Process  *exec.Cmd
}

func NewWork(clientID uuid.UUID, command []string) *Work {
	return &Work{
		ID:       uuid.New(),
		ClientID: uuid.New(),
		Command:  command,
		Process:  exec.Command(command[0], command[1:]...),
	}
}

func (w *Work) Start() error {
	w.Process.Stdout = &w.Output
	w.Process.Stderr = &w.Output
	err := w.Process.Start()
	if err != nil {
		return err
	}
	w.Status = ActiveStatus
	processErr := w.Process.Wait()
	if processErr == nil {
		w.Status = KilledStatus
		return processErr
	}
	return nil
}

func (w *Work) QueryStatus() Status {
	return w.Status
}

func (w *Work) Stop() error {
	if w.Status != ActiveStatus {
		return errors.New("This work is not running already")
	}
	if err := w.Process.Process.Kill(); err != nil {
		return errors.New("Failed to kill the process")
	}
	w.Status = StoppedStatus
	return nil
}

func (w *Work) Log() string {
	return w.Output.String()
}
