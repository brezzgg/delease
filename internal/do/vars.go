package do

import (
	"runtime"

	"github.com/brezzgg/delease/internal/models"
)

const (
	OsVarArgs = "os.args"
	OsVarOs   = "os.os"
	OsVarArch = "os.arch"
)

func (d *Do) GetOsVars(args string) *models.VarSource {
	v := &models.VarSource{}
	v.SetSource(map[string]string{
		OsVarArgs: args,
		OsVarOs: runtime.GOOS,
		OsVarArch: runtime.GOARCH,
	})
	return v
}
