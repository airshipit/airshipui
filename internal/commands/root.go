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
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/spf13/cobra"

	"opendev.org/airship/airshipui/internal/configs"
	"opendev.org/airship/airshipui/internal/webservice"
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
}

func launch(cmd *cobra.Command, args []string) {
	// set default config path
	// TODO: do we want to make this a flag that can be passed in?
	airshipUIConfigPath, err := getDefaultConfigPath()
	if err != nil {
		log.Printf("Error setting config path %s", err)
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	waitgrp := sync.WaitGroup{}

	// Read AirshipUI config file
	if err := configs.SetUIConfig(airshipUIConfigPath); err == nil {
		// launch any plugins marked as autoStart: true in airshipui.json
		for _, dashboard := range configs.UIConfig.Dashboards {
			if dashboard.Executable != nil {
				if dashboard.Executable.AutoStart {
					waitgrp.Add(1)
					go RunBinaryWithOptions(
						ctx,
						dashboard.Executable.Filepath,
						dashboard.Executable.Args,
						&waitgrp,
						sigs,
					)
				}
			}
		}
	} else {
		log.Printf("config %s", err)
		webservice.SendAlert(configs.Info, fmt.Sprintf("%s", err), true)
	}

	// start the web service and related sundries
	webservice.WebServer()

	// cancel running plugins and wait for shut down
	cancel()
	waitgrp.Wait()
}

// Execute is called from the main program and kicks this whole shindig off
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
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
