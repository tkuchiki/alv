package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tkuchiki/alv/cmd/alv/option"
)

func getOptions(rootCmd, cmd *cobra.Command) (*option.Options, error) {
	file, err := rootCmd.PersistentFlags().GetString("file")
	if err != nil {
		return nil, err
	}

	out, err := rootCmd.PersistentFlags().GetString("out")
	if err != nil {
		return nil, err
	}

	logFormat, err := rootCmd.PersistentFlags().GetString("log-format")
	if err != nil {
		return nil, err
	}

	queryString, err := rootCmd.PersistentFlags().GetBool("query-string")
	if err != nil {
		return nil, err
	}

	queryStringIgnoreValues, err := rootCmd.PersistentFlags().GetBool("qs-ignore-values")
	if err != nil {
		return nil, err
	}

	matchingGroups, err := rootCmd.PersistentFlags().GetStringSlice("matching-group")
	if err != nil {
		return nil, err
	}

	filters, err := rootCmd.PersistentFlags().GetStringSlice("filter")
	if err != nil {
		return nil, err
	}

	var keyOpt option.KeyOption
	switch logFormat {
	case "json":
		uriKey, err := cmd.Flags().GetString("uri-key")
		if err != nil {
			return nil, err
		}

		methodKey, err := cmd.Flags().GetString("method-key")
		if err != nil {
			return nil, err
		}

		userKey, err := cmd.Flags().GetString("user-key")
		if err != nil {
			return nil, err
		}

		keyOpt = option.KeyOption{
			UriKey:    uriKey,
			MethodKey: methodKey,
			UserKey:   userKey,
		}

	}

	return &option.Options{
		File:                    file,
		Output:                  out,
		LogFormat:               logFormat,
		QueryString:             queryString,
		QueryStringIgnoreValues: queryStringIgnoreValues,
		MatchingGroups:          matchingGroups,
		Filters:                 filters,
		KeyOption:               keyOpt,
	}, nil
}
