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
package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// AirshipuiProps basic structure for a given external dashboard
// TODO: solidify the struct requirements for the input
// TODO: maybe move where props gathering and parsing lives
type AirshipuiProps struct {
	AuthMethod struct {
		Type  string   `json:"type,omitempty"`
		Value []string `json:"values,omitempty"`
		URL   string   `json:"url,omitempty"`
	} `json:"authMethod"`
	ExtDashboard []interface{} `json:"external_dashboards"`
}

// AirshipuiPropsCache the file so we don't have to reread every execution
// TODO: maybe move where props gathering and parsing lives
var AirshipuiPropsCache AirshipuiProps

// TODO: add watcher to the json file to reload conf on change
// TODO: maybe move where props gathering and parsing lives
// Get dashboard info if present in $HOME/.airshipui/airshipui.json
func GetConfsFromFile() error {
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

	err = json.Unmarshal(byteValue, &AirshipuiPropsCache)

	if err != nil {
		return err
	}
	return err
}
