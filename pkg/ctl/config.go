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
	"encoding/json"
	"fmt"

	ctlconfig "opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/log"
)

// HandleConfigRequest will flop between requests so we don't have to have them all mapped as function calls
// This will wait for the sub component to complete before responding.  The assumption is this is an async request
func HandleConfigRequest(user *string, request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.CTLConfig,
		SubComponent: request.SubComponent,
		Name:         request.Name,
	}

	var err error
	var message *string

	client, err := NewClient(AirshipConfigPath, KubeConfigPath, request)
	if err != nil {
		e := fmt.Sprintf("Error initializing airshipctl client: %s", err)
		response.Error = &e
		return response
	}

	subComponent := request.SubComponent
	switch subComponent {
	case configs.GetCurrentContext:
		context := client.Config.CurrentContext
		message = &context
	case configs.GetContexts:
		response.Data = GetContexts(client)
	case configs.GetEncryptionConfigs:
		response.Data = GetEncryptionConfigs(client)
	case configs.GetManagementConfigs:
		response.Data = GetManagementConfigs(client)
	case configs.GetManifests:
		response.Data = GetManifests(client)
	case configs.Init:
		err = InitAirshipConfig(AirshipConfigPath)
	case configs.SetContext:
		response.Data, err = SetContext(client, request)
		str := fmt.Sprintf("Context '%s' has been modified", request.Name)
		message = &str
	case configs.SetEncryptionConfig:
		response.Data, err = SetEncryptionConfig(client, request)
		str := fmt.Sprintf("Encryption configuration '%s' has been modified", request.Name)
		message = &str
	case configs.SetManagementConfig:
		err = SetManagementConfig(client, request)
		str := fmt.Sprintf("Management configuration '%s' has been modified", request.Name)
		message = &str
	case configs.SetManifest:
		response.Data, err = SetManifest(client, request)
		str := fmt.Sprintf("Manifest '%s' has been modified", request.Name)
		message = &str
	case configs.UseContext:
		err = UseContext(client, request)
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

// InitAirshipConfig wrapper function for CTL's CreateConfig using the specified path
func InitAirshipConfig(path *string) error {
	return ctlconfig.CreateConfig(*path)
}

// Context wrapper struct to include context name with CTL's Context
type Context struct {
	Name string `json:"name"`
	ctlconfig.Context
}

// GetContexts returns a slice of wrapper Context structs so we know the name of each
// for display in the UI
func GetContexts(client *Client) []Context {
	contexts := []Context{}
	for name, context := range client.Config.Contexts {
		contexts = append(contexts, Context{
			name,
			ctlconfig.Context{
				NameInKubeconf:          context.NameInKubeconf,
				Manifest:                context.Manifest,
				EncryptionConfig:        context.EncryptionConfig,
				ManagementConfiguration: context.ManagementConfiguration,
			},
		})
	}

	return contexts
}

// Manifest wraps CTL's Manifest to include the manifest name
type Manifest struct {
	Name     string              `json:"name"`
	Manifest *ctlconfig.Manifest `json:"manifest"`
}

// GetManifests returns a slice of wrapper Manifest structs so we know the name of each
// for display in the UI
func GetManifests(client *Client) []Manifest {
	manifests := []Manifest{}

	for name, manifest := range client.Config.Manifests {
		manifests = append(manifests, Manifest{
			Name:     name,
			Manifest: manifest,
		})
	}

	return manifests
}

// ManagementConfig wrapper struct for CTL's ManagementConfiguration that
// includes a name
type ManagementConfig struct {
	Name string `json:"name"`
	ctlconfig.ManagementConfiguration
}

// GetManagementConfigs function to retrieve all management configs
func GetManagementConfigs(client *Client) []ManagementConfig {
	configs := []ManagementConfig{}
	for name, conf := range client.Config.ManagementConfiguration {
		configs = append(configs, ManagementConfig{
			name,
			ctlconfig.ManagementConfiguration{
				Insecure:            conf.Insecure,
				SystemActionRetries: conf.SystemActionRetries,
				SystemRebootDelay:   conf.SystemRebootDelay,
				Type:                conf.Type,
				UseProxy:            conf.UseProxy,
			},
		})
	}
	return configs
}

// EncryptionConfig wrapper struct for CTL's EncryptionConfiguration that
// includes a name
type EncryptionConfig struct {
	Name string `json:"name"`
	ctlconfig.EncryptionConfig
}

// GetEncryptionConfigs returns a slice of wrapper EncryptionConfig structs so we
// know the name of each for display in the UI
func GetEncryptionConfigs(client *Client) []EncryptionConfig {
	configs := []EncryptionConfig{}
	for name, config := range client.Config.EncryptionConfigs {
		configs = append(configs, EncryptionConfig{
			name,
			ctlconfig.EncryptionConfig{
				EncryptionKeyFileSource:   config.EncryptionKeyFileSource,
				EncryptionKeySecretSource: config.EncryptionKeySecretSource,
			},
		})
	}

	return configs
}

// SetContext wrapper function for CTL's RunSetContext, using a UI client
func SetContext(client *Client, message configs.WsMessage) (bool, error) {
	bytes, err := json.Marshal(message.Data)
	if err != nil {
		return false, err
	}

	var opts ctlconfig.ContextOptions
	err = json.Unmarshal(bytes, &opts)
	if err != nil {
		return false, err
	}

	err = opts.Validate()
	if err != nil {
		return false, err
	}

	return ctlconfig.RunSetContext(&opts, client.Config, true)
}

// SetEncryptionConfig wrapper function for CTL's RunSetEncryptionConfig, using a UI client
func SetEncryptionConfig(client *Client, message configs.WsMessage) (bool, error) {
	bytes, err := json.Marshal(message.Data)
	if err != nil {
		return false, err
	}

	var opts ctlconfig.EncryptionConfigOptions
	err = json.Unmarshal(bytes, &opts)
	if err != nil {
		return false, err
	}

	err = opts.Validate()
	if err != nil {
		return false, err
	}

	return ctlconfig.RunSetEncryptionConfig(&opts, client.Config, true)
}

// SetManagementConfig sets the specified management configuration with values
// received from the frontend client
// TODO(mfuller): there's currently no setter for this in the CTL config pkg
// so we'll set the values manually and then persist the config
func SetManagementConfig(client *Client, message configs.WsMessage) error {
	bytes, err := json.Marshal(message.Data)
	if err != nil {
		return err
	}

	if mCfg, found := client.Config.ManagementConfiguration[message.Name]; found {
		err = json.Unmarshal(bytes, mCfg)
		if err != nil {
			return err
		}

		err = client.Config.PersistConfig()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Management configuration '%s' not found", message.Name)
	}

	return nil
}

// SetManifest wrapper function for CTL's RunSetManifest, using a UI client
func SetManifest(client *Client, message configs.WsMessage) (bool, error) {
	bytes, err := json.Marshal(message.Data)
	if err != nil {
		return false, err
	}

	var opts ctlconfig.ManifestOptions
	err = json.Unmarshal(bytes, &opts)
	if err != nil {
		return false, err
	}

	log.Infof("Unmarshaled options: %+v", opts)
	err = opts.Validate()
	if err != nil {
		return false, err
	}

	return ctlconfig.RunSetManifest(&opts, client.Config, true)
}

// UseContext wrapper function for CTL's RunUseConfig, using a UI client
func UseContext(client *Client, message configs.WsMessage) error {
	return ctlconfig.RunUseContext(message.Name, client.Config)
}
