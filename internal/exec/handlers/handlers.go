package handlers

import "mvdan.cc/sh/v3/interp"

func Get() []func(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
	return []func(next interp.ExecHandlerFunc) interp.ExecHandlerFunc{
		FormatHandler,
		StartHandler,
	}
}
