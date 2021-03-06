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

	"opendev.org/airship/airshipctl/pkg/secret/generate"
	"opendev.org/airship/airshipui/pkg/configs"
)

// HandleSecretRequest will flop between requests so we don't have to have them all mapped as function calls
// This will wait for the sub component to complete before responding.  The assumption is this is an async request
func HandleSecretRequest(user *string, request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.Secret,
		SubComponent: request.SubComponent,
	}

	var err error
	var message *string

	subComponent := request.SubComponent
	switch subComponent {
	case configs.Generate:
		message = generatePassphrase()
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

// generatePassphrase will generate a master passphrase to be used for encryption / decryption
func generatePassphrase() *string {
	engine := generate.NewEncryptionKeyEngine(nil)
	masterPassphrase := engine.GenerateEncryptionKey()
	return &masterPassphrase
}
