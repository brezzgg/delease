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
			b, err := parser.FindConfig(config, wd)
			if err != nil {
				lg.Fatal(ErrBadConfig(err))
			}
			root, err := parser.Parse(b)
			if err != nil {
				lg.Fatal(ErrParseFailed(err))
			}

			// if -a: apply vars
			if applie {
				for k := range root.Tasks {
					err = parser.ApplyVars(root, k, nil)
					if err != nil {
						lg.Fatal(ErrApplyVars(err))
					}
				}
			}

			lg.Info("parse successful", lg.C{"root": root})
		},
	}

	configTaskCmd = &cobra.Command{
		Use:   "task [task name]",
		Short: "Print task",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			b, err := parser.FindConfig(config, wd)
			if err != nil {
				lg.Fatal(ErrBadConfig(err))
			}
			root, err := parser.Parse(b)
			if err != nil {
				lg.Fatal(ErrParseFailed(err))
			}
			if len(root.Tasks) == 0 {
				return
			}

			// if -a: apply vars
			if applie {
				for k := range root.Tasks {
					err = parser.ApplyVars(root, k, nil)
					if err != nil {
						lg.Fatal(ErrApplyVars(err))
					}
				}
			}

			// show task
			task, err := root.GetTask(args[0])
			if err != nil {
				lg.Fatal(ErrTaskNotFound(err))
			}
			lg.Info("ok", lg.C{args[0]: task})
		},
	}

	configTasksCmd = &cobra.Command{
		Use:   "tasks",
		Short: "Print all tasks",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := parser.FindConfig(config, wd)
			if err != nil {
				lg.Fatal(ErrBadConfig(err))
			}
			root, err := parser.Parse(b)
			if err != nil {
				lg.Fatal(ErrParseFailed(err))
			}
			if len(root.Tasks) == 0 {
				return
			}

			// if -a: apply vars
			if applie {
				for k := range root.Tasks {
					err := parser.ApplyVars(root, k, nil)
					if err != nil {
						lg.Fatal(ErrApplyVars(err))
					}
				}
			}

			switch {
			case taskNames:
				// print only names
				o := []string{}
				for key := range root.Tasks {
					o = append(o, key)
				}
				lg.Info("ok", lg.C{"n": len(o), "tasks": o})
			case printCmd:
				// print only cmds
				c := make(lg.C, len(root.Tasks))
				for key, val := range root.Tasks {
					c[key] = val.Cmds
				}
				lg.Info("ok", c)
			default:
				// print full
				lg.Info("ok", lg.C{"n": len(root.Tasks), "tasks": root.Tasks})
			}
		},
	}
)
