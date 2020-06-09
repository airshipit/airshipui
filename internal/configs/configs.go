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
	"path/filepath"

	"opendev.org/airship/airshipctl/pkg/config"
)

var (
	UiConfig Config
)

// Config basic structure to hold configuration params for Airship UI
type Config struct {
	AuthMethod AuthMethod `json:"authMethod,omitempty"`
	Plugins    []Plugin   `json:"plugins,omitempty"`
	Clusters   []Cluster  `json:"clusters,omitempty"`
}

type AuthMethod struct {
	Type  string   `json:"type,omitempty"`
	Value []string `json:"values,omitempty"`
	URL   string   `json:"url,omitempty"`
}

type Plugin struct {
	Name      string `json:"name,omitempty"`
	Dashboard struct {
		Protocol string `json:"protocol,omitempty"`
		FQDN     string `json:"fqdn,omitempty"`
		Port     uint16 `json:"port,omitempty"`
		Path     string `json:"path,omitempty"`
	} `json:"dashboard"`
	Executable struct {
		AutoStart bool     `json:"autoStart,omitempty"`
		Filepath  string   `json:"filepath,omitempty"`
		Args      []string `json:"args,omitempty"`
	} `json:"executable"`
}

// Dashboard structure
type Dashboard struct {
	Name     string `json:"name,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	FQDN     string `json:"fqdn,omitempty"`
	Port     uint16 `json:"port,omitempty"`
	Path     string `json:"path,omitempty"`
}

// Namespace structure
type Namespace struct {
	Name       string      `json:"name,omitempty"`
	Dashboards []Dashboard `json:"dashboards,omitempty"`
}

// Cluster basic structure describing a cluster
type Cluster struct {
	Name       string      `json:"name,omitempty"`
	BaseFqdn   string      `json:"baseFqdn,omitempty"`
	Namespaces []Namespace `json:"namespaces,omitempty"`
}

// WsRequestType is used to set the specific types allowable for WsRequests
type WsRequestType string

// WsComponentType is used to set the specific component types allowable for WsRequests
type WsComponentType string

// WsSubComponentType is used to set the specific subcomponent types allowable for WsRequests
type WsSubComponentType string

const (
	AirshipCTL WsRequestType = "airshipctl"
	Electron   WsRequestType = "electron"
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
)

// WsMessage is a request / return structure used for websockets
type WsMessage struct {
	// base components of a message
	Type         WsRequestType      `json:"type,omitempty"`
	Component    WsComponentType    `json:"component,omitempty"`
	SubComponent WsSubComponentType `json:"subComponent,omitempty"`
	Timestamp    int64              `json:"timestamp,omitempty"`

	// additional conditional components that may or may not be involved in the request / response
	Error           string `json:"error,omitempty"`
	Fade            bool   `json:"fade,omitempty"`
	HTML            string `json:"html,omitempty"`
	IsAuthenticated bool   `json:"isAuthenticated,omitempty"`
	Message         string `json:"message,omitempty"`

	// information related to the init of the UI
	Dashboards      []Cluster              `json:"dashboards,omitempty"`
	Plugins         []Plugin               `json:"plugins,omitempty"`
	Authentication  AuthMethod             `json:"authentication,omitempty"`
	AuthInfoOptions config.AuthInfoOptions `json:"authInfoOptions,omitempty"`
	ContextOptions  config.ContextOptions  `json:"contextOptions,omitempty"`
	ClusterOptions  config.ClusterOptions  `json:"clusterOptions,omitempty"`
}

// TODO: add watcher to the json file to reload conf on change
func GetConfigFromFile() error {
	var fileName string
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	fileName = filepath.FromSlash(home + "/.airship/airshipui.json")

	jsonFile, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &UiConfig)

	if err != nil {
		return err
	}

	return err
}
