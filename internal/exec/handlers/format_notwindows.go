//go:build !windows
package handlers

import (
	"context"

	"mvdan.cc/sh/v3/interp"
)

func formatHandler(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
	return func(ctx context.Context, args []string) error {
		return next(ctx, args)
	}
}
