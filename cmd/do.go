package cmd

import (
	"context"

	"github.com/brezzgg/delease/internal/do"
	"github.com/brezzgg/delease/internal/parser"
	"github.com/brezzgg/go-packages/lg"
	"github.com/spf13/cobra"
)

var (
	doArgs string
	doVars []string
	doEnvs []string
)

var (
	doCmd = &cobra.Command{
		Use: "do [tasks...]",
		Short: "Do tasks",
		Args: cobra.RangeArgs(0, 128),
		Run: func(cmd *cobra.Command, args []string) {
			parser := parser.New(config, wd)
			root, err := parser.Parse()
			if err != nil {
				lg.Fatal(ErrParseFailed(err))
			}

			d := do.New(
				root,
				do.WithArgs(doArgs),
				do.WithEnvs(doEnvs),
				do.WithVars(doVars),
			)
			if err := d.Execute(context.Background(), args); err != nil {
				lg.Fatal(err)
			}
		},
	}
)
