/*
 Copyright (c) 2020 AT&T. All Rights Reserved.

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
package webservice

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// basic structure for a given external dashboard
// TODO: solidify the struct requirements for the input
type extPlugins struct {
	ExtDashboard []interface{} `json:"external_dashboards"`
}

// cache the file so we don't have to reread every execution
var pluginCache map[string]interface{}

// getPlugins updates the pluginCache from file if needed
func getPlugins() map[string]interface{} {
	if pluginCache == nil {
		err := getPluginsFromFile()
		if err != nil {
			log.Printf("Error attempting to get plugins from file: %s\n", err)
		}
	}
	return pluginCache
}

// TODO: add watcher to the json file to reload conf on change
// Get dashboard links for Plugins if present in $HOME/.airshipui/plugins.json
func getPluginsFromFile() error {
	var fileName string
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error determining the home directory %s\n", err)
	}

	fileName = filepath.FromSlash(home + "/.airship/plugins.json")

	jsonFile, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return err
	}

	var plugins extPlugins
	err = json.Unmarshal(byteValue, &plugins)

	if err != nil {
		return err
	}

	log.Printf("Plugins found: %v\n", plugins)

	pluginCache = map[string]interface{}{
		"type":      "plugins",
		"component": "dropdown",
		"timestamp": time.Now().UnixNano() / 1000000,
		"plugins":   plugins,
	}
	return err
}
