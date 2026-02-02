package models

type Task struct {
	Before  []string `yaml:"before,omitempty"`
	After   []string `yaml:"after,omitempty"`
	Cmds    []string `yaml:"cmds,omitempty"`
	Vars    Vars     `yaml:"vars,omitempty"`
	Dir     string   `yaml:"dir"`
	Env     Envs     `yaml:"env"`
	Default bool     `yaml:"default"`
}

func (r *Task) GetVars() Vars {
	return r.Vars
}

type Envs map[string]string
