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
	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/pkg/log"
	"opendev.org/airship/airshipui/pkg/configs"
	uiLog "opendev.org/airship/airshipui/pkg/log"
	"opendev.org/airship/airshipui/pkg/webservice"
)

// CTLFunctionMap is a function map for the CTL functions that is referenced in the webservice
var CTLFunctionMap = map[configs.WsComponentType]func(configs.WsMessage) configs.WsMessage{
	configs.Baremetal: HandleBaremetalRequest,
	configs.Document:  HandleDocumentRequest,
}

// maintain the state of a potentially long running process
var runningRequests map[configs.WsSubComponentType]bool = make(map[configs.WsSubComponentType]bool)

// Client provides a library of functions that enable external programs (e.g. Airship UI) to perform airshipctl
// functionality in exactly the same manner as the CLI.
type Client struct {
	settings *environment.AirshipCTLSettings
}

// LogInterceptor is just a struct to hold a pointer to the remote channel
type LogInterceptor struct {
	response configs.WsMessage
}

// Init allows for the circular reference to the webservice package to be broken and allow for the sending
// of arbitrary messages from any package to the websocket
func Init() {
	webservice.AppendToFunctionMap(configs.CTL, map[configs.WsComponentType]func(configs.WsMessage) configs.WsMessage{
		configs.Baremetal: HandleBaremetalRequest,
		configs.Document:  HandleDocumentRequest,
	})
}

// NewDefaultClient initializes the airshipctl client for external usage with default logging.
func NewDefaultClient() *Client {
	settings := &environment.AirshipCTLSettings{}
	// ensure no error if airship config doesn't exist
	settings.Create = true
	settings.InitConfig()

	client := &Client{
		settings: settings,
	}

	// set verbosity to true
	client.settings.Debug = true

	return client
}

// NewClient initializes the airshipctl client for external usage with the logging overridden.
func NewClient(request configs.WsMessage) *Client {
	client := NewDefaultClient()

	// init the interceptor to send messages to the UI
	// TODO: Unsure how this will be handled with overlapping runs
	log.Init(client.settings.Debug, NewLogInterceptor(request))

	return client
}

// NewLogInterceptor will construct a channel writer for use with the logger
func NewLogInterceptor(request configs.WsMessage) *LogInterceptor {
	// TODO: determine if we're only getting stub responses and if we don't have to pick things out that we care about
	// This is a stub response used by the writer to kick out messages to the UI
	response := configs.WsMessage{
		Type:      configs.UI,
		Component: configs.Log,
		SessionID: request.SessionID,
	}

	return &LogInterceptor{
		response: response,
	}
}

// Write satisfies the implementation of io.Writer.
// The intention is to hijack the log output for a progress bar on the UI
func (cw *LogInterceptor) Write(data []byte) (n int, err error) {
	response := cw.response
	response.Message = string(data)
	if err = webservice.WebSocketSend(response); err != nil {
		uiLog.Errorf("Error receiving / sending message: %s\n", err)
		return len(data), err
	}

	return len(data), nil
}
