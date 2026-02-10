//go:build windows

package handlers

import (
	"context"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/kolesnikovae/go-winjob"
	"golang.org/x/sys/windows"
)

func startCmd(ctx context.Context, cmd *exec.Cmd) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_SUSPENDED,
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	job, err := winjob.Create("")
	if err != nil {
		cmd.Process.Kill()
		return &ErrStartProcessGroup{err.Error()}
	}
	defer job.Close()

	if err := job.Assign(cmd.Process); err != nil {
		cmd.Process.Kill()
		return &ErrStartProcessGroup{err.Error()}
	}

	if err := winjob.ResumeProcess(cmd.Process.Pid); err != nil {
		cmd.Process.Kill()
		return &ErrStartProcessGroup{err.Error()}
	}

	waitDone := make(chan error, 1)
	go func() {
		waitDone <- cmd.Wait()
	}()

	select {
	case err := <-waitDone:
		return err

	case <-ctx.Done():
		if KillTimeout > 0 {
			cmd.Process.Signal(os.Interrupt)
			select {
			case <-time.After(KillTimeout):
				job.Terminate()
			case <-waitDone:
				return ctx.Err()
			}
		} else {
			job.Terminate()
		}

		select {
		case <-waitDone:
		case <-time.After(2 * time.Second):
		}

		return ctx.Err()
	}
}
