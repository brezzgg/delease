package handlers

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
)

type ErrStartProcessGroup struct {
	Msg string
}

func (e *ErrStartProcessGroup) Error() string {
	return fmt.Sprintf("process group: %s", e.Msg)
}

var KillTimeout = 1 * time.Second

func StartHandler(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
	return func(ctx context.Context, args []string) error {
		hc := interp.HandlerCtx(ctx)
		path, err := interp.LookPathDir(hc.Dir, hc.Env, args[0])
		if err != nil {
			fmt.Fprintf(hc.Stderr, "%s: %v\n", filepath.Base(args[0]), err)
			return interp.ExitStatus(127)
		}
		cmd := exec.Cmd{
			Path:   path,
			Args:   args,
			Env:    execEnv(hc.Env),
			Dir:    hc.Dir,
			Stdin:  hc.Stdin,
			Stdout: hc.Stdout,
			Stderr: hc.Stderr,
		}

		return handleError(startCmd(ctx, &cmd), hc)
	}
}

func handleError(err error, hc interp.HandlerContext) error {
	switch err := err.(type) {
	case *exec.ExitError:
		return interp.ExitStatus(err.ExitCode())
	case *exec.Error:
		fmt.Fprintf(hc.Stderr, "%v\n", err)
		return interp.ExitStatus(127)
	default:
		return err
	}
}

func execEnv(env expand.Environ) []string {
	list := make([]string, 0, 64)
	for name, vr := range env.Each {
		if !vr.IsSet() {
			for i, kv := range list {
				if strings.HasPrefix(kv, name+"=") {
					list[i] = ""
				}
			}
		}
		if vr.Exported && vr.Kind == expand.String {
			list = append(list, name+"="+vr.String())
		}
	}
	return list
}
