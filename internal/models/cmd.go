package models

import (
	"regexp"
	"runtime"
	"strings"

	"github.com/brezzgg/go-packages/lg"
	"gopkg.in/yaml.v3"
)

type CmdSource struct {
	YamlSliceSource[*Command]
}

func (c *CmdSource) Compile(ctx *VarContext) ([]string, error) {
	if c == nil || c.Len() == 0 {
		return nil, nil
	}

	cmds := c.GetSource()
	result := make([]string, 0, len(cmds))

	for i, cmd := range cmds {
		compiled, err := cmd.Compile(ctx)
		if err != nil {
			return nil, lg.Ef("command[%d]: %w", i, err)
		}

		if compiled != "" {
			result = append(result, compiled)
		}
	}

	return result, nil
}

type Command struct {
	Raw  string
	Vars []CommandVar
	Os   string
}

func (c *Command) Compile(ctx *VarContext) (string, error) {
	if c == nil {
		return "", lg.Ef("command is nil")
	}

	if c.Os != "" {
		if !strings.EqualFold(c.Os, runtime.GOOS) {
			return "", nil
		}
	}

	if len(c.Vars) == 0 {
		return c.Raw, nil
	}

	result := c.Raw

	for _, cmdVar := range c.Vars {
		val, ok := ctx.Get(cmdVar.Name)
		if !ok {
			return "", lg.Ef("variable %q required but not found", cmdVar.Name)
		}
		result = strings.ReplaceAll(result, cmdVar.Raw, val)
	}

	return result, nil
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

	varRe := regexp.MustCompile(`\$\{\{\s*([a-zA-Z0-9._-]{1,256})\s*\}\}`)

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

var _ yaml.Unmarshaler = (*Command)(nil)

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
