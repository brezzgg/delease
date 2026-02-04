package models

type VarSource struct {
	YamlMapSource[string]
}

func (v *VarSource) Merge(oth *VarSource, force bool) *VarSource {
	if v == nil {
		if oth == nil {
			return &VarSource{}
		}
		return oth
	}
	if oth == nil {
		return v
	}

	merged := v.YamlMapSource.Merge(oth.YamlMapSource, force)
	return &VarSource{YamlMapSource: merged}
}

var _ Mergeable[*VarSource] = (*VarSource)(nil)
