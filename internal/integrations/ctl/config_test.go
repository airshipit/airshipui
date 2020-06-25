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
	"log"
	"testing"

	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipui/internal/configs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"opendev.org/airship/airshipui/testutil"
)

// TODO: Determine if this should be broken out into it's own file
// setup the airshipCTL env prior to running
func initCTL(t *testing.T) {
	conf, configPath, kubeConfigPath, cleanup := testutil.InitConfig(t)
	defer cleanup(t)
	// point airshipctl client toward test configs
	c.settings = &environment.AirshipCTLSettings{
		AirshipConfigPath: configPath,
		KubeConfigPath:    kubeConfigPath,
		Config:            conf,
	}

	err := c.settings.Config.LoadConfig(
		c.settings.AirshipConfigPath,
		c.settings.KubeConfigPath,
	)

	if err != nil {
		log.Fatal(err)
	}
}

func TestHandleDefaultConfigRequest(t *testing.T) {
	initCTL(t)
	// get the default html
	html, err := getConfigHTML()
	require.NoError(t, err)

	// simulate incoming WsMessage from websocket client
	request := configs.WsMessage{
		Type:         configs.AirshipCTL,
		Component:    configs.CTLConfig,
		SubComponent: configs.GetDefaults,
	}

	response := HandleConfigRequest(request)

	expected := configs.WsMessage{
		Type:         configs.AirshipCTL,
		Component:    configs.CTLConfig,
		SubComponent: configs.GetDefaults,
		HTML:         html,
	}

	assert.Equal(t, expected, response)
}

func TestHandleUnknownConfigSubComponent(t *testing.T) {
	request := configs.WsMessage{
		Type:         configs.AirshipCTL,
		Component:    configs.CTLConfig,
		SubComponent: "fake_subcomponent",
	}

	response := HandleConfigRequest(request)

	expected := configs.WsMessage{
		Type:         configs.AirshipCTL,
		Component:    configs.CTLConfig,
		SubComponent: "fake_subcomponent",
		Error:        "Subcomponent fake_subcomponent not found",
	}

	assert.Equal(t, expected, response)
}
