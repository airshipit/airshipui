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

package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"opendev.org/airship/airshipctl/pkg/config"
)

// variables related to UI config
var (
	UIConfig Config
)

// Config basic structure to hold configuration params for Airship UI
type Config struct {
	AuthMethod *AuthMethod `json:"authMethod,omitempty"`
	Dashboards []Dashboard `json:"dashboards,omitempty"`
}

// AuthMethod structure to hold authentication parameters
type AuthMethod struct {
	Type  string   `json:"type,omitempty"`
	Value []string `json:"values,omitempty"`
	URL   string   `json:"url,omitempty"`
}

// Executable structure to hold parameters for launching an executable plugin
type Executable struct {
	AutoStart bool     `json:"autoStart,omitempty"`
	Filepath  string   `json:"filepath,omitempty"`
	Args      []string `json:"args,omitempty"`
}

// Dashboard structure
type Dashboard struct {
	Name       string      `json:"name,omitempty"`
	BaseURL    string      `json:"baseURL,omitempty"`
	Path       string      `json:"path,omitempty"`
	IsProxied  bool        `json:"isProxied,omitempty"`
	Executable *Executable `json:"executable,omitempty"`
}

// WsRequestType is used to set the specific types allowable for WsRequests
type WsRequestType string

// WsComponentType is used to set the specific component types allowable for WsRequests
type WsComponentType string

// WsSubComponentType is used to set the specific subcomponent types allowable for WsRequests
type WsSubComponentType string

// constants related to specific request/component/subcomponent types for WsRequests
const (
	AirshipCTL WsRequestType = "airshipctl"
	AirshipUI  WsRequestType = "airshipui"
	Alert      WsRequestType = "alert"

	Authcomplete WsComponentType = "authcomplete"
	Error        WsComponentType = "danger"  // Error corresponds to a red alert message if used as an alert
	Info         WsComponentType = "info"    // Info corresponds to a blue alert message if used as an alert
	Warning      WsComponentType = "warning" // Warning corresponds to an orange alert message if used as an alert
	Success      WsComponentType = "success" // Success corresponds to a green alert message if used as an alert
	SetConfig    WsComponentType = "setConfig"
	Initialize   WsComponentType = "initialize"
	Keepalive    WsComponentType = "keepalive"
	CTLConfig    WsComponentType = "config"
	Baremetal    WsComponentType = "baremetal"
	Document     WsComponentType = "document"

	GetDefaults   WsSubComponentType = "getDefaults"
	SetContext    WsSubComponentType = "context"
	SetCluster    WsSubComponentType = "cluster"
	SetCredential WsSubComponentType = "credential"
	GenerateISO   WsSubComponentType = "generateISO"
	DocPull       WsSubComponentType = "docPull"
	Yaml          WsSubComponentType = "yaml"
	YamlWrite     WsSubComponentType = "yamlWrite"
	GetYaml       WsSubComponentType = "getYaml"
	GetSource     WsSubComponentType = "getSource"
	GetRendered   WsSubComponentType = "getRendered"
)

// WsMessage is a request / return structure used for websockets
type WsMessage struct {
	// base components of a message
	Type         WsRequestType      `json:"type,omitempty"`
	Component    WsComponentType    `json:"component,omitempty"`
	SubComponent WsSubComponentType `json:"subComponent,omitempty"`
	Timestamp    int64              `json:"timestamp,omitempty"`

	// additional conditional components that may or may not be involved in the request / response
	Error           string      `json:"error,omitempty"`
	IsAuthenticated bool        `json:"isAuthenticated,omitempty"`
	Message         string      `json:"message,omitempty"`
	Data            interface{} `json:"data,omitempty"`
	YAML            string      `json:"yaml,omitempty"`
	Name            string      `json:"name,omitempty"`

	// information related to the init of the UI
	Dashboards      []Dashboard             `json:"dashboards,omitempty"`
	Authentication  *AuthMethod             `json:"authentication,omitempty"`
	AuthInfoOptions *config.AuthInfoOptions `json:"authInfoOptions,omitempty"`
	ContextOptions  *config.ContextOptions  `json:"contextOptions,omitempty"`
	ClusterOptions  *config.ClusterOptions  `json:"clusterOptions,omitempty"`
}

// SetUIConfig sets the UIConfig object with values obtained from
// airshipui.json, located at 'filename'
// TODO: add watcher to the json file to reload conf on change
func SetUIConfig(filename string) error {
	bytes, err := getBytesFromFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &UIConfig)
	if err != nil {
		return err
	}

	return nil
}

func getBytesFromFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
