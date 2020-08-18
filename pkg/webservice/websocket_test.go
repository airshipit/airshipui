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
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/log"
)

const (
	// client messages
	keepalive        string = `{"type":"ui","component":"keepalive"}`
	unknownType      string = `{"type":"fake_type","component":"initialize"}`
	unknownComponent string = `{"type":"ui","component":"fake_component"}`
)

var client *websocket.Conn

func init() {
	u := url.URL{Scheme: "ws", Host: serverAddr, Path: "/ws"}
	var err error
	client, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)
	// get server response to "initialize" message from client which is sent by default
	var response configs.WsMessage
	err = client.ReadJSON(&response)
	if err != nil {
		log.Fatal(err)
	}
}

func TestKeepalive(t *testing.T) {
	// get server response to "keepalive" message from client
	response, err := getResponse(keepalive)
	require.NoError(t, err)

	expected := configs.WsMessage{
		SessionID: response.SessionID,
		Type:      configs.UI,
		Component: configs.Keepalive,
		// don't fail on timestamp diff
		Timestamp: response.Timestamp,
	}

	assert.Equal(t, expected, response)
}

func TestUnknownType(t *testing.T) {
	response, err := getResponse(unknownType)
	require.NoError(t, err)

	expected := configs.WsMessage{
		SessionID: response.SessionID,
		Type:      "fake_type",
		Component: configs.Initialize,
		// don't fail on timestamp diff
		Timestamp: response.Timestamp,
		Error:     "Requested type: fake_type, not found",
	}

	assert.Equal(t, expected, response)
}

func TestUnknownComponent(t *testing.T) {
	response, err := getResponse(unknownComponent)
	require.NoError(t, err)

	expected := configs.WsMessage{
		SessionID: response.SessionID,
		Type:      configs.UI,
		Component: "fake_component",
		// don't fail on timestamp diff
		Timestamp: response.Timestamp,
		Error:     "Requested component: fake_component, not found",
	}

	assert.Equal(t, expected, response)
}

func getResponse(message string) (configs.WsMessage, error) {
	err := client.WriteJSON(json.RawMessage(message))

	time.Sleep(50 * time.Millisecond)

	if err != nil {
		return configs.WsMessage{}, err
	}

	var response configs.WsMessage
	err = client.ReadJSON(&response)
	if err != nil {
		return configs.WsMessage{}, err
	}

	return response, nil
}
