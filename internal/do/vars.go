package do

import (
	"os"
	"runtime"
	"strings"

	"github.com/brezzgg/delease/internal/models"
)

type OsVar string

const (
	OsVarArgs OsVar = "os.args"
	OsVarOs   OsVar = "os.os"
	OsVarArch OsVar = "os.arch"
)

func (v OsVar) String() string {
	return string(v)
}

func (d *Do) GetOsVars(args string) *models.VarSource {
	v := &models.VarSource{}
	v.SetSource(map[string]string{
		OsVarArgs.String(): args,
		OsVarOs.String():   runtime.GOOS,
		OsVarArch.String(): runtime.GOARCH,
	})
	return v
}

func (d *Do) GetEnv() *models.EnvSource {
	envs := os.Environ()
	m := make(map[string]string, 10)
	for _, v := range envs {
		if before, after, ok := strings.Cut(v, "="); ok {
			m[before] = after
		}
	}
	r := &models.EnvSource{}
	r.SetSource(m)
	return r
}
