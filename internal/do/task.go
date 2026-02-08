package do

import (
	"strings"

	"github.com/brezzgg/delease/internal/models"
)

type Stager interface {
	Stage() ([]*Task, error)
}

type Task struct {
	Name   string
	Task   *models.Task
	Params *models.VarSource
	Head   bool
}

func NewTasks(tasks []string) []*Task {
	tasks = normalizeTasks(tasks)
	res := make([]*Task, 0, len(tasks))
	for _, name := range tasks {
		res = append(res, &Task{Name: name})
	}
	return res
}

func normalizeTasks(args []string) []string {
	if len(args) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(args))
	var buffer strings.Builder
	inQuotes := false

	for i := range args {
		arg := args[i]
		quoteCount := strings.Count(arg, `"`)

		if !inQuotes {
			if quoteCount%2 == 1 {
				buffer.WriteString(arg)
				inQuotes = true
			} else {
				result = append(result, arg)
			}
		} else {
			buffer.WriteString(" ")
			buffer.WriteString(arg)

			if quoteCount%2 == 1 {
				result = append(result, buffer.String())
				buffer.Reset()
				inQuotes = false
			}
		}
	}

	if buffer.Len() > 0 {
		result = append(result, buffer.String())
	}

	return result
}
