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
	"path/filepath"

	"opendev.org/airship/airshipctl/pkg/document"
	"opendev.org/airship/airshipctl/pkg/document/pull"
	"opendev.org/airship/airshipui/pkg/configs"
)

var (
	index map[string]interface{}
)

// HandleDocumentRequest will flop between requests so we don't have to have them all mapped as function calls
func HandleDocumentRequest(request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.Document,
		SubComponent: request.SubComponent,
	}

	var err error
	var message string
	var id string
	switch request.SubComponent {
	case configs.DocPull:
		message, err = c.docPull()
	case configs.YamlWrite:
		id = request.ID
		response.Name, response.YAML, err = writeYamlFile(id, request.YAML)
		message = fmt.Sprintf("File '%s' saved successfully", response.Name)
	case configs.GetYaml:
		id = request.ID
		response.Name, response.YAML, err = getYaml(id)
	case configs.GetPhaseTree:
		response.Data, err = GetPhaseTree()
	case configs.GetPhaseDocuments:
		id = request.ID
		response.Data, err = GetPhaseDocuments(request.ID)
	case configs.GetPhaseSourceFiles:
		id = request.ID
		response.Data, err = GetPhaseSourceFiles(request.ID)
	case configs.GetTarget:
		message = getTarget()
	default:
		err = fmt.Errorf("Subcomponent %s not found", request.SubComponent)
	}

	if err != nil {
		response.Error = err.Error()
	} else {
		response.Message = message
		response.ID = id
	}

	return response
}

func getTarget() string {
	m, err := c.settings.Config.CurrentContextManifest()
	if err != nil {
		return "unknown"
	}

	return filepath.Join(m.TargetPath, m.SubPath)
}

func getYaml(id string) (string, string, error) {
	obj := index[id]
	switch t := obj.(type) {
	case string:
		return getFileYaml(t)
	case document.Document:
		return getDocumentYaml(t)
	default:
		return "", "", fmt.Errorf("ID %s not found in index", id)
	}
}

func getDocumentYaml(doc document.Document) (string, string, error) {
	title := doc.GetName()
	bytes, err := doc.AsYAML()
	if err != nil {
		return "", "", err
	}

	return title, base64.StdEncoding.EncodeToString(bytes), nil
}

func getFileYaml(path string) (string, string, error) {
	ccm, err := c.settings.Config.CurrentContextManifest()
	if err != nil {
		return "", "", err
	}

	sitePath := filepath.Join(ccm.TargetPath, ccm.SubPath)

	// this path is making the assumption that the subPath
	// in airship config is going to be pointing to a site
	// in airshipctl/manifests/site/{SITENAME}
	manifestsDir := filepath.Join(sitePath, "..", "..")

	title, err := filepath.Rel(manifestsDir, path)
	if err != nil {
		return "", "", err
	}

	file, err := os.Open(path)
	if err != nil {
		return "", "", err
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", "", err
	}

	return title, base64.StdEncoding.EncodeToString(bytes), nil
}

func writeYamlFile(id, yaml64 string) (string, string, error) {
	path, ok := index[id].(string)
	if !ok {
		return "", "", fmt.Errorf("ID %s not found", id)
	}

	yaml, err := base64.StdEncoding.DecodeString(yaml64)
	if err != nil {
		return "", "", err
	}

	err = ioutil.WriteFile(path, yaml, 0600)
	if err != nil {
		return "", "", err
	}

	return getFileYaml(path)
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
