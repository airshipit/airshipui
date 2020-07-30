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
	"log"
	"net/http"
	"time"

	"opendev.org/airship/airshipui/internal/configs"
)

// semaphore to signal the UI to authenticate
var isAuthenticated bool

// handle an auth complete attempt
func handleAuth(http.ResponseWriter, *http.Request) {
	// TODO: handle the response body to capture the credentials
	err := ws.WriteJSON(configs.WsMessage{
		Type:      configs.AirshipUI,
		Component: configs.Authcomplete,
		Timestamp: time.Now().UnixNano() / 1000000,
	})

	// error sending the websocket request
	if err != nil {
		onError(err)
	} else {
		isAuthenticated = true
	}
}

// WebServer will run the handler functions for WebSockets
func WebServer() {
	webServerMux := http.NewServeMux()

	// some things may need a redirect so we'll give them a url to do that with
	webServerMux.HandleFunc("/auth", handleAuth)

	// hand off the websocket upgrade over http
	webServerMux.HandleFunc("/ws", onOpen)

	// establish routing to static angular client
	webServerMux.HandleFunc("/", serveFile)

	// TODO: Figureout if we need to toggle the proxies on and off
	// start proxies for web based use
	startProxies()

	// TODO: pull ports out into conf files
	log.Println("Attempting to start webservice on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", webServerMux))
}
