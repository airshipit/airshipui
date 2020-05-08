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
)

var (
	UiConfig Config
)

// Config basic structure to hold configuration params for Airship UI
type Config struct {
	AuthMethod struct {
		Type  string   `json:"type,omitempty"`
		Value []string `json:"values,omitempty"`
		URL   string   `json:"url,omitempty"`
	} `json:"authMethod"`
	Plugins  []Plugin  `json:"plugins"`
	Clusters []Cluster `json:"clusters"`
}

type Plugin struct {
	Name      string `json:"name"`
	Dashboard struct {
		Protocol string `json:"protocol"`
		FQDN     string `json:"fqdn"`
		Port     uint16 `json:"port"`
		Path     string `json:"path"`
	} `json:"dashboard"`
	Executable struct {
		AutoStart bool     `json:"autoStart"`
		Filepath  string   `json:"filepath"`
		Args      []string `json:"args"`
	} `json:"executable"`
}

// Dashboard structure
type Dashboard struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Hostname string `json:"hostname,omitempty"`
	FQDN     string `json:"fqdn,omitempty"`
	Port     uint16 `json:"port"`
	Path     string `json:"path"`
}

// Namespace structure
type Namespace struct {
	Name       string      `json:"name"`
	Dashboards []Dashboard `json:"dashboards"`
}

// Cluster basic structure describing a cluster
type Cluster struct {
	Name       string      `json:"name"`
	BaseFqdn   string      `json:"baseFqdn"`
	Namespaces []Namespace `json:"namespaces"`
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
