//go:build !windows

package handlers

import (
	"context"
	"os/exec"
	"syscall"
	"time"
)

func startCmd(ctx context.Context, cmd *exec.Cmd) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	waitDone := make(chan error, 1)
	go func() {
		waitDone <- cmd.Wait()
	}()

	select {
	case err := <-waitDone:
		return err

	case <-ctx.Done():
		pgid := cmd.Process.Pid

		if KillTimeout > 0 {
			syscall.Kill(-pgid, syscall.SIGTERM)
			select {
			case <-time.After(KillTimeout):
				syscall.Kill(-pgid, syscall.SIGKILL)
			case <-waitDone:
				return ctx.Err()
			}
		} else {
			syscall.Kill(-pgid, syscall.SIGKILL)
		}

		select {
		case <-waitDone:
		case <-time.After(2 * time.Second):
		}

		return ctx.Err()
	}
}
