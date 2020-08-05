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
	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipui/pkg/configs"
)

// CTLFunctionMap is a function map for the CTL functions that is referenced in the webservice
var CTLFunctionMap = map[configs.WsComponentType]func(configs.WsMessage) configs.WsMessage{
	configs.Baremetal: HandleBaremetalRequest,
	configs.Document:  HandleDocumentRequest,
}

// maintain the state of a potentially long running process
var runningRequests map[configs.WsSubComponentType]bool = make(map[configs.WsSubComponentType]bool)

// Client provides a library of functions that enable external programs (e.g. Airship UI) to perform airshipctl
// functionality in exactly the same manner as the CLI.
type Client struct {
	settings *environment.AirshipCTLSettings
}

// NewClient initializes the airshipctl client for external usage.
func NewClient() *Client {
	settings := &environment.AirshipCTLSettings{}
	settings.InitConfig()

	c := &Client{
		settings: settings,
	}

	// set verbosity to true
	c.settings.Debug = true

	return c
}

// initilize the connection to airshipctl
var c *Client = NewClient()
