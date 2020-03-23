/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"opendev.org/airship/airshipctl/pkg/config"
	ctlenv "opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipui/internal/environment"
)

var (
	settings *ctlenv.AirshipCTLSettings
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "airshipui",
	Short:   "airshipui is a graphical user interface for airship",
	Run:     launchOctant,
	Version: environment.Version(),
}

func init() {
	settings = &ctlenv.AirshipCTLSettings{}

	// Add options to rootCmd for configuration of airshipctl and kube config
	settings.InitFlags(rootCmd)

	// Load the airshipctl settings
	settings.InitConfig()

	// Add a 'version' command, in addition to the '--version' option that is auto created
	rootCmd.AddCommand(newVersionCmd())

	// Add flags for the underlying octant dashboard
	addDashboardFlags(rootCmd)
}

func launchOctant(cmd *cobra.Command, args []string) {
	// only process args if there are any
	if cmd.Flags().NFlag() > 0 {
		args = append(args, getFlags(cmd)...)
		fmt.Printf("Executing Octant with the following args: %v\n", args)
	}

	kubeConfig := ""
	airshipKubeConfig := settings.KubeConfigPath()

	// If the kubeconfig specified on the command line (or defaulted to) does not exist,
	// then do not pass specify a kubeconfig.  This will permit the underlying octant
	// to use its own default, which can be either ~/.kube/config or grabbed from the
	// environment
	if fileExists(airshipKubeConfig) {
		kubeConfig = airshipKubeConfig
	}

	RunOctantWithOptions(kubeConfig, args)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Determine whether the given filename exists and is accessible
func fileExists(filename string) bool {
	f, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return !f.IsDir()
}

// this function pulls the passed command line options and renders unto octant what is octant's
func getFlags(cmd *cobra.Command) []string {
	var args []string

	// This will apply all command line arguments to the octant execution depending on its variance from the default values
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Changed {
			name := flag.Name
			// only pass octant specific flags
			// this will need to be refactored if any additional non octant flags are added
			if (name != config.FlagConfigFilePath) && (name != "debug") {
				value := flag.Value
				switch value.Type() {
				case "bool":
					args = append(args, "--"+name)
				default:
					args = append(args, "--"+name, value.String())
				}
			}
		}
	})

	return args
}
