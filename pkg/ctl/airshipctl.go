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
	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipctl/pkg/log"
	"opendev.org/airship/airshipui/pkg/configs"
	uiLog "opendev.org/airship/airshipui/pkg/log"
	"opendev.org/airship/airshipui/pkg/webservice"
)

// AirshipConfigPath location of airship config (default $HOME/.airship.config)
// TODO(mfuller): are we going to retrieve these from the environment / cli options?
// leaving them both unset (nil) for now so that the default locations will be used
var AirshipConfigPath *string

// KubeConfigPath location of kubeconfig used by airshipctl (default $HOME/.airship/kubeconfig)
var KubeConfigPath *string

// CTLFunctionMap is a function map for the CTL functions that is referenced in the webservice
var CTLFunctionMap = map[configs.WsComponentType]func(*string, configs.WsMessage) configs.WsMessage{
	configs.Baremetal: HandleBaremetalRequest,
	configs.Cluster:   HandleClusterRequest,
	configs.CTLConfig: HandleConfigRequest,
	configs.Document:  HandleDocumentRequest,
	configs.Image:     HandleImageRequest,
	configs.Phase:     HandlePhaseRequest,
	configs.Secret:    HandleSecretRequest,
}

// maintain the state of a potentially long running process
var runningRequests map[configs.WsSubComponentType]bool = make(map[configs.WsSubComponentType]bool)

// Client provides a library of functions that enable external programs (e.g. Airship UI) to perform airshipctl
// functionality in exactly the same manner as the CLI.
type Client struct {
	Config *config.Config
	Debug  bool // this is a placeholder until I figure out how / where to set this in airshipctl
}

// LogInterceptor is just a struct to hold a pointer to the remote channel
type LogInterceptor struct {
	response configs.WsMessage
}

// Init allows for the circular reference to the webservice package to be broken and allow for the sending
// of arbitrary messages from any package to the websocket
func Init() {
	webservice.AppendToFunctionMap(
		configs.CTL,
		map[configs.WsComponentType]func(*string, configs.WsMessage) configs.WsMessage{
			configs.Baremetal: HandleBaremetalRequest,
			configs.Document:  HandleDocumentRequest,
			configs.Phase:     HandlePhaseRequest,
		})
}

// NewDefaultClient initializes the airshipctl client for external usage with default logging.
func NewDefaultClient(airshipConfigPath, kubeConfigPath *string) (*Client, error) {
	cfgFactory := config.CreateFactory(airshipConfigPath, kubeConfigPath)

	conf, err := cfgFactory()
	if err != nil {
		return nil, err
	}

	client := &Client{
		Config: conf,
	}

	// TODO(mfuller): how do you do this now?
	// set verbosity to true

	return client, nil
}

// NewClient initializes the airshipctl client for external usage with the logging overridden.
func NewClient(airshipConfigPath, kubeConfigPath *string, request configs.WsMessage) (*Client, error) {
	client, err := NewDefaultClient(airshipConfigPath, kubeConfigPath)
	if err != nil {
		return nil, err
	}

	// init the interceptor to send messages to the UI
	// TODO: Unsure how this will be handled with overlapping runs
	log.Init(client.Debug, NewLogInterceptor(request))

	return client, nil
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
	s := string(data)
	response.Message = &s
	if err = webservice.WebSocketSend(response); err != nil {
		uiLog.Errorf("Error receiving / sending message: %s\n", err)
		return len(data), err
	}

	return len(data), nil
}
