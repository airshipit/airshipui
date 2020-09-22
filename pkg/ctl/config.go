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

	"opendev.org/airship/airshipui/pkg/configs"
)

// HandleConfigRequest will flop between requests so we don't have to have them all mapped as function calls
// This will wait for the sub component to complete before responding.  The assumption is this is an async request
func HandleConfigRequest(request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.Baremetal,
		SubComponent: request.SubComponent,
	}

	var err error
	var message string

	subComponent := request.SubComponent
	switch subComponent {
	case configs.GetContext:
		err = fmt.Errorf("Subcomponent %s not implemented", request.SubComponent)
	case configs.GetEncryptionConfig:
		err = fmt.Errorf("Subcomponent %s not implemented", request.SubComponent)
	case configs.GetManagementConfig:
		err = fmt.Errorf("Subcomponent %s not implemented", request.SubComponent)
	case configs.GetManifest:
		err = fmt.Errorf("Subcomponent %s not implemented", request.SubComponent)
	case configs.Init:
		err = fmt.Errorf("Subcomponent %s not implemented", request.SubComponent)
	case configs.SetContext:
		err = fmt.Errorf("Subcomponent %s not implemented", request.SubComponent)
	case configs.SetEncryptionConfig:
		err = fmt.Errorf("Subcomponent %s not implemented", request.SubComponent)
	case configs.SetManagementConfig:
		err = fmt.Errorf("Subcomponent %s not implemented", request.SubComponent)
	case configs.SetManifest:
		err = fmt.Errorf("Subcomponent %s not implemented", request.SubComponent)
	case configs.UseContext:
		err = fmt.Errorf("Subcomponent %s not implemented", request.SubComponent)
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
