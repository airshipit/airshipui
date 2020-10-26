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

package task

import (
	"fmt"
	"time"

	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/log"
	"opendev.org/airship/airshipui/pkg/webservice"
)

// RunningTasks serves as a cache for currently running tasks
// TODO(mfuller): keeping a backend cache may not be necessary since
// task objects get attached to the event processors injected into phase
// clients and will always be updated directly from there. But we may
// want to use it to ensure the frontend can retrieve running tasks
// in the event of a browser refresh
var RunningTasks = map[string]Task{}

// Task simple structure to hold details about a long running task
type Task struct {
	ID        string
	SessionID string
	Name      string
	Progress  Progress
	Running   bool // TODO(mfuller): this is probably only necessary on the frontend
}

// Progress structure to store and pass progress data for a running task
type Progress struct {
	StartTime   int64    `json:"startTime"`
	EndTime     int64    `json:"endTime"`
	LastUpdated int64    `json:"lastUpdated"`
	TotalSteps  int      `json:"totalSteps"`
	CurrentStep int      `json:"currentStep"`
	Message     string   `json:"message"`
	Errors      []string `json:"errors"`
}

// HandleTaskRequest handles incoming WS messages for tasks
// TODO(mfuller): it's unclear how often this will happen. Task requests
// and updates will almost always come from event processing. Is this needed?
func HandleTaskRequest(request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.UI,
		Component:    configs.Task,
		SubComponent: request.SubComponent,
	}

	var err error
	var message *string

	switch request.SubComponent {
	// case configs.TaskStart:
	// 	response.ID, response.Data = StartTask(request.SessionID, request.Name)
	case configs.TaskRemove:
		message, err = RemoveTask(request.ID)
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

// NewTask returns a pointer to a new Task built with a session ID, name, and UUID
func NewTask(sessionID, taskID, name string) *Task {
	task := Task{
		ID:        taskID,
		SessionID: sessionID,
		Name:      name,
		Progress: Progress{
			StartTime:   time.Now().UnixNano() / 1000000,
			TotalSteps:  0, // will steps be determinable at task start?
			CurrentStep: 1,
			Errors:      []string{},
		},
		Running: true,
	}

	RunningTasks[task.ID] = task

	return &task
}

// RemoveTask removes a Task from RunningTasks and sends confirmation
// message to UI. This function is intended to be called by the frontend
// client by clicking a "remove" button in the task manager
func RemoveTask(id string) (*string, error) {
	if t, ok := RunningTasks[id]; ok {
		delete(RunningTasks, id)
		msg := fmt.Sprintf("Removed task '%s'", t.Name)
		return &msg, nil
	}

	return nil, fmt.Errorf("Task with id %s not found", id)
}

// UpdateTask updates a task with new progress details
// TODO(mfuller): I don't know if this function is even necessary
// since most updates are going to come from event processing,
// so the message will likely be fired from the processor directly
func UpdateTask(sessionID, id string, progress Progress) {
	if t, ok := RunningTasks[id]; ok {
		t.SendTaskMessage(configs.TaskUpdate, progress)
	}

	// this is the only reason we need session ID, otherwise we'd have
	// to just log a message and walk away...
	m := fmt.Sprintf("Task with id %s not found", id)
	err := webservice.WebSocketSend(configs.WsMessage{
		SessionID: sessionID,
		Type:      configs.UI,
		Error:     &m,
	})

	if err != nil {
		log.Errorf("Error sending message for task %s", err)
	}
}

// SendTaskMessage allows a running Task to push progress updates to the frontend client
func (t *Task) SendTaskMessage(subComponent configs.WsSubComponentType, progress Progress) {
	err := webservice.WebSocketSend(configs.WsMessage{
		SessionID:    t.SessionID,
		ID:           t.ID,
		Name:         t.Name,
		Timestamp:    time.Now().UnixNano() / 1000000,
		Type:         configs.UI,
		Component:    configs.Task,
		SubComponent: subComponent,
		Message:      &t.Name,
		Data:         progress,
	})

	if err != nil {
		log.Errorf("Error sending message for task %s", err)
	}
}
