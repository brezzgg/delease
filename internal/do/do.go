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
		if d.root.Do != nil && d.root.Do.Len() > 0 {
			for _, name := range d.root.Do.GetSource() {
				if _, ok := d.root.Tasks.Get(name); ok {
					tasks = append(tasks, name)
				}
			}
		}
		if len(tasks) == 0 {
			return lg.Ef("default task not set, use 'task.default: true'")
		}
	}

	osVars := GetOsVars(d.runtimeArgs)
	runtimeVars := osVars
	if d.runtimeVars != nil {
		runtimeVars = osVars.Merge(d.runtimeVars, true)
	}

	rootCtx := models.NewRootVarContext(d.root.Var, runtimeVars)
	compiler := NewTaskCompiler(rootCtx)

	taskLoader := NewTaskLoader(tasks, d.root)
	taskArr, err := taskLoader.Load()
	if err != nil {
		return err
	}
	if len(taskArr) == 0 {
		return lg.Ef("nothing to do")
	}

	environ := GetEnv()

	for _, task := range taskArr {
		compiledCmds, err := compiler.Compile(task)
		if err != nil {
			return lg.Ef("failed to compile task: %w", err)
		}

		env := environ.Merge(d.root.Env, true).Merge(task.Env, true)
		if d.runtimeEnvs != nil {
			env = env.Merge(d.runtimeEnvs, true)
		}

		wd, err := d.GetDir(task.Dir)
		if err != nil {
			return lg.Ef("bad working directory: %w", err)
		}

		var exe exec.Executor = &exec.Sh{}
		exe.Setup(wd, compiledCmds, env.StringSlice(), func(m string, t exec.MsgType) {
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
