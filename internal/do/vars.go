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

func GetOsVars(args string) *models.VarSource {
	v := &models.VarSource{}
	v.SetSource(map[string]*models.Var{
		OsVarArgs.String(): models.NewVar(args, models.VarTypeOs),
		OsVarOs.String():   models.NewVar(runtime.GOOS, models.VarTypeOs),
		OsVarArch.String(): models.NewVar(runtime.GOARCH, models.VarTypeOs),
	})
	return v
}

func GetEnv() *models.EnvSource {
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
