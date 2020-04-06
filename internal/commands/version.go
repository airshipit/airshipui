/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// version will be overridden by ldflags supplied in Makefile
	version = "(dev-version)"
)

func newVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Version for airshipui binary",
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()

			fmt.Fprintln(out, "airshipui version", Version())
		},
	}
	return versionCmd
}

func Version() string {
	return version
}
