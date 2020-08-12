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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"opendev.org/airship/airshipui/pkg/configs"
)

func TestSendAlert(t *testing.T) {
	client, err := NewTestClient()
	require.NoError(t, err)
	defer client.Close()

	// construct and send alert from server to client
	SendAlert(configs.Error, "Test Alert")

	response, err := MessageReader(client)
	require.NoError(t, err)

	expected := configs.WsMessage{
		Type:      configs.Alert,
		Component: configs.Error,
		Message:   "Test Alert",
		// don't fail on timestamp diff
		Timestamp: response.Timestamp,
	}

	assert.Equal(t, expected, response)
}

func TestSendAlertNoWebSocket(t *testing.T) {
	// test requires that ws == nil
	conn := ws
	ws = nil
	defer func() {
		ws = conn
		Alerts = nil
	}()

	// queue should be empty
	Alerts = nil

	SendAlert(configs.Info, "Test Alert")

	// ws is nil, so the queue should now have 1 Alert
	assert.Len(t, Alerts, 1)

	expected := configs.WsMessage{
		Type:      configs.Alert,
		Component: configs.Info,
		Message:   "Test Alert",
		// don't fail on timestamp diff
		Timestamp: Alerts[0].Timestamp,
	}

	assert.Equal(t, expected, Alerts[0])
}
