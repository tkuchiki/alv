package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:     "alv",
		Short:   "Access Log Visualizer",
		Long:    `Access Log Visualizer`,
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringP("file", "", "", "The access log file")
	rootCmd.PersistentFlags().StringP("out", "", "", "The sankey diagram")
	rootCmd.PersistentFlags().StringP("log-format", "l", "json", "The log format")
	rootCmd.PersistentFlags().BoolP("query-string", "q", false, "Include the URI query string")
	rootCmd.PersistentFlags().BoolP("qs-ignore-values", "", false, " Ignore the value of the query string. Replace all values with xxx (only use with -q)")
	rootCmd.PersistentFlags().StringSliceP("matching-group", "m", []string{}, "Specifies URI matching group")
	rootCmd.PersistentFlags().StringSliceP("filter", "f", []string{}, " Only the logs are profiled that match the conditions")

	rootCmd.AddCommand(NewTrafficCmd(rootCmd))
	rootCmd.SetVersionTemplate(fmt.Sprintln(version))

	return rootCmd
}

func Execute() error {
	rootCmd := NewRootCmd()
	return rootCmd.Execute()
}
