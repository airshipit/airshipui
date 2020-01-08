/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	env "opendev.org/airship/airshipui/internal/environment"
)

func newVersionCmd() *cobra.Command {

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Version for airshipui binary",
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()

			fmt.Fprintln(out, "airshipui version", env.Version())
		},
	}
	return versionCmd
}
