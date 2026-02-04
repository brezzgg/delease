package cmd

import (
	"context"

	"github.com/brezzgg/delease/internal/do"
	"github.com/brezzgg/delease/internal/parser"
	"github.com/brezzgg/go-packages/lg"
	"github.com/spf13/cobra"
)

var (
	doCmd = &cobra.Command{
		Use: "do [tasks...]",
		Short: "Do tasks",
		Args: cobra.RangeArgs(0, 128),
		Run: func(cmd *cobra.Command, args []string) {
			b, err := parser.FindConfig(config, wd)
			if err != nil {
				lg.Fatal(ErrBadConfig(err))
			}
			root, err := parser.Parse(b)
			if err != nil {
				lg.Fatal(ErrParseFailed(err))
			}

			d := do.New(root)
			if err := d.Execute(context.Background(), args); err != nil {
				lg.Fatal(err)
			}
		},
	}
)
