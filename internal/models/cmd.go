package models

import (
	"regexp"
	"strings"

	"github.com/brezzgg/go-packages/lg"
)

var varRe = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9._-]{1,256})\s*\}\}`)

type CmdSource struct {
	YamlSliceSource[string]
}

func (c *CmdSource) ApplyVars(vars *VarSource) (*CmdSource, error) {
	if c == nil {
		return nil, nil
	}
	if vars == nil {
		vars = &VarSource{}
	}
	cmds := c.GetCopy()
	for i, cmd := range cmds {
		res, err := c.replace(cmd, vars)
		if err != nil {
			return nil, lg.Ef("cmds line %d, error: %w", i, err)
		}
		cmds[i] = res
	}
	res := &CmdSource{}
	res.SetSource(cmds)
	return res, nil
}

func (c *CmdSource) replace(cmd string, vars *VarSource) (string, error) {
	var (
		notFound []string
		varcp    = vars.GetMapCopy()
	)
	res := varRe.ReplaceAllStringFunc(cmd, func(s string) string {
		match := varRe.FindStringSubmatch(s)
		v, ok := varcp[match[1]]
		if !ok {
			notFound = append(notFound, match[1])
			return s
		} else {
			return v
		}
	})
	if len(notFound) > 0 {
		return "", lg.Ef("var(s) '%s' required", strings.Join(notFound, ","))
	}
	return res, nil
}
