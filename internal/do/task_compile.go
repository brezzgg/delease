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

func (c *TaskCompiler) Compile(task *models.Task) ([]string, error) {
	if task == nil {
		return nil, lg.Ef("task is nil")
	}

	if task.Cmds == nil || task.Cmds.Len() == 0 {
		return []string{}, nil
	}

	taskCtx := c.rootCtx.Child(task.Vars)

	return task.Cmds.Compile(taskCtx)
}
