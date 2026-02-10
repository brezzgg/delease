package exec

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/brezzgg/delease/internal/exec/handlers"
	"github.com/brezzgg/go-packages/lg"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type Sh struct {
	runner *interp.Runner
	lines  []string
	ow     *SyncWriter
	ew     *SyncWriter
}

func (s *Sh) Setup(wd string, lines, env []string, log Logger) error {
	s.ow, s.ew = NewSyncWriter(log, MsgTypeStdout), NewSyncWriter(log, MsgTypeStderr)

	handlers.KillTimeout = time.Millisecond * 750

	runner, err := interp.New(
		interp.Dir(wd),
		interp.Env(expand.ListEnviron(env...)),
		interp.StdIO(nil, s.ow, s.ew),
		interp.ExecHandlers(handlers.Get()...),
	)
	if err != nil {
		return err
	}

	s.runner = runner
	s.lines = lines
	return nil
}

func (s *Sh) RunLine(ctx context.Context, line int) Result {
	if line < 0 || line >= len(s.lines) {
		return Result{
			Error: lg.Ef("line %d out of bounds (total: %d)", line, len(s.lines)),
		}
	}

	file, err := syntax.NewParser().Parse(strings.NewReader(s.lines[line]), "")
	if err != nil {
		return Result{
			Error: lg.Ef("parse error: %w", err),
		}
	}

	var msgSb strings.Builder
	msgSb.WriteString(lg.T(s.lines[line]))
	split := strings.Split(msgSb.String(), "\n")
	msgSb.Reset()
	for i, v := range split {
		if i == 0 {
			msgSb.WriteString("$ ")
			msgSb.WriteString(v)
		} else {
			msgSb.WriteString("  ")
			msgSb.WriteString(v)
		}
		msgSb.WriteString("\n")
	}

	s.ow.log(msgSb.String(), MsgTypeStdout)
	err = s.runner.Run(ctx, file)

	var exitErr interp.ExitStatus
	if errors.As(err, &exitErr) {
		return Result{
			Code: int(exitErr),
		}
	}
	if err != nil {
		return Result{
			Error: lg.Ef("runner error: %w", err),
		}
	}
	return Result{}
}

func (s *Sh) Run(ctx context.Context, resChan chan<- Result) {
	var last Result
	for i := range s.lines {
		if ctx.Err() != nil {
			resChan <- Result{Error: ctx.Err()}
			return
		}

		last = s.RunLine(ctx, i)
		if last.Error != nil || last.Code != 0 {
			resChan <- last
			return
		}
	}
	s.ow.Flush()
	s.ew.Flush()
	resChan <- last
}
