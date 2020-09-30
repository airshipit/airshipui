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

	"opendev.org/airship/airshipctl/pkg/events"
	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/log"
	"opendev.org/airship/airshipui/pkg/webservice"
)

// UIEventProcessor basic structure to hold eventsChan, session ID, and errors
type UIEventProcessor struct {
	errors     []error
	eventsChan chan<- events.Event
	sessionID  string
}

// NewUIEventProcessor returns instance of UIEventProcessor for current session ID
func NewUIEventProcessor(id string) events.EventProcessor {
	eventsCh := make(chan events.Event)
	return &UIEventProcessor{
		errors:     []error{},
		eventsChan: eventsCh,
		sessionID:  id,
	}
}

// Process implements EventProcessor interface
func (p *UIEventProcessor) Process(ch <-chan events.Event) error {
	for e := range ch {
		switch e.Type {
		case events.ApplierType:
			log.Errorf("Processing for apply events are not yet implemented")
			p.errors = append(p.errors, e.ErrorEvent.Error)
		case events.ErrorType:
			log.Errorf("Received error on event channel %v", e.ErrorEvent)
			p.errors = append(p.errors, e.ErrorEvent.Error)
		case events.ClusterctlType:
			p.processClusterctlEvent(e.ClusterctlEvent)
		case events.IsogenType:
			p.processIsogenEvent(e.IsogenEvent)
		case events.StatusPollerType:
			log.Errorf("Processing for status poller events are not yet implemented")
			p.errors = append(p.errors, e.ErrorEvent.Error)
		case events.WaitType:
			log.Errorf("Processing for wait events are not yet implemented")
			p.errors = append(p.errors, e.ErrorEvent.Error)
		default:
			log.Errorf("Unknown event type received: %d", e.Type)
			p.errors = append(p.errors, e.ErrorEvent.Error)
		}
	}
	return checkErrors(p.errors)
}

func (p *UIEventProcessor) processIsogenEvent(e events.IsogenEvent) {
	eventType := "isogen"
	msg := e.Message
	switch e.Operation {
	case events.IsogenStart:
		if msg == "" {
			msg = "starting ISO generation"
		}
	case events.IsogenValidation:
		if msg == "" {
			msg = "validation in progress"
		}
	case events.IsogenEnd:
		if msg == "" {
			msg = "ISO generation complete"
		}
	}

	// TODO(mfuller): what shall we do with these events? Pushing
	// them as toasts for now
	sendEventMessage(p.sessionID, eventType, msg)
}

func (p *UIEventProcessor) processClusterctlEvent(e events.ClusterctlEvent) {
	eventType := "clusterctl"
	msg := e.Message

	switch e.Operation {
	case events.ClusterctlInitStart:
		if msg == "" {
			msg = "starting init"
		}
	case events.ClusterctlInitEnd:
		if msg == "" {
			msg = "init completed"
		}
	case events.ClusterctlMoveStart:
		if msg == "" {
			msg = "starting move"
		}
	case events.ClusterctlMoveEnd:
		if msg == "" {
			msg = "move completed"
		}
	}

	sendEventMessage(p.sessionID, eventType, msg)
}

func sendEventMessage(sessionID, eventType, message string) {
	err := webservice.WebSocketSend(configs.WsMessage{
		SessionID:    sessionID,
		Type:         configs.CTL,
		Component:    configs.Document, // probably will change to configs.Phase soon
		SubComponent: configs.Run,
		Message:      fmt.Sprintf("%s: %s", eventType, message),
	})
	if err != nil {
		log.Errorf("Error sending message %s", err)
	}
}

// Check list of errors, and verify that these errors we are able to tolerate
// currently we simply check if the list is empty or not
func checkErrors(errs []error) error {
	if len(errs) != 0 {
		return events.ErrEventReceived{
			Errors: errs,
		}
	}
	return nil
}
