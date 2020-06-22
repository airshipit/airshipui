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
	"io/ioutil"
	"testing"

	"opendev.org/airship/airshipui/internal/configs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipctl/pkg/environment"
)

const (
	testConfigHTML    string = "testdata/config.html"
	testKubeConfig    string = "testdata/kubeconfig.yaml"
	testAirshipConfig string = "testdata/config.yaml"
)

func TestHandleDefaultConfigRequest(t *testing.T) {
	html, err := ioutil.ReadFile(testConfigHTML)
	require.NoError(t, err)

	// point airshipctl client toward test configs
	c.settings = &environment.AirshipCTLSettings{
		AirshipConfigPath: testAirshipConfig,
		KubeConfigPath:    testKubeConfig,
		Config:            config.NewConfig(),
	}

	err = c.settings.Config.LoadConfig(
		c.settings.AirshipConfigPath,
		c.settings.KubeConfigPath,
	)
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
		HTML:         string(html),
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
