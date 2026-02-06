//go:build windows
package handlers

import (
	"slices"
	"context"
	"os"
	"path/filepath"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
)

func execHandler(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
	return func(ctx context.Context, args []string) error {
		hc := interp.HandlerCtx(ctx)
		cmdName := args[0]
		if hasWindowsExtension(cmdName) {
			return next(ctx, args)
		}
		resolvedPath, err := lookPathWithExt(hc.Dir, hc.Env, cmdName)
		if err == nil {
			args[0] = resolvedPath
		}
		return next(ctx, args)
	}
}

func hasWindowsExtension(name string) bool {
    ext := strings.ToLower(filepath.Ext(name))
    extensions := []string{".exe", ".com", ".bat", ".cmd", ".ps1"}
    return slices.Contains(extensions, ext)
}

func lookPathWithExt(dir string, env expand.Environ, name string) (string, error) {
    pathext := env.Get("PATHEXT").String()
    if pathext == "" {
        pathext = ".COM;.EXE;.BAT;.CMD;.VBS;.VBE;.JS;.JSE;.WSF;.WSH;.MSC"
    }
    
    extensions := strings.Split(pathext, ";")
    
    path, err := interp.LookPathDir(dir, env, name)
    if err == nil {
        return path, nil
    }
    
    for _, ext := range extensions {
        nameWithExt := name + strings.ToLower(ext)
        path, err := interp.LookPathDir(dir, env, nameWithExt)
        if err == nil {
            return path, nil
        }
    }
    
    return "", os.ErrNotExist
}
