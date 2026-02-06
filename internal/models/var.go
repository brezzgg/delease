package models

type VarSource struct {
	YamlMapSource[string]
}

func (v *VarSource) Merge(oth *VarSource, force bool) *VarSource {
	if r := PremergeCheck(v, oth); r != nil {
		return r
	}

	merged := v.YamlMapSource.Merge(oth.YamlMapSource, force)
	return &VarSource{YamlMapSource: merged}
}

var _ Mergeable[*VarSource] = (*VarSource)(nil)
