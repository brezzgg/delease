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
			for _, v := range current.Before.GetSource() {
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
			for _, v := range current.After.GetSource() {
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
