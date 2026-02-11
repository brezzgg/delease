package cmd

import (
	"slices"
	"strings"

	"github.com/brezzgg/delease/internal/models"
	"github.com/brezzgg/delease/internal/parser"
	"github.com/brezzgg/go-packages/lg"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		parser := parser.New(config, wd)
		root, err := parser.Parse()
		if err != nil {
			lg.Fatal(ErrParseFailed(err))
		}

		if root.Tasks == nil || root.Tasks.Len() == 0 {
			return
		}

		type task struct {
			name string
			desc string
		}

		doTasks := make(map[string]struct{})
		if root.Do != nil && root.Do.Len() > 0 {
			params := models.NewParamsParser()
			for _, taskName := range root.Do.GetSource() {
				name, _, err := params.Parse(taskName)
				if err != nil {
					continue
				}
				doTasks[name] = struct{}{}
			}
		}

		var descTasks, otherTasks []task
		maxNameLen := 0

		for name, ta := range root.Tasks.GetSource() {
			if _, isDo := doTasks[name]; isDo {
				name = "*" + name
			}
			t := task{
				name: name,
				desc: ta.Desc,
			}

			nameLen := len(name)
			if nameLen > maxNameLen {
				maxNameLen = nameLen
			}

			if ta.Desc != "" {
				descTasks = append(descTasks, t)
			} else {
				otherTasks = append(otherTasks, t)
			}
		}

		slices.SortFunc(descTasks, func(a, b task) int {
			return strings.Compare(a.name, b.name)
		})
		slices.SortFunc(otherTasks, func(a, b task) int {
			return strings.Compare(a.name, b.name)
		})

		for _, t := range descTasks {
			name := t.name
			lg.Info(lg.F("%-*s%s", maxNameLen+4, name, t.desc))
		}

		for _, t := range otherTasks {
			name := t.name
			lg.Info(name)
		}
	},
}
