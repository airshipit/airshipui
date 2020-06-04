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

var headless bool
var remote bool

func init() {
	// Add a 'version' command, in addition to the '--version' option that is auto created
	rootCmd.AddCommand(newVersionCmd())

	// add the remote & headless options in case people want to run a split setup
	rootCmd.Flags().BoolVar(&headless, "headless", false, "start the system in headless webserver only, no ui.")
	rootCmd.Flags().BoolVar(&remote, "remote", false, "start the system in remote ui only, no webserver.")
}

func launch(cmd *cobra.Command, args []string) {
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
		webservice.SendAlert(webservice.Info, fmt.Sprintf("%s", err))
	}

	// just a little ditty to see if we should open the ui or the webservice or both
	switch handleStartType() {
	case "headless":
		// start webservice and listen for the the ctl + c to exit
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			log.Println("Exiting the webservice")
			os.Exit(0)
		}()
		webservice.WebServer()
	case "remote":
		// start the electron app
		startElectron()
	default:
		// start webservice and electron
		go webservice.WebServer()
		startElectron()
	}

	// cancel running plugins and wait for shut down
	cancel()
	waitgrp.Wait()
}

func startElectron() {
	err := electron.RunElectron()
	if err != nil {
		log.Printf("Exit %s", err)
	}
}

// TODO: determine if cobra can make flags exclusive without the extra logic
func handleStartType() string {
	st := "default"
	if remote && headless {
		log.Fatalf("Cannot set both --remote and --headless flags")
	}

	if remote {
		st = "remote"
	} else if headless {
		st = "headless"
	}

	return st
}

// Execute is called from the main program and kicks this whole shindig off
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
