package handlers

import "mvdan.cc/sh/v3/interp"

func FormatHandler(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
	return formatHandler(next)
}
