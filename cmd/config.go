package cmd

import (
	"github.com/brezzgg/delease/internal/models"
	"github.com/brezzgg/delease/internal/parser"
	"github.com/brezzgg/go-packages/lg"
	"github.com/spf13/cobra"
)

var (
	printCmd  bool
	compile   bool
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

			var compileC lg.C
			// if -a: compile vars
			if compile {
				compileC = make(lg.C)
				ctx := models.NewExampleVarContext(root.Var)

				for name, task := range root.Tasks.GetSource() {
					taskCtx := ctx.Child(task.Vars)
					compiled, err := task.Cmds.Compile(taskCtx)
					if err != nil {
						lg.Fatal(ErrCompileVars(err))
					}
					compileC[name] = compiled
				}
			}

			if compileC == nil {
				compileC = lg.C{"help": "use `-a` to see compiled commands"}
			}
			lg.Info("parse successful", lg.C{"root": root}, lg.C{"compiled_cmds": compileC})
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

			task, ok := root.Tasks.Get(args[0])
			if !ok {
				lg.Fatal(ErrTaskNotFound(args[0]))
			}

			// if -a: compile vars
			if compile {
				ctx := models.NewExampleVarContext(root.Var)
				taskCtx := ctx.Child(task.Vars)
				compiled, err := task.Cmds.Compile(taskCtx)
				if err != nil {
					lg.Fatal(ErrCompileVars(err))
				}
				lg.Info("ok", lg.C{
					"task":          task,
					"compiled_cmds": compiled,
				})
			} else {
				lg.Info("ok", lg.C{args[0]: task})
			}
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

			switch {
			case taskNames:
				keys := root.Tasks.Keys()
				lg.Info("ok", lg.C{"n": len(keys), "tasks": keys})
			case printCmd:
				c := make(lg.C, root.Tasks.Len())
				if compile {
					ctx := models.NewExampleVarContext(root.Var)
					for key, task := range root.Tasks.GetSource() {
						taskCtx := ctx.Child(task.Vars)
						compiled, err := task.Cmds.Compile(taskCtx)
						if err != nil {
							lg.Fatal(ErrCompileVars(err))
						}
						c[key] = compiled
					}
				} else {
					for key, task := range root.Tasks.GetSource() {
						c[key] = task.Cmds
					}
				}
				lg.Info("ok", c)
			default:
				lg.Info("ok", lg.C{"n": root.Tasks.Len(), "tasks": root.Tasks.GetSource()})
			}
		},
	}
)
