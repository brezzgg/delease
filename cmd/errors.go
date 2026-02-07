package cmd

import (
	"strings"

	"github.com/brezzgg/go-packages/lg"
)

func ErrBadConfig(err error) error {
	return lg.Ef("bad config: %w", err)
}

func ErrParseFailed(err error) error {
	text := strings.ReplaceAll(err.Error(), "\n", " ")
	return lg.Ef("parse error: %s", text)
}

func ErrCompileVars(err error) error {
	return lg.Ef("vars compile failed: %w", err)
}

func ErrTaskNotFound(name string) error {
	return lg.Ef("task not found: %s", name)
}
