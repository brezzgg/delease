package do

import (
	"strings"

	"github.com/brezzgg/delease/internal/models"
	"github.com/brezzgg/go-packages/lg"
)

type ParamsStage struct {
	tasks []*Task
	root  *models.Root
}

func NewParamsStage(root *models.Root, tasks []*Task) *ParamsStage {
	return &ParamsStage{
		tasks: tasks,
		root:  root,
	}
}

func (p *ParamsStage) Stage() ([]*Task, error) {
	if len(p.tasks) == 0 {
		if p.root.Do != nil && p.root.Do.Len() > 0 {
			for _, name := range p.root.Do.GetSource() {
				if _, ok := p.root.Tasks.Get(name); ok {
					p.tasks = append(p.tasks, &Task{Name: name})
				}
			}
		}
		if len(p.tasks) == 0 {
			return nil, lg.Ef("default tasks not set, use 'do: [task1, task2, ..., taskn]'")
		}
	}

	parser := models.NewParamsParser()
	res := make([]*Task, 0, len(p.tasks))
	for _, task := range p.tasks {
		if len(strings.TrimSpace(task.Name)) == 0 {
			continue
		}

		if strings.Contains(task.Name, ":") {
			var err error
			if err = p.parse(parser, task); err != nil {
				return nil, lg.Ef("param %s: %w", task.Name, err)
			}
		}
		res = append(res, task)
	}
	return res, nil
}

func (p *ParamsStage) parse(parser *models.ParamsParser, task *Task) error {
	name, params, err := parser.Parse(task.Name)
	if err != nil {
		return err
	}
	task.Name = name
	task.Params = params
	return nil
}

var _ Stager = (*ParamsStage)(nil)
