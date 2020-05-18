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
package ctl

import (
	"log"

	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/pkg/version"
	"opendev.org/airship/airshipui/internal/configs"
)

// client provides a library of functions that enable external programs (e.g. Airship UI) to perform airshipctl
// functionality in exactly the same manner as the CLI.
type client struct {
	settings *environment.AirshipCTLSettings
}

// NewClient initializes the airshipctl client for external usage.
func NewClient() *client {
	settings := &environment.AirshipCTLSettings{}
	settings.InitConfig()

	c := &client{
		settings: settings,
	}

	return c
}

// initilize the connection to airshipctl
var c *client = NewClient()

// GetAirshipCTLVersion will kick out what version of airshipctl we're using
func GetAirshipCTLVersion() string {
	return version.Get().GitVersion
}

// GetDefaults will send to the UI the basics of what airshipctl we know about
func GetDefaults(configs.WsMessage) configs.WsMessage {
	config, err := getDefaultHTML()
	if err != nil {
		config = "Error attempting to get data for AirshipCTL configs: " + err.Error()
		log.Println(err)
	}

	return configs.WsMessage{
		Type:      configs.AirshipCTL,
		Component: configs.Info,
		HTML:      config,
	}
}
