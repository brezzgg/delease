package parser

import (
	"bytes"
	"maps"
	"regexp"
	"strings"

	"github.com/brezzgg/delease/internal/parser/models"
	"github.com/brezzgg/go-packages/lg"
	"gopkg.in/yaml.v3"
)

func Parse(b []byte) (*models.Root, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(b))
	decoder.KnownFields(true)

	root := &models.Root{}

	if err := decoder.Decode(root); err != nil {
		return nil, lg.Ef("unmarshal error: %w", err)
	}

	return root, nil
}

func SetVar(v models.Variabler, key, value string, force bool) {
	vars := v.GetVars()
	if force {
		vars[key] = value
	} else {
		if _, ok := vars[key]; !ok {
			vars[key] = value
		}
	}
}

func GetVar(v models.Variabler, key string) (string, bool) {
	val, ok := v.GetVars()[key]
	return val, ok
}

func MergeVars(l, r models.Vars, force bool) models.Vars {
	if l == nil && r == nil {
		return models.Vars{}
	}
	if l == nil && r != nil {
		return r
	}
	if l != nil && r == nil {
		return l
	}
	res := make(models.Vars, len(r)+len(l))
	maps.Copy(res, l)
	for k, v := range r {
		if force {
			res[k] = v
		} else {
			if _, ok := res[k]; !ok {
				res[k] = v
			}
		}
	}
	return res
}

func MergeVariablers(l, r models.Variabler, force bool) models.Vars {
	if l == nil && r == nil {
		return models.Vars{}
	}
	if l == nil && r != nil {
		return r.GetVars()
	}
	if l != nil && r == nil {
		return l.GetVars()
	}
	return MergeVars(l.GetVars(), r.GetVars(), force)
}

func ApplyVars(r *models.Root, task string, forceVars models.Vars) error {
	if r == nil {
		return lg.Ef("root is nil")
	}

	t, err := r.GetTask(task)
	if err != nil {
		return lg.Ef("get task: %w", err)
	}

	vars := MergeVars(r.Vars, t.Vars, true)
	vars = MergeVars(vars, forceVars, true)

	cmds := make([]string, 0, len(t.Cmds))
	for i, cmd := range r.Tasks[task].Cmds {
		res, err := ReplaceCmd(cmd, vars)
		if err != nil {
			return lg.Ef("task '%s': cmds line %d, error: %w", task, i, err)
		}
		cmds = append(cmds, res)
	}

	r.Tasks[task].Cmds = cmds
	return nil
}

var varRe = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9._-]{1,256})\s*\}\}`)

func ReplaceCmd(cmd string, vars models.Vars) (string, error) {
	var notFound []string
	res := varRe.ReplaceAllStringFunc(cmd, func(s string) string {
		match := varRe.FindStringSubmatch(s)
		v, ok := vars[match[1]]
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
