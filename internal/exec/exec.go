package exec

import (
	"context"
)

type Executor interface {
	Setup(wd string, lines, env []string, log Logger) error
	RunLine(ctx context.Context, line int) Result
	Run(ctx context.Context, resChan chan<- Result)
}

type Result struct {
	Code  int
	Error error
}

type (
	Logger  func(m string, t MsgType)
	MsgType int
)

const (
	MsgTypeStdout MsgType = iota
	MsgTypeStderr
	MsgTypeInternal
)
