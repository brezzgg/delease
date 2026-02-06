package handlers

import "mvdan.cc/sh/v3/interp"

func ExecHandler(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
	return execHandler(next)
}
