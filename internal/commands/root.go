/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/
package commands

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipui/internal/electron"
	"opendev.org/airship/airshipui/internal/webservice"
)

var (
	settings *environment.AirshipCTLSettings
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "airshipui",
	Short:   "airshipui is a graphical user interface for airship",
	Run:     launch,
	Version: Version(),
}

func init() {
	settings = &environment.AirshipCTLSettings{}

	// Add options to rootCmd for configuration of airshipctl and kube config
	settings.InitFlags(rootCmd)

	// Load the airshipctl settings
	settings.InitConfig()

	// Add a 'version' command, in addition to the '--version' option that is auto created
	rootCmd.AddCommand(newVersionCmd())

	// Add flags for the underlying octant dashboard
	addDashboardFlags(rootCmd)
}

func launch(cmd *cobra.Command, args []string) {
	// only process args if there are any
	if cmd.Flags().NFlag() > 0 {
		args = append(args, getFlags(cmd)...)
		log.Printf("Executing AirshipUI with the following args: %v\n", args)
	}

	// start the webservice
	go webservice.WebServer()

	// start the electron app
	electron.RunElectron()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
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

// some day this may need to get refactored if the options become transportable from external sources
func addDashboardFlags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = true

	// octant specific flags
	cmd.Flags().StringP("context", "", "", "initial context")
	cmd.Flags().String("kubeconfig", "", "absolute path to kubeConfig file")
	cmd.Flags().StringP("namespace", "n", "", "initial namespace")
	cmd.Flags().StringP("plugin-path", "", "", "plugin path")
	cmd.Flags().BoolP("verbose", "v", false, "turn on debug logging")
}
