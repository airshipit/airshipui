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
	"path/filepath"

	"opendev.org/airship/airshipctl/pkg/bootstrap/isogen"
	"opendev.org/airship/airshipui/internal/configs"
)

// HandleBaremetalRequest will flop between requests so we don't have to have them all mapped as function calls
// This will wait for the sub component to complete before responding.  The assumption is this is an async request
func HandleBaremetalRequest(request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.AirshipCTL,
		Component:    configs.Baremetal,
		SubComponent: request.SubComponent,
	}

	var err error
	var message string
	subComponent := request.SubComponent
	switch subComponent {
	case configs.GetDefaults:
		response.HTML, err = getBaremetalHTML()
	case configs.GenerateISO:
		// since this is long running cache it up
		runningRequests[subComponent] = true
		message, err = c.generateIso()
		// now that we're done forget we did anything
		delete(runningRequests, subComponent)
	default:
		err = fmt.Errorf("Subcomponent %s not found", request.SubComponent)
	}

	if err != nil {
		response.Error = err.Error()
	} else {
		response.Message = message
	}

	return response
}

func (c *Client) generateIso() (string, error) {
	var message string
	err := isogen.GenerateBootstrapIso(c.settings)
	if err == nil {
		message = fmt.Sprintf("Success")
	}

	return message, err
}

func getBaremetalHTML() (string, error) {
	p := ctlPage{
		Title:      "Baremetal",
		Version:    getAirshipCTLVersion(),
		ButtonText: "Generate ISO",
	}

	if _, ok := runningRequests[configs.GenerateISO]; ok {
		p.Disabled = "disabled"
		p.ButtonText = "In Progress"
	}

	return getHTML(filepath.Join(basepath, "/templates/baremetal.html"), p)
}
