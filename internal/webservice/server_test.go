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
	"net/http"
	"net/url"
	"testing"
	"time"

	"opendev.org/airship/airshipui/internal/configs"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	serverAddr string = "localhost:8080"

	// client messages
	initialize       string = `{"type":"airshipui","component":"initialize"}`
	keepalive        string = `{"type":"airshipui","component":"keepalive"}`
	unknownType      string = `{"type":"fake_type","component":"initialize"}`
	unknownComponent string = `{"type":"airshipui","component":"fake_component"}`
	document         string = `{"type":"airshipctl","component":"document","subcomponent":"getDefaults"}`
	baremetal        string = `{"type":"airshipctl","component":"baremetal","subcomponent":"getDefaults"}`
	config           string = `{"type":"airshipctl","component":"config","subcomponent":"getDefaults"}`
)

func init() {
	go WebServer()
}

func TestHandleAuth(t *testing.T) {
	client, err := NewTestClient()
	require.NoError(t, err)
	defer client.Close()

	isAuthenticated = false

	// trigger web server's handleAuth function
	_, err = http.Get("http://localhost:8080/auth")
	require.NoError(t, err)

	var response configs.WsMessage
	err = client.ReadJSON(&response)
	require.NoError(t, err)

	expected := configs.WsMessage{
		Type:      configs.AirshipUI,
		Component: configs.Authcomplete,
		// don't fail on timestamp diff
		Timestamp: response.Timestamp,
	}

	// isAuthenticated should now be true after auth complete
	assert.Equal(t, isAuthenticated, true)
	assert.Equal(t, expected, response)
}

func NewTestClient() (*websocket.Conn, error) {
	var err error
	var client *websocket.Conn
	u := url.URL{Scheme: "ws", Host: serverAddr, Path: "/ws"}
	// allow multiple attempts to establish websocket in case server isn't ready
	for i := 0; i < 5; i++ {
		client, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err == nil {
			return client, nil
		}
		time.Sleep(2 * time.Second)
	}
	return nil, err
}
