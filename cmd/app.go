package cmd

import (
	"github.com/brezzgg/go-packages/lg"
)

var (
	config string
	wd     string
)

func Run() {
	rootCmd.PersistentFlags().StringVarP(
		&config, "config", "C", "", "specify config file",
	)
	rootCmd.PersistentFlags().StringVarP(
		&wd, "working-directory", "D", "", "specify working directory",
	)

	configCmd.PersistentFlags().BoolVarP(
		&applie, "applie", "a", false, "apply variables",
	)
	configTaskCmd.PersistentFlags().BoolVarP(
		&printCmd, "print-cmd", "c", false, "print task commands",
	)
	configTasksCmd.PersistentFlags().BoolVarP(
		&printCmd, "print-cmd", "c", false, "print tasks commands",
	)
	configTasksCmd.PersistentFlags().BoolVarP(
		&taskNames, "names", "n", false, "print names",
	)

	rootCmd.AddCommand(
		uiCmd,
		configCmd,
	)
	configCmd.AddCommand(
		configParseCmd,
		configTaskCmd,
		configTasksCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		lg.Fatal(err.Error())
	}
}
