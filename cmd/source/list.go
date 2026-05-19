// Copyright (c) 2021, SailPoint Technologies, Inc. All rights reserved.
package source

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	v3 "github.com/sailpoint-oss/golang-sdk/v2/api_v3"
	"github.com/sailpoint-oss/sailpoint-cli/internal/config"
	"github.com/sailpoint-oss/sailpoint-cli/internal/output"
	"github.com/sailpoint-oss/sailpoint-cli/internal/sdk"

	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	var filters string
	var outputFile string
	var prettyPrint bool
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all sources in Identity Security Cloud",
		Long:    "\nList all sources in Identity Security Cloud\n\n",
		Example: "sail source list | sail source ls",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			apiClient, err := config.InitAPIClient(false)
			if err != nil {
				return err
			}

			sources, resp, err := sailpoint.PaginateWithDefaults[v3.Source](
				apiClient.V3.SourcesAPI.ListSources(context.TODO()).Filters(filters),
			)
			if err != nil {
				return sdk.HandleSDKError(resp, err)
			}

			defer resp.Body.Close()

			var entries [][]string

			for _, v := range sources {
				entries = append(entries, []string{v.Name, *v.Id})
			}

			// Read response body
			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}

			// Check if response is JSON and pretty print if requested
			if prettyPrint {
				var jsonData interface{}
				if err := json.Unmarshal(responseBody, &jsonData); err == nil {
					prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
					if err == nil {
						responseBody = prettyJSON
					}
				}
			}

			// Output to file or stdout
			if outputFile != "" {
				if err := writeToFile(outputFile, responseBody); err != nil {
					return fmt.Errorf("failed to write to file: %w", err)
				}
				fmt.Printf("Response saved to %s\n", outputFile)
			} else {
				//fmt.Fprint(cmd.OutOrStdout(), string(responseBody))
				output.WriteTable(cmd.OutOrStdout(), []string{"Name", "ID"}, entries, "Name")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&filters, "query", "q", "", "Filter to search and return matching sources")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file to save the response (if not specified, prints basic info to stdout)")
	cmd.Flags().BoolVarP(&prettyPrint, "pretty", "p", false, "Pretty print JSON response")
	return cmd
}

// writeToFile writes data to a file
func writeToFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}


//example: sailpoint-cli.exe source ls -q "connectorName  in (\"Active Directory\",\"Workday\")"