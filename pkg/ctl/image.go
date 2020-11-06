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
	"fmt"

	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipctl/pkg/phase"
	"opendev.org/airship/airshipui/pkg/configs"
)

// HandleImageRequest will flop between requests so we don't have to have them all mapped as function calls
// This will wait for the sub component to complete before responding.  The assumption is this is an async request
func HandleImageRequest(user *string, request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.Baremetal,
		SubComponent: request.SubComponent,
	}

	var err error
	var message *string

	client, err := NewClient(configs.UIConfig.AirshipConfigPath, request)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	subComponent := request.SubComponent
	switch subComponent {
	case configs.Generate:
		// since this is long running cache it up
		// TODO: Test before running the geniso
		runningRequests[subComponent] = true
		message, err = client.generateIso()
		// now that we're done forget we did anything
		delete(runningRequests, subComponent)
	default:
		err = fmt.Errorf("Subcomponent %s not found", request.SubComponent)
	}

	if err != nil {
		e := err.Error()
		response.Error = &e
	} else {
		response.Message = message
	}

	return response
}

// generate iso now just runs a phase and not an individual command
func (c *Client) generateIso() (*string, error) {
	cfgFactory := config.CreateFactory(configs.UIConfig.AirshipConfigPath)
	p := &phase.RunCommand{
		Factory: cfgFactory,
	}
	p.Options.PhaseID.Name = config.BootstrapPhase
	s := "Success"
	return &s, p.RunE()
}
