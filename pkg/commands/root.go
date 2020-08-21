/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package commands

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"

	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/log"
	"opendev.org/airship/airshipui/pkg/webservice"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "airshipui",
	Short:   "airshipui is a graphical user interface for airship",
	Run:     launch,
	Version: Version(),
}

func init() {
	// Add a 'version' command, in addition to the '--version' option that is auto created
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.Flags().IntVar(
		&log.LogLevel,
		"loglevel",
		6,
		"This well set the log level, anything at or below that level will be viewed, all others suppressed\n"+
			"  6 -- Trace\n"+
			"  5 -- Debug\n"+
			"  4 -- Info\n"+
			"  3 -- Warn\n"+
			"  2 -- Error\n"+
			"  1 -- Fatal\n",
	)
}

func launch(cmd *cobra.Command, args []string) {
	// set default config path
	// TODO: do we want to make this a flag that can be passed in?
	airshipUIConfigPath, err := getDefaultConfigPath()
	if err != nil {
		log.Errorf("Error setting config path %s", err)
	}

	// Read AirshipUI config file
	if err := configs.SetUIConfig(airshipUIConfigPath); err != nil {
		log.Errorf("config %s", err)
	}

	// start webservice and listen for the the ctl + c to exit
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Exiting the webservice")
		os.Exit(0)
	}()
	webservice.WebServer()
}

// Execute is called from the main program and kicks this whole shindig off
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func getDefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.FromSlash(home + "/.airship/airshipui.json"), nil
}
