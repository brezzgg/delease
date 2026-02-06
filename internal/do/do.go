package do

import (
	"context"
	"fmt"
	"strings"

	"github.com/brezzgg/delease/internal/exec"
	"github.com/brezzgg/delease/internal/models"
	"github.com/brezzgg/go-packages/lg"
)

type Do struct {
	root *models.Root

	runtimeArgs string
	runtimeEnvs *models.EnvSource
	runtimeVars *models.VarSource
}

type DoOption func(*Do)

func New(root *models.Root, opts ...DoOption) *Do {
	d := &Do{
		root: root,
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

func WithArgs(args string) DoOption {
	return func(d *Do) {
		d.runtimeArgs = args
	}
}

func WithEnvs(envs []string) DoOption {
	return func(d *Do) {
		if len(envs) == 0 {
			return 
		}
		data := make(map[string]string, len(envs))
		for _, v := range envs {
			if strings.Contains(v, "=") {
				before, after, ok := strings.Cut(v, "=")
				if ok {
					data[before] = after
				}
			}
		}
		src := &models.EnvSource{}
		src.SetSource(data)
		d.runtimeEnvs = src
	}
}

func WithVars(vars []string) DoOption {
	return func(d *Do) {
		if len(vars) == 0 {
			return 
		}
		data := make(map[string]string, len(vars))
		for _, v := range vars {
			if strings.Contains(v, "=") {
				before, after, ok := strings.Cut(v, "=")
				if ok {
					data[before] = after
				}
			}
		}
		src := &models.VarSource{}
		src.SetSource(data)
		d.runtimeVars = src
	}
}

func (d *Do) Execute(ctx context.Context, tasks []string) error {
	if len(tasks) == 0 {
		for name, task := range d.root.Tasks.GetMap() {
			if task.Default {
				tasks = append(tasks, name)
				continue
			}
		}
		if len(tasks) == 0 {
			return lg.Ef("default task not set, use 'task.default: true'")
		}
	}

	if d.root.Applied() {
		return lg.Ef("root cant be applied")
	}

	vars := d.GetOsVars(d.runtimeArgs)
	if d.runtimeVars != nil {
		vars = vars.Merge(d.runtimeVars, true)
	}

	newRoot, err := d.root.ApplyVars(vars)
	if err != nil {
		return lg.Ef("faild to apply vars: %w", err)
	}
	d.root = newRoot

	taskLoader := NewTaskLoader(tasks, d.root)
	taskArr, err := taskLoader.Load()
	if err != nil {
		return err
	}
	if len(taskArr) == 0 {
		return lg.Ef("nothing to do")
	}

	environ := d.GetEnv()

	for _, task := range taskArr {
		env := environ
		env = env.Merge(d.root.Env, true)
		env = env.Merge(task.Env, true)
		if d.runtimeEnvs != nil {
			env = env.Merge(d.runtimeEnvs, true)
		}

		cmds := task.Cmds.Get()
		var lines []string
		for _, v := range cmds {
			lines = append(lines, v.Cmd)
		}

		wd, err := d.GetDir(task.Dir)
		if err != nil {
			return lg.Ef("bad working directory: %w", err)
		}

		var exe exec.Executor = &exec.Sh{}
		exe.Setup(wd, lines, env.StringSlice(), func(m string, t exec.MsgType) {
			if t == exec.MsgTypeStdout {
				fmt.Print(m)
			} else {
				fmt.Print(m)
			}
		})

		ch := make(chan exec.Result, 1)
		exe.Run(ctx, ch)
		res := <-ch
		if res.Error != nil {
			return lg.Ef("exec error: %s", res.Error)
		}
		if res.Code != 0 {
			return lg.Ef("app exit with status code %d", res.Code)
		}
	}

	return nil
}
