package do

import (
	"strings"

	"github.com/brezzgg/delease/internal/models"
	"github.com/brezzgg/go-packages/lg"
)

type TaskLoader struct {
	tasks []string
	root  *models.Root
	head  string
}

func NewTaskLoader(tasks []string, root *models.Root) *TaskLoader {
	return &TaskLoader{
		tasks: tasks,
		root:  root,
	}
}

func (t *TaskLoader) Load() ([]*models.Task, error) {
	if !t.root.Applied() {
		return nil, lg.Ef("root must be applied")
	}
	tasks := []*models.Task{}
	for _, task := range t.tasks {
		t.head = task
		r, err := t.load(task, true)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, r...)
		t.head = ""
	}
	res := []*models.Task{}
	for _, v := range tasks {
		if v != nil && v.Cmds != nil && v.Cmds.Len() > 0 {
			res = append(res, v)
		}
	}
	return res, nil
}

func (t *TaskLoader) load(task string, recursive bool) ([]*models.Task, error) {
	if t.root.Tasks == nil {
		return nil, nil
	}
	current, ok := t.root.Tasks.Get(task)
	if !ok {
		return nil, lg.Ef("task %s not found", task)
	}
	var after, before []*models.Task
	if recursive {
		if current.Before != nil {
			for _, v := range current.Before.Get() {
				var (
					tasks []*models.Task
					err   error
				)
				if strings.HasPrefix(v, "~") {
					tasks, err = t.load(v[1:], false)
				} else {
					if t.head == v {
						return nil, lg.Ef("recursive call in task %s detected", t.head)
					}
					tasks, err = t.load(v, true)
				}
				if err != nil {
					return nil, lg.Ef("%s.before: %w", v, err)
				}
				before = append(before, tasks...)
			}
		}
		if current.After != nil {
			for _, v := range current.After.Get() {
				var (
					tasks []*models.Task
					err   error
				)
				if strings.HasPrefix(v, "~") {
					tasks, err = t.load(v[1:], false)
				} else {
					if t.head == v {
						return nil, lg.Ef("recursive call in task %s detected", t.head)
					}
					tasks, err = t.load(v, true)
				}
				if err != nil {
					return nil, lg.Ef("%s.after: %w", v, err)
				}
				after = append(after, tasks...)
			}
		}
	}
	res := before
	res = append(res, current)
	res = append(res, after...)
	return res, nil
}

//
// func (d *Do) preload(name string) ([]*models.Task, error) {
// 	task, err := d.load(name)
// 	if err != nil {
// 		return nil, lg.Ef("task %s: %w", name, err)
// 	}
//
// 	beforeTasks, afterTasks, err := d.loadBeforeAfter(task)
// 	if err != nil {
// 		return nil, lg.Ef("task %s: %w", name, err)
// 	}
//
// 	res := beforeTasks
// 	res = append(res, task)
// 	res = append(res, afterTasks...)
//
// 	return res, nil
// }
//
// func (d *Do) loadBeforeAfter(task *models.Task) (before, after []*models.Task, err error) {
// 	if task.Before != nil {
// 		for _, v := range task.Before.Get() {
// 			t, e := d.load(v)
// 			if e != nil {
// 				return nil, nil, lg.Ef("failed to load before %s: %w", v, e)
// 			}
// 			before = append(before, t)
// 		}
// 	}
//
// 	if task.After != nil {
// 		for _, v := range task.After.Get() {
// 			t, e := d.load(v)
// 			if e != nil {
// 				return nil, nil, lg.Ef("failed to load after %s: %w", v, e)
// 			}
// 			after = append(after, t)
// 		}
// 	}
//
// 	return before, after, err
// }
//
// func (d *Do) load(name string) (*models.Task, error) {
// 	applied, err := d.root.ApplyVarsToTask(name, d.GetOsVars(d.args))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if applied.Tasks.Len() != 1 {
// 		return nil, lg.Ef("bad task len %d", applied.Tasks.Len())
// 	}
//
// 	for _, v := range applied.Tasks.GetMap() {
// 		return v, nil
// 	}
// 	return nil, lg.Ef("internal")
// }
