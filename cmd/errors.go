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
	return lg.Ef("parse failed: %s", text)
}

func ErrApplyVars(err error) error {
	return lg.Ef("vars apply failed: %w", err)
}

func ErrTaskNotFound(err error) error {
	return lg.Ef("task not found: %w", err)
}
