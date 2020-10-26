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
	"time"

	"opendev.org/airship/airshipctl/pkg/events"
	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/log"
	"opendev.org/airship/airshipui/pkg/task"
	applyevent "sigs.k8s.io/cli-utils/pkg/apply/event"
)

// TODO(mfuller): I'll need to implement at least some no-op event
// processors for the remaining types, otherwise tasks don't get added
// to the frontend, and I can't process errors for them either

// UIEventProcessor basic structure to hold eventsChan, session ID, and errors
type UIEventProcessor struct {
	errors     []error
	eventsChan chan<- events.Event
	sessionID  string
	task       *task.Task
}

// NewUIEventProcessor returns instance of UIEventProcessor for current session ID
func NewUIEventProcessor(sessionID string, task *task.Task) events.EventProcessor {
	eventsCh := make(chan events.Event)
	return &UIEventProcessor{
		errors:     []error{},
		eventsChan: eventsCh,
		sessionID:  sessionID,
		task:       task,
	}
}

// Process implements EventProcessor interface
func (p *UIEventProcessor) Process(ch <-chan events.Event) error {
	for e := range ch {
		switch e.Type {
		case events.ApplierType:
			p.processApplierEvent(e.ApplierEvent)
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
	return p.checkErrors()
}

// TODO(mfuller): this function currently only adds errors if present,
// otherwise it sends a task message with the entire applyevent.Event
// object. At some point, we'll probably want to see how the printer
// being used in ctl is determining what to print out to console and
// do something similar
func (p *UIEventProcessor) processApplierEvent(e applyevent.Event) {
	var sub configs.WsSubComponentType
	eventType := "kubernetes applier"
	var msg string

	if e.Type == applyevent.ErrorType {
		p.errors = append(p.errors, e.ErrorEvent.Err)
		return
	}

	if e.Type == applyevent.ApplyType {
		switch e.ApplyEvent.Type {
		case applyevent.ApplyEventCompleted:
			sub = configs.TaskEnd
			msg = "completed"
			p.task.Progress.EndTime = time.Now().UnixNano() / 1000000
		default:
			sub = configs.TaskUpdate
			msg = fmt.Sprintf("%+v", e)
		}
	}

	message := fmt.Sprintf("%s: %s", eventType, msg)
	p.task.Progress.LastUpdated = time.Now().UnixNano() / 1000000
	p.task.Progress.Message = message

	p.task.SendTaskMessage(sub, p.task.Progress)
}

func (p *UIEventProcessor) processIsogenEvent(e events.IsogenEvent) {
	var sub configs.WsSubComponentType
	eventType := "isogen"
	msg := e.Message
	switch e.Operation {
	case events.IsogenStart:
		sub = configs.TaskUpdate
		if msg == "" {
			msg = "starting ISO generation"
		}
	case events.IsogenValidation:
		sub = configs.TaskUpdate
		p.task.Progress.LastUpdated = time.Now().UnixNano() / 1000000
		if msg == "" {
			msg = "validation in progress"
		}
	case events.IsogenEnd:
		sub = configs.TaskEnd
		if msg == "" {
			msg = "ISO generation complete"
		}
		p.task.Progress.EndTime = time.Now().UnixNano() / 1000000
	}

	message := fmt.Sprintf("%s: %s", eventType, msg)

	p.task.Progress.LastUpdated = time.Now().UnixNano() / 1000000
	p.task.Progress.Message = message

	p.task.SendTaskMessage(sub, p.task.Progress)
}

func (p *UIEventProcessor) processClusterctlEvent(e events.ClusterctlEvent) {
	var sub configs.WsSubComponentType
	eventType := "clusterctl"
	msg := e.Message

	switch e.Operation {
	case events.ClusterctlInitStart:
		sub = configs.TaskUpdate
		if msg == "" {
			msg = "starting init"
		}
	case events.ClusterctlInitEnd:
		sub = configs.TaskEnd
		p.task.Progress.EndTime = time.Now().UnixNano() / 1000000
		if msg == "" {
			msg = "init completed"
		}
	case events.ClusterctlMoveStart:
		sub = configs.TaskUpdate
		if msg == "" {
			msg = "starting move"
		}
	case events.ClusterctlMoveEnd:
		sub = configs.TaskEnd
		p.task.Progress.EndTime = time.Now().UnixNano() / 1000000
		if msg == "" {
			msg = "move completed"
		}
	}

	message := fmt.Sprintf("%s: %s", eventType, msg)

	p.task.Progress.LastUpdated = time.Now().UnixNano() / 1000000
	p.task.Progress.Message = message

	p.task.SendTaskMessage(sub, p.task.Progress)
}

// Check list of errors, and verify that these errors we are able to tolerate
// currently we simply check if the list is empty or not
func (p *UIEventProcessor) checkErrors() error {
	log.Infof("p.errors: %+v", p.errors)
	if len(p.errors) != 0 {
		for _, e := range p.errors {
			p.task.Progress.Errors = append(p.task.Progress.Errors, e.Error())
		}
		p.task.SendTaskMessage(configs.TaskUpdate, p.task.Progress)
		return events.ErrEventReceived{
			Errors: p.errors,
		}
	}
	return nil
}
