package cmd

import (
	"github.com/brezzgg/delease/gui"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "delease",
	}

	uiCmd = &cobra.Command{
		Use:   "ui",
		Short: "Open ui",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	configCmd = &cobra.Command{
		Use:     "config",
		Aliases: []string{"conf", "cfg"},
		Short: "Manage configuration",
	}
)
