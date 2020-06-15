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

	"opendev.org/airship/airshipctl/pkg/bootstrap/isogen"
	"opendev.org/airship/airshipui/internal/configs"
)

// HandleBaremetalRequest will flop between requests so we don't have to have them all mapped as function calls
func HandleBaremetalRequest(request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.AirshipCTL,
		Component:    configs.Baremetal,
		SubComponent: request.SubComponent,
	}

	var err error
	var message string
	switch request.SubComponent {
	case configs.GetDefaults:
		response.HTML, err = getBaremetalHTML()
	case configs.DocPull:
		message, err = c.docPull()
	case configs.GenerateISO:
		message, err = c.generateIso()
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

func (c *client) generateIso() (string, error) {
	var message string
	err := isogen.GenerateBootstrapIso(c.settings)
	if err == nil {
		message = fmt.Sprintf("Success")
	}

	return message, err
}

func getBaremetalHTML() (string, error) {
	return getHTML("./internal/integrations/ctl/templates/baremetal.html", ctlPage{
		Title:   "Baremetal",
		Version: getAirshipCTLVersion(),
	})
}
