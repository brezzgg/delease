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
		data := make(map[string]*models.Var, len(vars))
		for _, v := range vars {
			if strings.Contains(v, "=") {
				before, after, ok := strings.Cut(v, "=")
				if ok {
					data[before] = models.NewVar(after, models.VarTypeStatic)
				}
			}
		}
		src := &models.VarSource{}
		src.SetSource(data)
		d.runtimeVars = src
	}
}

func (d *Do) Execute(ctx context.Context, taskNames []string) error {
	var (
		tasks = NewTasks(taskNames)
		err   error
	)

	paramParser := NewParamsStage(d.root, tasks)
	tasks, err = paramParser.Stage()
	if err != nil {
		return lg.Ef("params stage error: %w", err)
	}

	osVars := GetOsVars(d.runtimeArgs)
	runtimeVars := osVars
	if d.runtimeVars != nil {
		runtimeVars = osVars.Merge(d.runtimeVars, true)
	}

	rootCtx := models.NewVarContext(d.root.Var, runtimeVars)
	compiler := NewTaskCompiler(rootCtx)

	taskLoader := NewLoadStage(tasks, d.root)
	tasks, err = taskLoader.Stage()
	if err != nil {
		return lg.Ef("load stage error: %w", err)
	}
	if len(tasks) == 0 {
		return lg.Ef("nothing to do")
	}

	environ := GetEnv()

	for _, task := range tasks {
		var cmds []string
		cmds, err = compiler.Compile(task)
		if err != nil {
			return lg.Ef("compile error: %w", err)
		}

		env := environ.Merge(d.root.Env, true).Merge(task.Task.Env, true)
		if d.runtimeEnvs != nil {
			env = env.Merge(d.runtimeEnvs, true)
		}

		wd, err := d.GetDir(task.Task.Dir)
		if err != nil {
			return lg.Ef("bad working directory: %w", err)
		}

		var exe exec.Executor = &exec.Sh{}
		exe.Setup(wd, cmds, env.StringSlice(), func(m string, t exec.MsgType) {
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
