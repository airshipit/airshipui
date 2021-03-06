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
	"opendev.org/airship/airshipctl/pkg/document/pull"
	"opendev.org/airship/airshipui/pkg/configs"
)

// HandleDocumentRequest will flop between requests so we don't have to have them all mapped as function calls
func HandleDocumentRequest(user *string, request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.Document,
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
	case configs.Pull:
		message, err = client.docPull()
	case configs.Plugin:
		err = fmt.Errorf("Subcomponent %s not implemented", request.SubComponent)
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

func (c *Client) docPull() (*string, error) {
	var message *string
	cfgFactory := config.CreateFactory(configs.UIConfig.AirshipConfigPath)
	// 2nd arg is noCheckout, I assume we want to checkout the repo,
	// so setting to false
	err := pull.Pull(cfgFactory, false)
	if err == nil {
		s := "Success"
		message = &s
	}

	return message, err
}
