package do

import (
	"github.com/brezzgg/delease/internal/models"
	"github.com/brezzgg/go-packages/lg"
)

type TaskCompiler struct {
	rootCtx *models.VarContext
}

func NewTaskCompiler(rootCtx *models.VarContext) *TaskCompiler {
	return &TaskCompiler{rootCtx: rootCtx}
}

func (c *TaskCompiler) Compile(task *Task) ([]string, error) {
	if task == nil || task.Task == nil {
		return nil, lg.Ef("task is nil")
	}

	if task.Task.Cmds == nil || task.Task.Cmds.Len() == 0 {
		return []string{}, nil
	}

	vars := task.Task.Vars.Merge(task.Params, true)
	taskCtx := c.rootCtx.Child(vars)

	return task.Task.Cmds.Compile(taskCtx)
}
