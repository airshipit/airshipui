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
	"bytes"
	"text/template"

	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/pkg/version"
	"opendev.org/airship/airshipui/internal/configs"
)

// maintain the state of a potentially long running process
var runningRequests map[configs.WsSubComponentType]bool = make(map[configs.WsSubComponentType]bool)

// ctlPage struct is used for templated HTML
type ctlPage struct {
	ClusterRows    string
	ContextRows    string
	CredentialRows string
	Title          string
	Version        string
	Disabled       string
	ButtonText     string
}

// client provides a library of functions that enable external programs (e.g. Airship UI) to perform airshipctl
// functionality in exactly the same manner as the CLI.
type client struct {
	settings *environment.AirshipCTLSettings
}

// NewClient initializes the airshipctl client for external usage.
func NewClient() *client {
	settings := &environment.AirshipCTLSettings{}
	settings.InitConfig()

	c := &client{
		settings: settings,
	}

	return c
}

// initilize the connection to airshipctl
var c *client = NewClient()

// GetAirshipCTLVersion will kick out what version of airshipctl we're using
func getAirshipCTLVersion() string {
	return version.Get().GitVersion
}

func getHTML(templateFile string, contents ctlPage) (string, error) {
	// go templates need an io writer, since we need a string this buffer can be converted
	var buff bytes.Buffer

	// TODO: make the node path dynamic or setable at compile time
	t, err := template.ParseFiles(templateFile)

	if err != nil {
		return "", err
	}

	// parse and merge the template
	err = template.Must(t, err).Execute(&buff, contents)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}
