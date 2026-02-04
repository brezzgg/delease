package do

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/brezzgg/delease/internal/exec"
	"github.com/brezzgg/delease/internal/models"
	"github.com/brezzgg/go-packages/lg"
)

type Do struct {
	root *models.Root
	args string
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
		d.args = args
	}
}

func (d *Do) preload(name string) ([]*models.Task, error) {
	task, err := d.load(name)
	if err != nil {
		return nil, lg.Ef("task %s: %w", name, err)
	}

	beforeTasks, afterTasks, err := d.loadBeforeAfter(task)
	if err != nil {
		return nil, lg.Ef("task %s: %w", name, err)
	}

	res := beforeTasks
	res = append(res, task)
	res = append(res, afterTasks...)

	return res, nil
}

func (d *Do) loadBeforeAfter(task *models.Task) (before, after []*models.Task, err error) {
	if task.Before != nil {
		for _, v := range task.Before.Get() {
			t, e := d.load(v)
			if e != nil {
				return nil, nil, lg.Ef("failed to load before %s: %w", v, e)
			}
			before = append(before, t)
		}
	}

	if task.After != nil {
		for _, v := range task.After.Get() {
			t, e := d.load(v)
			if e != nil {
				return nil, nil, lg.Ef("failed to load after %s: %w", v, e)
			}
			after = append(after, t)
		}
	}

	return before, after, err
}

func (d *Do) load(name string) (*models.Task, error) {
	applied, err := d.root.ApplyVarsToTask(name, d.GetOsVars(d.args))
	if err != nil {
		return nil, err
	}
	if applied.Tasks.Len() != 1 {
		return nil, lg.Ef("bad task len %d", applied.Tasks.Len())
	}

	for _, v := range applied.Tasks.GetMap() {
		return v, nil
	}
	return nil, lg.Ef("internal")
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

	taskMatrix := make([][]*models.Task, 0, len(tasks))

	for _, name := range tasks {
		t, err := d.preload(name)
		if err != nil {
			return err
		}
		taskMatrix = append(taskMatrix, t)
	}

	environ := d.environ()

	for _, row := range taskMatrix {
		for _, task := range row {
			env := environ
			env = env.Merge(d.root.Env, true)
			env = env.Merge(task.Env, true)

			lines := task.Cmds.Get()

			var (
				wd  = task.Dir
				err error
			)
			if wd == "" {
				wd, err = os.Getwd()
				if err != nil {
					return lg.Ef("bad wd: %s", err)
				}
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
	}
	return nil
}

func (d *Do) environ() *models.EnvSource {
	envs := os.Environ()
	m := make(map[string]string, 10)
	for _, v := range envs {
		if before, after, ok := strings.Cut(v, "="); ok {
			m[before] = after
		}
	}
	r := &models.EnvSource{}
	r.SetSource(m)
	return r
}
