package cmd

import (
	"github.com/brezzgg/go-packages/lg"
)

var (
	config string
	wd     string
)

func Run() {
	rootCmd.PersistentFlags().BoolVarP(
		&verbose, "verbose", "v", false, "verbose output",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&indent, "indent", "i", false, "indent output",
	)
	rootCmd.PersistentFlags().StringVarP(
		&config, "config", "C", "", "specify config file",
	)
	rootCmd.PersistentFlags().StringVarP(
		&wd, "working-directory", "D", "", "specify working directory",
	)

	configCmd.PersistentFlags().BoolVarP(
		&compile, "compile", "a", false, "compile variables",
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

	doCmd.PersistentFlags().StringVarP(
		&doArgs, "args", "A", "", "set args that can be used as {{ os.args }}",
	)
	doCmd.PersistentFlags().StringArrayVarP(
		&doVars, "vars", "V", nil, "set vars that can be used as {{ var-name }}",
	)
	doCmd.PersistentFlags().StringArrayVarP(
		&doEnvs, "envs", "E", nil, "set envs that can be used as $env-name",
	)

	rootCmd.AddCommand(
		uiCmd,
		configCmd,
		doCmd,
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
