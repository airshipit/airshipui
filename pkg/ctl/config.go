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
	"os"
	"path/filepath"

	ctlconfig "opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipui/pkg/configs"
)

// ConfigFunctionMap is being used to call the appropriate function based on the SubComponentType,
// since the linter seems to think there are too many cases for a switch / case
var ConfigFunctionMap = map[configs.WsSubComponentType]func(configs.WsMessage) configs.WsMessage{
	configs.SetAirshipConfig:     SetAirshipConfig,
	configs.GetAirshipConfigPath: GetAirshipConfigPath,
	configs.GetCurrentContext:    GetCurrentContext,
	configs.GetContexts:          GetContexts,
	configs.GetEncryptionConfigs: GetEncryptionConfigs,
	configs.GetManagementConfigs: GetManagementConfigs,
	configs.GetManifests:         GetManifests,
	configs.Init:                 InitAirshipConfig,
	configs.SetContext:           SetContext,
	configs.SetEncryptionConfig:  SetEncryptionConfig,
	configs.SetManagementConfig:  SetManagementConfig,
	configs.SetManifest:          SetManifest,
	configs.UseContext:           UseContext,
}

// helper function to create most of the relevant bits of the response message
func newResponse(request configs.WsMessage) configs.WsMessage {
	return configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.CTLConfig,
		SubComponent: request.SubComponent,
		Name:         request.Name,
	}
}

// HandleConfigRequest will find the appropriate subcomponent function in the function map
// and wait for it to complete before returning the response message
func HandleConfigRequest(user *string, request configs.WsMessage) configs.WsMessage {
	var response configs.WsMessage

	if handler, ok := ConfigFunctionMap[request.SubComponent]; ok {
		response = handler(request)
	} else {
		response = newResponse(request)
		err := fmt.Sprintf("Subcomponent %s not found", request.SubComponent)
		response.Error = &err
	}

	return response
}

// GetAirshipConfigPath returns value stored in AirshipConfigPath
func GetAirshipConfigPath(request configs.WsMessage) configs.WsMessage {
	response := newResponse(request)

	// leave message empty if the file doesn't exist
	if configFileExists(configs.UIConfig.AirshipConfigPath) {
		response.Message = configs.UIConfig.AirshipConfigPath
	}

	return response
}

// SetAirshipConfig sets the AirshipConfigPath to the value specified by
// UI client
func SetAirshipConfig(request configs.WsMessage) configs.WsMessage {
	response := newResponse(request)

	configs.UIConfig.AirshipConfigPath = request.Message
	err := configs.UIConfig.Persist()
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	msg := fmt.Sprintf("Config file set to '%s'", *configs.UIConfig.AirshipConfigPath)
	response.Message = &msg

	return response
}

// GetCurrentContext returns the name of the currently configured context
func GetCurrentContext(request configs.WsMessage) configs.WsMessage {
	response := newResponse(request)

	client, err := NewClient(configs.UIConfig.AirshipConfigPath, request)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	response.Message = &client.Config.CurrentContext

	return response
}

// InitAirshipConfig wrapper function for CTL's CreateConfig using the specified path
// TODO(mfuller): we'll need to persist this info in airshipui.json so that we can
// set AirshipConfigPath at app launch
func InitAirshipConfig(request configs.WsMessage) configs.WsMessage {
	response := newResponse(request)

	confPath := *request.Message
	if confPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			e := err.Error()
			response.Error = &e
			return response
		}
		confPath = filepath.Join(home, ".airship", "config")
	}

	err := ctlconfig.CreateConfig(confPath)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	configs.UIConfig.AirshipConfigPath = &confPath
	// save this location back to airshipui config file so we'll remember it for next time
	err = configs.UIConfig.Persist()
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	msg := fmt.Sprintf("Config file set to '%s'", *configs.UIConfig.AirshipConfigPath)

	response.Message = &msg

	return response
}

// Context wrapper struct to include context name with CTL's Context
type Context struct {
	Name string `json:"name"`
	ctlconfig.Context
}

// GetContexts returns a slice of wrapper Context structs so we know the name of each
// for display in the UI
func GetContexts(request configs.WsMessage) configs.WsMessage {
	response := newResponse(request)

	client, err := NewClient(configs.UIConfig.AirshipConfigPath, request)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

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

	response.Data = contexts

	return response
}

// Manifest wraps CTL's Manifest to include the manifest name
type Manifest struct {
	Name     string              `json:"name"`
	Manifest *ctlconfig.Manifest `json:"manifest"`
}

// GetManifests returns a slice of wrapper Manifest structs so we know the name of each
// for display in the UI
func GetManifests(request configs.WsMessage) configs.WsMessage {
	response := newResponse(request)

	client, err := NewClient(configs.UIConfig.AirshipConfigPath, request)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	manifests := []Manifest{}

	for name, manifest := range client.Config.Manifests {
		manifests = append(manifests, Manifest{
			Name:     name,
			Manifest: manifest,
		})
	}

	response.Data = manifests

	return response
}

// ManagementConfig wrapper struct for CTL's ManagementConfiguration that
// includes a name
type ManagementConfig struct {
	Name string `json:"name"`
	ctlconfig.ManagementConfiguration
}

// GetManagementConfigs function to retrieve all management configs
func GetManagementConfigs(request configs.WsMessage) configs.WsMessage {
	response := newResponse(request)

	client, err := NewClient(configs.UIConfig.AirshipConfigPath, request)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

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

	response.Data = configs

	return response
}

// EncryptionConfig wrapper struct for CTL's EncryptionConfiguration that
// includes a name
type EncryptionConfig struct {
	Name string `json:"name"`
	ctlconfig.EncryptionConfig
}

// GetEncryptionConfigs returns a slice of wrapper EncryptionConfig structs so we
// know the name of each for display in the UI
func GetEncryptionConfigs(request configs.WsMessage) configs.WsMessage {
	response := newResponse(request)

	client, err := NewClient(configs.UIConfig.AirshipConfigPath, request)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

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

	response.Data = configs

	return response
}

// SetContext wrapper function for CTL's RunSetContext, using a UI client
func SetContext(request configs.WsMessage) configs.WsMessage {
	response := newResponse(request)

	client, err := NewClient(configs.UIConfig.AirshipConfigPath, request)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	bytes, err := json.Marshal(request.Data)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	var opts ctlconfig.ContextOptions
	err = json.Unmarshal(bytes, &opts)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	err = opts.Validate()
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	_, err = ctlconfig.RunSetContext(&opts, client.Config, true)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	msg := fmt.Sprintf("Context '%s' has been modified", request.Name)
	response.Message = &msg

	return response
}

// SetEncryptionConfig wrapper function for CTL's RunSetEncryptionConfig, using a UI client
func SetEncryptionConfig(request configs.WsMessage) configs.WsMessage {
	response := newResponse(request)

	client, err := NewClient(configs.UIConfig.AirshipConfigPath, request)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	bytes, err := json.Marshal(request.Data)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	var opts ctlconfig.EncryptionConfigOptions
	err = json.Unmarshal(bytes, &opts)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	err = opts.Validate()
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	_, err = ctlconfig.RunSetEncryptionConfig(&opts, client.Config, true)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	msg := fmt.Sprintf("Encryption configuration '%s' has been modified", request.Name)
	response.Message = &msg

	return response
}

// SetManagementConfig sets the specified management configuration with values
// received from the frontend client
// TODO(mfuller): there's currently no setter for this in the CTL config pkg
// so we'll set the values manually and then persist the config
func SetManagementConfig(request configs.WsMessage) configs.WsMessage {
	response := newResponse(request)

	client, err := NewClient(configs.UIConfig.AirshipConfigPath, request)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	bytes, err := json.Marshal(request.Data)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	if mCfg, found := client.Config.ManagementConfiguration[request.Name]; found {
		err = json.Unmarshal(bytes, mCfg)
		if err != nil {
			e := err.Error()
			response.Error = &e
			return response
		}

		err = client.Config.PersistConfig()
		if err != nil {
			e := err.Error()
			response.Error = &e
			return response
		}
	} else {
		e := fmt.Sprintf("Management configuration '%s' not found", request.Name)
		response.Error = &e
		return response
	}

	msg := fmt.Sprintf("Management configuration '%s' has been modified", request.Name)
	response.Message = &msg

	return response
}

// SetManifest wrapper function for CTL's RunSetManifest, using a UI client
func SetManifest(request configs.WsMessage) configs.WsMessage {
	response := newResponse(request)

	client, err := NewClient(configs.UIConfig.AirshipConfigPath, request)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	bytes, err := json.Marshal(request.Data)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	var opts ctlconfig.ManifestOptions
	err = json.Unmarshal(bytes, &opts)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	err = opts.Validate()
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	_, err = ctlconfig.RunSetManifest(&opts, client.Config, true)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	msg := fmt.Sprintf("Manifest '%s' has been modified", request.Name)
	response.Message = &msg

	return response
}

// UseContext wrapper function for CTL's RunUseConfig, using a UI client
func UseContext(request configs.WsMessage) configs.WsMessage {
	response := newResponse(request)

	client, err := NewClient(configs.UIConfig.AirshipConfigPath, request)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	err = ctlconfig.RunUseContext(request.Name, client.Config)
	if err != nil {
		e := err.Error()
		response.Error = &e
		return response
	}

	msg := fmt.Sprintf("Using context '%s'", request.Name)
	response.Message = &msg

	return response
}
