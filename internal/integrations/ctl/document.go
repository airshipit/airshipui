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

	"opendev.org/airship/airshipctl/pkg/document/pull"
	"opendev.org/airship/airshipui/internal/configs"
)

// HandleDocumentRequest will flop between requests so we don't have to have them all mapped as function calls
func HandleDocumentRequest(request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.AirshipCTL,
		Component:    configs.Document,
		SubComponent: request.SubComponent,
	}

	var err error
	var message string
	switch request.SubComponent {
	case configs.GetDefaults:
		response.HTML, err = getDocumentHTML()
	case configs.DocPull:
		message, err = c.docPull()
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

func (c *client) docPull() (string, error) {
	var message string
	settings := pull.Settings{AirshipCTLSettings: c.settings}
	err := settings.Pull()
	if err == nil {
		message = fmt.Sprintf("Success")
	}

	return message, err
}

func getDocumentHTML() (string, error) {
	return getHTML("./internal/integrations/ctl/templates/document.html", ctlPage{
		Title:   "Document",
		Version: getAirshipCTLVersion(),
	})
}
