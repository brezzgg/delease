package models

type EnvSource struct {
	YamlMapSource[string]
}

func (e *EnvSource) Merge(oth *EnvSource, force bool) *EnvSource {
	if r := PremergeCheck(e, oth); r != nil {
		return r
	}
	if oth == nil {
		return e
	}

	merged := e.YamlMapSource.Merge(oth.YamlMapSource, force)
	return &EnvSource{YamlMapSource: merged}
}

var _ Mergeable[*EnvSource] = (*EnvSource)(nil)

func (e *EnvSource) StringSlice() []string {
	res := make([]string, len(e.data))
	for k, v := range e.data {
		res = append(res, k+"="+v)
	}
	return res
}
