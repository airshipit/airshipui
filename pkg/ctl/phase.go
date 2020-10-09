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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"opendev.org/airship/airshipctl/pkg/document"
	"opendev.org/airship/airshipctl/pkg/events"
	"opendev.org/airship/airshipctl/pkg/phase"
	"opendev.org/airship/airshipctl/pkg/phase/ifc"
	"opendev.org/airship/airshipui/pkg/configs"
)

var (
	fileIndex map[string]string
	docIndex  map[string]document.Document
)

// HandlePhaseRequest will flop between requests so we don't have to have them all mapped as function calls
// This will wait for the sub component to complete before responding.  The assumption is this is an async request
func HandlePhaseRequest(user *string, request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.Phase,
		SubComponent: request.SubComponent,
	}

	var err error
	var message *string
	var id string
	var valid bool

	client, err := NewClient(AirshipConfigPath, KubeConfigPath, request)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	switch request.SubComponent {
	case configs.Run:
		err = client.RunPhase(request)
	case configs.ValidatePhase:
		valid, err = client.ValidatePhase(request.ID, request.SessionID)
		message = validateHelper(valid)
	case configs.YamlWrite:
		id = request.ID
		response.Name, response.YAML, err = client.writeYamlFile(id, request.YAML)
		s := fmt.Sprintf("File '%s' saved successfully", response.Name)
		message = &s
	case configs.GetYaml:
		id = request.ID
		message = request.Message
		response.Name, response.YAML, err = client.getYaml(id, *message)
	case configs.GetPhaseTree:
		response.Data, err = client.GetPhaseTree()
	case configs.GetPhase:
		id = request.ID
		s := "rendered"
		message = &s
		response.Name, response.Details, response.YAML, err = client.GetPhase(id)
	case configs.GetDocumentsBySelector:
		id = request.ID
		message = request.Message
		response.Data, err = GetDocumentsBySelector(request.ID, *message)
	case configs.GetTarget:
		message = client.getTarget()
	case configs.GetExecutorDoc:
		id = request.ID
		s := "rendered"
		message = &s
		response.Name, response.YAML, err = client.GetExecutorDoc(id)
	default:
		err = fmt.Errorf("Subcomponent %s not found", request.SubComponent)
	}

	if err != nil {
		e := err.Error()
		response.Error = &e
	} else {
		response.Message = message
		response.ID = id
	}

	return response
}

// this helper function will likely disappear once a clear workflow for
// phase validation takes shape in UI. For now, it simply returns a
// string message to be displayed as a toast in frontend client
func validateHelper(valid bool) *string {
	msg := "invalid"
	if valid {
		msg = "valid"
	}
	return &msg
}

// GetExecutorDoc returns the title and YAML of the executor document for the specified phase
func (c *Client) GetExecutorDoc(id string) (string, string, error) {
	helper, err := getHelper()
	if err != nil {
		return "", "", err
	}

	phaseID := ifc.ID{}

	err = json.Unmarshal([]byte(id), &phaseID)
	if err != nil {
		return "", "", err
	}

	ed, err := helper.ExecutorDoc(phaseID)
	if err != nil {
		return "", "", err
	}

	title := ed.GetName()
	bytes, err := ed.AsYAML()
	if err != nil {
		return "", "", err
	}

	return title, base64.StdEncoding.EncodeToString(bytes), nil
}

func (c *Client) getTarget() *string {
	var s string
	m, err := c.Config.CurrentContextManifest()
	if err != nil {
		s = "unknown"
		return &s
	}

	s = filepath.Join(m.TargetPath, m.SubPath)
	return &s
}

func (c *Client) getPhaseDetails(id ifc.ID) (string, error) {
	helper, err := getHelper()
	if err != nil {
		return "", err
	}

	pClient := phase.NewClient(helper)

	phase, err := pClient.PhaseByID(id)
	if err != nil {
		return "", err
	}

	return phase.Details()
}

func (c *Client) getYaml(id, message string) (string, string, error) {
	switch message {
	case "source":
		name, yaml, err := c.getFileYaml(id)
		return name, yaml, err
	case "rendered":
		name, yaml, err := c.getDocumentYaml(id)
		return name, yaml, err
	default:
		return "", "", fmt.Errorf("'%s' unrecognized document type", message)
	}
}

func (c *Client) getDocumentYaml(id string) (string, string, error) {
	doc, ok := docIndex[id]
	if !ok {
		return "", "", fmt.Errorf("document with ID '%s' not found", id)
	}
	title := doc.GetName()
	bytes, err := doc.AsYAML()
	if err != nil {
		return "", "", err
	}

	return title, base64.StdEncoding.EncodeToString(bytes), nil
}

func (c *Client) getFileYaml(id string) (string, string, error) {
	path, ok := fileIndex[id]
	if !ok {
		return "", "", fmt.Errorf("file with ID '%s' not found", id)
	}

	ccm, err := c.Config.CurrentContextManifest()
	if err != nil {
		return "", "", err
	}

	// this is making the assumption that the site definition
	// will always found at: targetPath/subPath
	sitePath := filepath.Join(ccm.TargetPath, ccm.SubPath)

	// TODO(mfuller): will this be true in treasuremap or
	// other external repos?
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

func (c *Client) writeYamlFile(id, yaml64 string) (string, string, error) {
	path, ok := fileIndex[id]
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

	return c.getFileYaml(id)
}

func getPhaseBundle(id ifc.ID) (document.Bundle, error) {
	helper, err := getHelper()
	if err != nil {
		return nil, err
	}

	pClient := phase.NewClient(helper)

	phase, err := pClient.PhaseByID(id)
	if err != nil {
		return nil, err
	}

	docRoot, err := phase.DocumentRoot()
	if err != nil {
		// if phase has no doc entrypoint defined, just
		// return nothing; otherwise, return the error
		if strings.Contains(err.Error(), "defined") {
			return nil, nil
		}
		return nil, err
	}

	b, err := document.NewBundleByPath(docRoot)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetPhase returns the name, description, and doc bundle for specified phase
func (c *Client) GetPhase(id string) (string, string, string, error) {
	phaseID := ifc.ID{}

	err := json.Unmarshal([]byte(id), &phaseID)
	if err != nil {
		return "", "", "", err
	}

	title := phaseID.Name

	details, err := c.getPhaseDetails(phaseID)
	if err != nil {
		return "", "", "", err
	}

	bundle, err := getPhaseBundle(phaseID)
	if err != nil {
		return "", "", "", err
	}

	// only return title if phase has no bundle
	if bundle == nil {
		return title, details, "", nil
	}

	var buf bytes.Buffer
	err = bundle.Write(&buf)
	if err != nil {
		return "", "", "", err
	}

	return title, details, base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// SelectorParams structure to hold data for constructing a document Selector
type SelectorParams struct {
	Name       string `json:"name,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	GVK        GVK    `json:"gvk,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Label      string `json:"label,omitempty"`
	Annotation string `json:"annotation,omitempty"`
}

// GVK small structure to hold group, version, kind for building a Selector
type GVK struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

// GetDocumentsBySelector returns a slice of KustomNodes representing all phase
// documents returned by applying the provided Selector
func GetDocumentsBySelector(id string, data string) ([]KustomNode, error) {
	docIndex = map[string]document.Document{}

	selector, err := getSelector(data)
	if err != nil {
		return nil, err
	}

	phaseID := ifc.ID{}
	err = json.Unmarshal([]byte(id), &phaseID)
	if err != nil {
		return nil, err
	}

	helper, err := getHelper()
	if err != nil {
		return nil, err
	}

	pClient := phase.NewClient(helper)

	phase, err := pClient.PhaseByID(phaseID)
	if err != nil {
		return nil, err
	}

	docRoot, err := phase.DocumentRoot()
	if err != nil {
		return nil, err
	}

	bundle, err := document.NewBundleByPath(docRoot)
	if err != nil {
		return nil, err
	}

	docs, err := bundle.Select(selector)
	if err != nil {
		return nil, err
	}

	results := []KustomNode{}

	for _, doc := range docs {
		// this is a workaround for a kustomize issue where cluster-scoped objects
		// are included in matching results when a namespace selector is specified
		// (https://github.com/kubernetes-sigs/kustomize/issues/2248)
		if selector.Namespace != "" && selector.Namespace != doc.GetNamespace() {
			continue
		}

		id := uuid.New().String()
		docIndex[id] = doc

		name := doc.GetNamespace()
		if name == "" {
			name = "[none]"
		}

		results = append(results, KustomNode{
			ID: id,
			Name: fmt.Sprintf("%s/%s/%s",
				name,
				doc.GetKind(),
				doc.GetName(),
			)},
		)
	}

	return results, nil
}

func getSelector(data string) (document.Selector, error) {
	params := SelectorParams{}
	err := json.Unmarshal([]byte(data), &params)
	if err != nil {
		return document.Selector{}, err
	}

	s := document.NewSelector()

	// build selector based on what we were given
	if params.Name != "" {
		s = s.ByName(params.Name)
	}
	if params.Namespace != "" {
		s = s.ByNamespace(params.Namespace)
	}
	if (GVK{}) != params.GVK {
		s = s.ByGvk(
			params.GVK.Group,
			params.GVK.Version,
			params.GVK.Kind,
		)
	}
	if params.Kind != "" {
		s = s.ByKind(params.Kind)
	}
	if params.Label != "" {
		s = s.ByLabel(params.Label)
	}
	if params.Annotation != "" {
		s = s.ByAnnotation(params.Annotation)
	}
	return s, nil
}

// ValidatePhase validates the specified phase
// (ifc.Phase.Validate isn't implemented yet, so this function
// currently always returns "valid")
func (c *Client) ValidatePhase(id, sessionID string) (bool, error) {
	phase, err := getPhaseIfc(id, sessionID)
	if err != nil {
		return false, err
	}

	err = phase.Validate()
	if err != nil {
		return false, err
	}

	return true, nil
}

// RunPhase runs the selected phase
func (c *Client) RunPhase(request configs.WsMessage) error {
	phase, err := getPhaseIfc(request.ID, request.SessionID)
	if err != nil {
		return err
	}

	opts := ifc.RunOptions{}
	if request.Data != nil {
		bytes, err := json.Marshal(request.Data)
		if err != nil {
			return err
		}

		err = json.Unmarshal(bytes, &opts)
		if err != nil {
			return err
		}
	}

	return phase.Run(opts)
}

// helper function to return a Phase interface based on a JSON
// string representation of an ifc.ID value
func getPhaseIfc(id, sessionID string) (ifc.Phase, error) {
	phaseID := ifc.ID{}

	err := json.Unmarshal([]byte(id), &phaseID)
	if err != nil {
		return nil, err
	}

	helper, err := getHelper()
	if err != nil {
		return nil, err
	}

	var procFunc phase.ProcessorFunc
	procFunc = func() events.EventProcessor {
		return NewUIEventProcessor(sessionID)
	}

	// inject event processor to phase client
	proc := phase.InjectProcessor(procFunc)

	client := phase.NewClient(helper, proc)

	phase, err := client.PhaseByID(phaseID)
	if err != nil {
		return nil, err
	}

	return phase, nil
}
