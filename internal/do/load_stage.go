package do

import (
	"strings"

	"github.com/brezzgg/delease/internal/models"
	"github.com/brezzgg/go-packages/lg"
)

type LoadStage struct {
	tasks   []*Task
	root    *models.Root
	head    string
	headIdx int
}

func NewLoadStage(tasks []*Task, root *models.Root) *LoadStage {
	return &LoadStage{
		tasks: tasks,
		root:  root,
	}
}

func (t *LoadStage) Stage() ([]*Task, error) {
	loaded := make([]*Task, 0, len(t.tasks))
	for _, headTask := range t.tasks {
		t.head = headTask.Name

		res, err := t.load(headTask, true)
		if err != nil {
			return nil, err
		}

		for idx, task := range res {
			task.Head = idx == t.headIdx
			loaded = append(loaded, task)
		}
	}

	res := make([]*Task, 0, len(loaded))
	for _, load := range loaded {
		if load.Task != nil && load.Task.Cmds != nil && load.Task.Cmds.Len() > 0 {
			res = append(res, load)
		}
	}

	return res, nil
}

func (t *LoadStage) load(task *Task, recursive bool) ([]*Task, error) {
	if t.root.Tasks == nil {
		return nil, nil
	}

	var ok bool
	task.Task, ok = t.root.Tasks.Get(task.Name)
	if !ok {
		return nil, lg.Ef("task %s not found", task.Name)
	}

	parser := models.NewParamsParser()

	var after, before []*Task
	if recursive {
		if task.Task.Before != nil {
			tasks, err := t.loadCalls(parser, task.Task.Before.GetSource())
			if err != nil {
				return nil, lg.Ef("%s.before: %w", task.Name, err)
			}
			before = tasks
		}
		if task.Task.After != nil {
			tasks, err := t.loadCalls(parser, task.Task.After.GetSource())
			if err != nil {
				return nil, lg.Ef("%s.after: %w", task.Name, err)
			}
			after = tasks
		}
	}

	res := before
	res = append(res, task)
	res = append(res, after...)

	t.headIdx = len(before)

	return res, nil
}

func (t *LoadStage) loadCalls(parser *models.ParamsParser, tasks []string) ([]*Task, error) {
	var res []*Task
	for _, taskName := range tasks {
		var (
			tasks  []*Task
			name   string
			params *models.VarSource
			err    error
		)

		if strings.HasPrefix(taskName, "~") {
			taskName = taskName[1:]
			name, params, err = parser.Parse(taskName)
			if err != nil {
				return nil, lg.Ef("params parser: %w", err)
			}
		} else {
			if t.head == taskName {
				return nil, lg.Ef("recursive call in task %s detected", t.head)
			}
			name, params, err = parser.Parse(taskName)
			if err != nil {
				return nil, lg.Ef("params parser: %w", err)
			}
		}

		tasks, err = t.load(&Task{
			Name:   name,
			Params: params,
		}, false)
		if err != nil {
			return nil, err
		}

		res = append(res, tasks...)
	}

	return res, nil
}

var _ Stager = (*LoadStage)(nil)
