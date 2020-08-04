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
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"opendev.org/airship/airshipctl/pkg/document/pull"
	"opendev.org/airship/airshipui/pkg/configs"
)

// HandleDocumentRequest will flop between requests so we don't have to have them all mapped as function calls
func HandleDocumentRequest(request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.AirshipCTL,
		Component:    configs.Document,
		SubComponent: request.SubComponent,
	}

	var err error
	var message string
	switch request.SubComponent {
	case configs.GetDefaults:
		response.Data = getGraphData()
	case configs.DocPull:
		message, err = c.docPull()
	case configs.Yaml:
		message = request.Message
		response.YAML, err = getYaml(message)
	case configs.YamlWrite:
		message = request.Message
		response.YAML, err = writeYaml(message, request.YAML)
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

// network graphs have nodes and edges defined, just attempting to put some dynamically defined data in it
func getGraphData() map[string]interface{} {
	return map[string]interface{}{
		"nodes": []map[string]string{
			{"id": "1", "label": ".airshipui"},
			{"id": "2", "label": c.settings.KubeConfigPath},
			{"id": "3", "label": c.settings.AirshipConfigPath},
		},
		"edges": []map[string]int64{
			{"from": 1, "to": 2},
			{"from": 1, "to": 3},
		},
	}
}

// getYaml reads the requested file and returns base64 encoded yaml for the front end to render
func getYaml(yamlType string) (string, error) {
	yamlFile, err := os.Open(getYamlFile(yamlType))
	if err != nil {
		return "", err
	}

	defer yamlFile.Close()

	// TODO: determine if this needs to be parsed as YAML as a validation effort
	bytes, err := ioutil.ReadAll(yamlFile)
	return base64.StdEncoding.EncodeToString(bytes), err
}

// a way to do a sanity check on the yaml passed from the frontend
func writeYaml(yamlType string, yaml64 string) (string, error) {
	// base64 decode
	yaml, err := base64.StdEncoding.DecodeString(yaml64)
	if err != nil {
		return "", err
	}

	// TODO: determine if we need to backup the existing before overwrite
	err = ioutil.WriteFile(getYamlFile(yamlType), yaml, 0600)
	if err != nil {
		return "", err
	}

	return getYaml(yamlType)
}

func getYamlFile(yamlType string) string {
	var fileName string
	switch yamlType {
	case "kube":
		fileName = c.settings.KubeConfigPath
	case "airship":
		fileName = c.settings.AirshipConfigPath
	}

	return fileName
}

func (c *Client) docPull() (string, error) {
	var message string
	settings := pull.Settings{AirshipCTLSettings: c.settings}
	err := settings.Pull()
	if err == nil {
		message = fmt.Sprintf("Success")
	}

	return message, err
}
