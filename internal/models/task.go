package models

type TaskSource struct {
	YamlMapSource[*Task]
}

func (t *TaskSource) Clean() {
	if t == nil {
		return
	}
	MapClean(t.GetSource())
}

func (t *TaskSource) Merge(oth *TaskSource, force bool) *TaskSource {
	if r := PremergeCheck(t, oth); r != nil {
		return r
	}

	merged := t.YamlMapSource.Merge(oth.YamlMapSource, force)
	return &TaskSource{YamlMapSource: merged}
}

var _ Mergeable[*TaskSource] = (*TaskSource)(nil)

type Task struct {
	Before  *TaskCallSource `yaml:"before" json:"before,omitempty"`
	After   *TaskCallSource `yaml:"after" json:"after,omitempty"`
	Cmds    *CmdSource      `yaml:"cmds" json:"cmds"`
	Dir     string          `yaml:"dir" json:"dir,omitempty"`
	Vars    *VarSource      `yaml:"vars" json:"vars,omitempty"`
	Env     *EnvSource      `yaml:"envs" json:"envs,omitempty"`
	Default bool            `yaml:"default" json:"default,omitempty"`
}

type TaskCallSource struct {
	YamlSliceSource[string]
}
