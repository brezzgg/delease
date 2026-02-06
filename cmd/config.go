package cmd

import (
	"github.com/brezzgg/delease/internal/parser"
	"github.com/brezzgg/go-packages/lg"
	"github.com/spf13/cobra"
)

var (
	printCmd  bool
	applie    bool
	taskNames bool
)

var (
	configParseCmd = &cobra.Command{
		Use:   "parse",
		Short: "Show all parsed data",
		Run: func(cmd *cobra.Command, args []string) {
			parser := parser.New(config, wd)
			root, err := parser.Parse()
			if err != nil {
				lg.Fatal(ErrParseFailed(err))
			}

			// if -a: apply vars
			if applie {
				if r, err := root.ApplyVars(nil); err != nil {
					lg.Fatal(ErrApplyVars(err))
				} else {
					root = r
				}
			}

			lg.Info("parse successful", lg.C{"root": root})
		},
	}

	configTaskCmd = &cobra.Command{
		Use:   "task [task name]",
		Short: "Print task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			parser := parser.New(config, wd)
			root, err := parser.Parse()
			if err != nil {
				lg.Fatal(ErrParseFailed(err))
			}
			if root.Tasks.Len() == 0 {
				return
			}

			// if -a: apply vars
			if applie {
				if r, err := root.ApplyVars(nil); err != nil {
					lg.Fatal(ErrApplyVars(err))
				} else {
					root = r
				}
			}

			// show task
			task, ok := root.Tasks.Get(args[0])
			if !ok {
				lg.Fatal(ErrTaskNotFound(args[0]))
			}
			lg.Info("ok", lg.C{args[0]: task})
		},
	}

	configTasksCmd = &cobra.Command{
		Use:   "tasks",
		Short: "Print all tasks",
		Run: func(cmd *cobra.Command, args []string) {
			parser := parser.New(config, wd)
			root, err := parser.Parse()
			if err != nil {
				lg.Fatal(ErrParseFailed(err))
			}
			if root.Tasks.Len() == 0 {
				return
			}

			// if -a: apply vars
			if applie {
				if r, err := root.ApplyVars(nil); err != nil {
					lg.Fatal(ErrApplyVars(err))
				} else {
					root = r
				}
			}

			switch {
			case taskNames:
				// print only names
				keys := root.Tasks.Keys()
				lg.Info("ok", lg.C{"n": len(keys), "tasks": keys})
			case printCmd:
				// print only cmds
				tasks := root.Tasks.GetMap()
				c := make(lg.C, len(tasks))
				for key, val := range tasks {
					c[key] = val.Cmds
				}
				lg.Info("ok", c)
			default:
				// print full
				lg.Info("ok", lg.C{"n": root.Tasks.Len(), "tasks": root.Tasks.GetMap()})
			}
		},
	}
)
