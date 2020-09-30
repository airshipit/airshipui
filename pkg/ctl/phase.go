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
	"encoding/json"
	"fmt"

	"opendev.org/airship/airshipctl/pkg/events"
	"opendev.org/airship/airshipctl/pkg/phase"
	"opendev.org/airship/airshipctl/pkg/phase/ifc"
	"opendev.org/airship/airshipui/pkg/configs"
)

// HandlePhaseRequest will flop between requests so we don't have to have them all mapped as function calls
// This will wait for the sub component to complete before responding.  The assumption is this is an async request
func HandlePhaseRequest(request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.Document, // setting this to Document for now since that's handling phase requests
		SubComponent: request.SubComponent,
	}

	var err error
	var message string
	var valid bool

	client, err := NewClient(AirshipConfigPath, KubeConfigPath, request)
	if err != nil {
		response.Error = err.Error()
		return response
	}

	subComponent := request.SubComponent
	switch subComponent {
	case configs.Plan:
		err = fmt.Errorf("Subcomponent %s not implemented", request.SubComponent)
	case configs.Render:
		err = fmt.Errorf("Subcomponent %s not implemented", request.SubComponent)
	case configs.Run:
		err = client.RunPhase(request)
	case configs.ValidatePhase:
		valid, err = client.ValidatePhase(request.ID, request.SessionID)
		message = validateHelper(valid)
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

// this helper function will likely disappear once a clear workflow for
// phase validation takes shape in UI. For now, it simply returns a
// string message to be displayed as a toast in frontend client
func validateHelper(valid bool) string {
	msg := "invalid"
	if valid {
		msg = "valid"
	}
	return msg
}

// ValidatePhase validates the specified phase
// (ifc.Phase.Validate isn't implemented yet, so this function
// currently always returns "valid")
func (c *Client) ValidatePhase(id, sessionID string) (bool, error) {
	phase, err := getPhaseIfc(id, sessionID)
	if err != nil {
		return false, err
	}

	err = phase.Validate()
	if err != nil {
		return false, err
	}

	return true, nil
}

// RunPhase runs the selected phase
func (c *Client) RunPhase(request configs.WsMessage) error {
	phase, err := getPhaseIfc(request.ID, request.SessionID)
	if err != nil {
		return err
	}

	opts := ifc.RunOptions{}
	if request.Data != nil {
		bytes, err := json.Marshal(request.Data)
		if err != nil {
			return err
		}

		err = json.Unmarshal(bytes, &opts)
		if err != nil {
			return err
		}
	}

	return phase.Run(opts)
}

// helper function to return a Phase interface based on a JSON
// string representation of an ifc.ID value
func getPhaseIfc(id, sessionID string) (ifc.Phase, error) {
	phaseID := ifc.ID{}

	err := json.Unmarshal([]byte(id), &phaseID)
	if err != nil {
		return nil, err
	}

	helper, err := getHelper()
	if err != nil {
		return nil, err
	}

	var procFunc phase.ProcessorFunc
	procFunc = func() events.EventProcessor {
		return NewUIEventProcessor(sessionID)
	}

	// inject event processor to phase client
	proc := phase.InjectProcessor(procFunc)

	client := phase.NewClient(helper, proc)

	phase, err := client.PhaseByID(phaseID)
	if err != nil {
		return nil, err
	}

	return phase, nil
}
