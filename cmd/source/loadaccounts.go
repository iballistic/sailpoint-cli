// Copyright (c) 2021, SailPoint Technologies, Inc. All rights reserved.
package source

import (
	"context"

	"github.com/charmbracelet/log"

	"github.com/sailpoint-oss/sailpoint-cli/internal/config"
	"github.com/spf13/cobra"
)

func newLoadCommand() *cobra.Command {
	var disableOptimization bool
	cmd := &cobra.Command{
		Use:     "load",
		Short:   "Starts an account aggregation on the specified source",
		Long:    "\nStarts an account aggregation on the specified source by its source id\n\n",
		Example: "sail source load 7e38967c99504eb3ac6095296007e015",
		Aliases: []string{"ld"},
		RunE: func(cmd *cobra.Command, args []string) error {

			apiClient, err := config.InitAPIClient(false)
			if err != nil {
				return err
			}

			for _, sourceID := range args {

				log.Info("Starting account aggregation", "sourceID", sourceID)
				request := apiClient.Beta.SourcesAPI.ImportAccounts(context.TODO(), sourceID)
				request = request.DisableOptimization("false")
				if disableOptimization {
					request = request.DisableOptimization("true")
					log.Info("Aggregation optimization is disabled")
				}

				_, _, err := request.Execute()
				if err != nil {
					return err
				} else {
					log.Info("Started account aggregation", "sourceID", sourceID)
				}

			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&disableOptimization, "DisableOptimization","d", false, "Disable Source Aggregation Optimization")
	return cmd
}
