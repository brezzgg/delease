package models

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type VarSource struct {
	YamlMapSource[*Var]
}

func (v *VarSource) Merge(oth *VarSource, force bool) *VarSource {
	if r := PremergeCheck(v, oth); r != nil {
		return r
	}

	merged := v.YamlMapSource.Merge(oth.YamlMapSource, force)
	return &VarSource{YamlMapSource: merged}
}

func (v *VarSource) UnmarshalYAML(value *yaml.Node) error {
	var data map[string]string
	if err := value.Decode(&data); err != nil {
		return err
	}

	result := make(map[string]*Var, len(data))
	for k, v := range data {
		result[k] = NewVarT(k, v)
	}
	v.SetSource(result)

	return nil
}

var (
	_ Mergeable[*VarSource] = (*VarSource)(nil)
	_ yaml.Unmarshaler      = (*VarSource)(nil)
)

type VarType string

const (
	VarTypeStatic  = "static"
	VarTypeDynamic = "dynamic"
	VarTypeOs      = "os"
	VarTypeAny     = "any"
)

type Var struct {
	Type    VarType
	Content string
}

func NewVar(raw string, varType VarType) *Var {
	return &Var{
		Type:    varType,
		Content: raw,
	}
}

func NewVarT(name, raw string) *Var {
	return &Var{
		Type:    GetVarType(name),
		Content: raw,
	}
}

func (v *Var) IsType(varType VarType) bool {
	if varType == VarTypeAny {
		return true
	}
	return v.Type == varType
}

func GetVarType(raw string) VarType {
	switch {
	case strings.HasPrefix(raw, "."):
		return VarTypeDynamic
	case strings.HasPrefix(raw, "os."):
		return VarTypeOs
	default:
		return VarTypeStatic
	}
}
