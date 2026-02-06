package models

type TaskSource struct {
	YamlMapSource[*Task]
	applied bool
}

func (t *TaskSource) ApplyVars(vars *VarSource) (*TaskSource, error) {
	if t == nil {
		return nil, nil
	}

	newData := make(map[string]*Task, t.Len())

	for _, key := range t.Keys() {
		task, ok := t.Get(key)
		if !ok {
			continue
		}
		newTask, err := task.ApplyVars(vars)
		if err != nil {
			return nil, err
		}
		newData[key] = newTask
	}

	res := &TaskSource{
		applied: true,
	}
	res.SetSource(newData)
	return res, nil
}

func (t *TaskSource) Applied() bool {
	return t.applied
}

func (t *TaskSource) Clean() {
	if t == nil {
		return
	}
	MapClean(t.GetMap())
}

func (t *TaskSource) Merge(oth *TaskSource, force bool) *TaskSource {
	if r := PremergeCheck(t, oth); r != nil {
		return r
	}

	merged := t.YamlMapSource.Merge(oth.YamlMapSource, force)
	return &TaskSource{YamlMapSource: merged}
}

var (
	_ Applier[*TaskSource]   = (*TaskSource)(nil)
	_ Mergeable[*TaskSource] = (*TaskSource)(nil)
)

type Task struct {
	Before  *TaskCallSource `yaml:"before" json:"before,omitempty"`
	After   *TaskCallSource `yaml:"after" json:"after,omitempty"`
	Cmds    *CmdSource      `yaml:"cmds" json:"cmds"`
	Dir     string          `yaml:"dir" json:"dir,omitempty"`
	Vars    *VarSource      `yaml:"vars" json:"vars,omitempty"`
	Env     *EnvSource      `yaml:"envs" json:"envs,omitempty"`
	Default bool            `yaml:"default" json:"default,omitempty"`

	appled bool
}

func (t *Task) ApplyVars(vars *VarSource) (*Task, error) {
	if t == nil {
		return nil, nil
	}
	if vars == nil {
		vars = &VarSource{}
	}

	localVars := t.Vars
	if localVars == nil {
		localVars = &VarSource{}
	}

	merged := vars.Merge(localVars, true)
	newCmds, err := t.Cmds.ApplyVars(merged)
	if err != nil {
		return nil, err
	}

	return &Task{
		Before:  t.Before,
		After:   t.After,
		Cmds:    newCmds,
		Dir:     t.Dir,
		Vars:    t.Vars,
		Env:     t.Env,
		Default: t.Default,
		appled:  true,
	}, nil
}

func (t *Task) Applied() bool {
	return t.appled
}

var _ Applier[*Task] = (*Task)(nil)

type TaskCallSource struct {
	YamlSliceSource[string]
}
