package models

import (
	"github.com/brezzgg/go-packages/lg"
	"gopkg.in/yaml.v3"
)

type Root struct {
	Include *IncludeSource `yaml:"includes" json:"includes,omitempty"`
	Var     *VarSource     `yaml:"vars" json:"vars,omitempty"`
	Tasks   *TaskSource    `yaml:"tasks" json:"tasks"`
	Env     *EnvSource     `yaml:"envs" json:"envs,omitempty"`

	applied bool
}

func (r *Root) UnmarshalYAML(value *yaml.Node) error {
	type rootAlias Root
	var root rootAlias

	if err := value.Decode(&root); err != nil {
		return err
	}

	root.Tasks.Clean()
	if root.Var == nil {
		root.Var = &VarSource{}
	}
	if root.Env == nil {
		root.Env = &EnvSource{}
	}
	if root.Include == nil {
		root.Include = &IncludeSource{}
	}

	r.Tasks = root.Tasks
	r.Var = root.Var
	r.Env = root.Env
	r.Include = root.Include

	return nil
}

func (r *Root) ApplyVars(vars *VarSource) (*Root, error) {
	if r == nil {
		return nil, lg.Ef("root is nil")
	}
	if r.Tasks == nil {
		return &Root{
			Var:     r.Var,
			Tasks:   r.Tasks,
			Env:     r.Env,
			applied: true,
		}, nil
	}
	global := r.Var
	if vars != nil {
		global = global.Merge(vars, true)
	}
	tasks := make(map[string]*Task, r.Tasks.Len())
	for k, v := range r.Tasks.GetMapCopy() {
		applied, err := v.ApplyVars(global)
		if err != nil {
			return nil, lg.Ef("task %s: %w", k, err)
		}
		tasks[k] = applied
	}
	taskSrc := &TaskSource{}
	taskSrc.SetSource(tasks)
	return &Root{
		Var:     r.Var,
		Tasks:   taskSrc,
		Env:     r.Env,
		applied: true,
	}, nil
}

func (r *Root) Applied() bool {
	return r.applied
}

func (r *Root) Merge(oth *Root, force bool) *Root {
	if res := PremergeCheck(r, oth); res != nil {
		return res
	}

	if r.applied != oth.applied {
		panic("applied root merges with not applied")
	}

	res := &Root{}
	res.Var = r.Var.Merge(oth.Var, force)
	res.Env = r.Env.Merge(oth.Env, force)
	res.Tasks = r.Tasks.Merge(oth.Tasks, force)
	res.Include = r.Include.Merge(oth.Include, force)

	return res
}

var (
	_ yaml.Unmarshaler = (*Root)(nil)
	_ Applier[*Root]   = (*Root)(nil)
	_ Mergeable[*Root] = (*Root)(nil)
)

type IncludeSource struct {
	YamlSliceSource[string]
}

func (i *IncludeSource) Merge(oth *IncludeSource, force bool) *IncludeSource {
	if r := PremergeCheck(i, oth); r != nil {
		return r
	}
	if force {
		return i.sort(i, oth)
	}
	return i.sort(oth, i)
}

func (i *IncludeSource) sort(left, right *IncludeSource) *IncludeSource {
	rIndex := make(map[string]int, right.Len())
	pre := right.GetCopy()

	for i, v := range pre {
		rIndex[v] = i
	}

	for _, v := range left.GetCopy() {
		idx, ok := rIndex[v]
		if ok {
			pre[idx] = ""
		}
		pre = append(pre, v)
	}

	res := []string{}
	for _, v := range pre {
		if v != "" {
			res = append(res, v)
		}
	}

	src := &IncludeSource{}
	src.SetSource(res)
	return src
}

var _ Mergeable[*IncludeSource] = (*IncludeSource)(nil)
