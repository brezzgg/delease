//go:build !windows
package handlers

import (
	"context"

	"mvdan.cc/sh/v3/interp"
)

func execHandler(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
	return func(ctx context.Context, args []string) error {
		return next(ctx, args)
	}
}
