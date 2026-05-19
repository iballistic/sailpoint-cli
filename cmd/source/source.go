// Copyright (c) 2021, SailPoint Technologies, Inc. All rights reserved.
package source

import (
	"github.com/spf13/cobra"
)

func NewSourceformCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "source",
		Short:   "Manage sources in Identity Security Cloud",
		Long:    "\nManage sources in Identity Security Cloud\n\n",
		Example: "sail source | sail src",
		Aliases: []string{"src"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newListCommand(),
		newLoadCommand(),
	)

	return cmd
}
