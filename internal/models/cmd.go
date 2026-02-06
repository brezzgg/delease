package models

import (
	"regexp"
	"runtime"
	"strings"

	"github.com/brezzgg/go-packages/lg"
	"gopkg.in/yaml.v3"
)

var varRe = regexp.MustCompile(`\$\{\{\s*([a-zA-Z0-9._-]{1,256})\s*\}\}`)

type CmdSource struct {
	YamlSliceSource[*Command]
	applied bool
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
		res, err := cmd.ApplyVars(vars)
		if err != nil {
			return nil, lg.Ef("cmds line %d, error: %w", i, err)
		}
		if res != "" {
			c := &Command{
				Cmd: res,
			}
			cmds[i] = c
		}
	}
	res := &CmdSource{}
	res.SetSource(cmds)
	return res, nil
}

func (c *CmdSource) Applied() bool {
	return c.applied
}

var _ Applier[*CmdSource] = (*CmdSource)(nil)

type Command struct {
	Raw  string
	Vars []CommandVar
	Os   string
	Cmd  string
}

func (c *Command) ApplyVars(vars *VarSource) (string, error) {
	if c == nil {
		return "", lg.Ef("command is nil")
	}
	if c.Os != "" {
		if !strings.EqualFold(c.Os, runtime.GOOS) {
			return "", nil
		}
	}
	if len(c.Vars) == 0 {
		return strings.Clone(c.Raw), nil
	}

	res := strings.Clone(c.Raw)

	for _, cmdVar := range c.Vars {
		if val, ok := vars.Get(cmdVar.Name); ok {
			res = strings.ReplaceAll(res, cmdVar.Raw, val)
		} else {
			return "", lg.Ef("var %s required", cmdVar.Name)
		}
	}

	return res, nil
}

func (c *Command) Applied() bool {
	return len(c.Vars) == 0
}

func (c *Command) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.MappingNode:
		m := make(map[string]string)
		if err := value.Decode(&m); err != nil {
			return err
		}
		if len(m) == 0 {
			return lg.Ef("command cant be empty")
		}
		for k, v := range m {
			switch k {
			case "cmd", "run":
				c.Raw = v
			case "os":
				c.Os = v
			default:
				return lg.Ef("unknown filed %s in command", k)
			}
		}

	case yaml.ScalarNode:
		if err := value.Decode(&c.Raw); err != nil {
			return err
		}
	default:
		return lg.Ef("unexpected value kind")
	}

	return c.ParseVars()
}

func (c *Command) ParseVars() error {
	if c == nil {
		return lg.Ef("command is nil")
	}

	varRe.ReplaceAllStringFunc(c.Raw, func(raw string) string {
		match := varRe.FindStringSubmatch(raw)
		name := match[1]

		var t CommandVarType
		if strings.HasPrefix(name, ".") {
			name = name[1:]
			t = CmdVarTypeDynamic
		} else if strings.HasPrefix(name, "os.") {
			t = CmdVarTypeOs
		} else {
			t = CmdVarTypeStatic
		}

		c.Vars = append(c.Vars, CommandVar{
			Name: name,
			Raw:  raw,
			Type: t,
		})

		return raw
	})

	return nil
}

var (
	_ yaml.Unmarshaler = (*Command)(nil)
	_ Applier[string]  = (*Command)(nil)
)

type CommandVarType int

const (
	CmdVarTypeStatic CommandVarType = iota
	CmdVarTypeDynamic
	CmdVarTypeOs
)

type CommandVar struct {
	Name string
	Raw  string
	Type CommandVarType
}
