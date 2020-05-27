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
	"sync"
	"syscall"

	"github.com/spf13/cobra"

	"opendev.org/airship/airshipui/internal/configs"
	"opendev.org/airship/airshipui/internal/electron"
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
	// only process args if there are any
	// TODO: what flags do we care about and what shall we do with them?
	if cmd.Flags().NFlag() > 0 {
		log.Printf("Executing AirshipUI with the following args: %v\n", args)
	}

	// start the webservice
	go webservice.WebServer()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	waitgrp := sync.WaitGroup{}

	// Read AirshipUI config file
	if err := configs.GetConfigFromFile(); err == nil {
		// launch any plugins marked as autoStart: true in airshipui.json
		for _, p := range configs.UiConfig.Plugins {
			if p.Executable.AutoStart {
				waitgrp.Add(1)
				go RunBinaryWithOptions(
					ctx,
					p.Executable.Filepath,
					p.Executable.Args,
					&waitgrp,
					sigs,
				)
			}
		}
	} else {
		log.Printf("config %s", err)
		webservice.Alerts = append(
			webservice.Alerts,
			webservice.Alert{
				Level:   "info",
				Message: fmt.Sprintf("%s", err),
			},
		)
	}

	// start the electron app
	err := electron.RunElectron()

	// cancel running plugins and wait for shut down
	cancel()
	waitgrp.Wait()

	if err != nil {
		log.Printf("Exit %s", err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
