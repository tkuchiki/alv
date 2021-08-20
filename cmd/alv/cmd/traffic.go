package cmd

import (
	"github.com/tkuchiki/alv/visualiser"

	"github.com/spf13/cobra"
)

func NewTrafficCmd(rootCmd *cobra.Command) *cobra.Command {
	var trafficCmd = &cobra.Command{
		Use:   "traffic",
		Short: "Visualize the traffic of access logs",
		Long:  `Visualize the traffic of access logs`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := getOptions(rootCmd, cmd)
			if err != nil {
				return err
			}

			v := visualiser.NewVisualizer(opts.LogFormat)
			err = v.SetInReader(opts.File)
			if err != nil {
				return err
			}
			defer v.CloseInReader()

			v.SetParser(opts)

			err = v.SetOutWriter(opts.Output)
			if err != nil {
				return err
			}
			defer v.CloseOutWriter()

			if len(opts.MatchingGroups) > 0 {
				err = v.SetUriMatchingGroups(opts.MatchingGroups)
				if err != nil {
					return err
				}
			}

			return v.Render()
		},
	}

	trafficCmd.PersistentFlags().StringP("uri-key", "", "uri", "Change the uri key")
	trafficCmd.PersistentFlags().StringP("method-key", "", "method", "Change the method key")
	trafficCmd.PersistentFlags().StringP("user-key", "", "user", "Change the user key")

	return trafficCmd
}
