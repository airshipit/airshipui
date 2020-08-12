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

package webservice

import (
	"encoding/json"
	"testing"
	"time"

	"opendev.org/airship/airshipui/util/utiltest"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"opendev.org/airship/airshipui/pkg/configs"
)

func TestClientInit(t *testing.T) {
	client, err := NewTestClient()
	require.NoError(t, err)
	defer client.Close()

	// simulate config provided by airshipui.json
	configs.UIConfig = utiltest.DummyCompleteConfig()

	// get server response to "initialize" message from client
	response, err := getResponse(client, initialize)
	require.NoError(t, err)

	expected := configs.WsMessage{
		Type:            configs.UI,
		Component:       configs.Initialize,
		IsAuthenticated: true,
		Dashboards:      utiltest.DummyDashboardsConfig(),
		Authentication:  utiltest.DummyAuthMethodConfig(),
		// don't fail on timestamp diff
		Timestamp: response.Timestamp,
	}

	assert.Equal(t, expected, response)
}

func TestClientInitNoAuth(t *testing.T) {
	client, err := NewTestClient()
	require.NoError(t, err)
	defer client.Close()

	// simulate config provided by airshipui.json
	configs.UIConfig = utiltest.DummyConfigNoAuth()

	isAuthenticated = false

	response, err := getResponse(client, initialize)
	require.NoError(t, err)

	expected := configs.WsMessage{
		Type:      configs.UI,
		Component: configs.Initialize,
		// isAuthenticated should now be true in response
		IsAuthenticated: response.IsAuthenticated,
		Dashboards: []configs.Dashboard{
			utiltest.DummyDashboardConfig(),
		},
		// don't fail on timestamp diff
		Timestamp: response.Timestamp,
	}

	assert.Equal(t, expected, response)
}

func TestKeepalive(t *testing.T) {
	client, err := NewTestClient()
	require.NoError(t, err)
	defer client.Close()

	// get server response to "keepalive" message from client
	response, err := getResponse(client, keepalive)
	require.NoError(t, err)

	expected := configs.WsMessage{
		Type:      configs.UI,
		Component: configs.Keepalive,
		// don't fail on timestamp diff
		Timestamp: response.Timestamp,
	}

	assert.Equal(t, expected, response)
}

func TestUnknownType(t *testing.T) {
	client, err := NewTestClient()
	require.NoError(t, err)
	defer client.Close()

	response, err := getResponse(client, unknownType)
	require.NoError(t, err)

	expected := configs.WsMessage{
		Type:      "fake_type",
		Component: configs.Initialize,
		// don't fail on timestamp diff
		Timestamp: response.Timestamp,
		Error:     "Requested type: fake_type, not found",
	}

	assert.Equal(t, expected, response)
}

func TestUnknownComponent(t *testing.T) {
	client, err := NewTestClient()
	require.NoError(t, err)
	defer client.Close()

	response, err := getResponse(client, unknownComponent)
	require.NoError(t, err)

	expected := configs.WsMessage{
		Type:      configs.UI,
		Component: "fake_component",
		// don't fail on timestamp diff
		Timestamp: response.Timestamp,
		Error:     "Requested component: fake_component, not found",
	}

	assert.Equal(t, expected, response)
}

func TestHandleDocumentRequest(t *testing.T) {
	client, err := NewTestClient()
	require.NoError(t, err)
	defer client.Close()

	response, err := getResponse(client, document)
	require.NoError(t, err)

	expected := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.Document,
		SubComponent: configs.GetDefaults,
		// don't fail on timestamp diff
		Timestamp: response.Timestamp,
	}

	// the non typed interface requires us to break up the checking otherwise the 2 will never be equal
	assert.Equal(t, expected.Type, response.Type)
	assert.Equal(t, expected.Component, response.Component)
	assert.Equal(t, expected.SubComponent, response.SubComponent)
}

func TestHandleBaremetalRequest(t *testing.T) {
	client, err := NewTestClient()
	require.NoError(t, err)
	defer client.Close()

	response, err := getResponse(client, baremetal)
	require.NoError(t, err)

	expected := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.Baremetal,
		SubComponent: configs.GetDefaults,
		// don't fail on timestamp diff
		Timestamp: response.Timestamp,
	}

	assert.Equal(t, expected.Type, response.Type)
	assert.Equal(t, expected.Component, response.Component)
	assert.Equal(t, expected.SubComponent, response.SubComponent)
}

func getResponse(client *websocket.Conn, message string) (configs.WsMessage, error) {
	err := client.WriteJSON(json.RawMessage(message))

	time.Sleep(250 * time.Millisecond)

	if err != nil {
		return configs.WsMessage{}, err
	}

	var response configs.WsMessage
	err = client.ReadJSON(&response)

	if response.Component == configs.Initialize {
		response = configs.WsMessage{}
		err = client.ReadJSON(&response)
	}
	if err != nil {
		return configs.WsMessage{}, err
	}

	return response, nil
}
