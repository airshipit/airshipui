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
	"os"
	"path/filepath"
	"time"

	"opendev.org/airship/airshipui/internal/configs"
)

// semaphore to signal the UI to authenticate
var isAuthenticated bool

// handle an auth complete attempt
func handleAuth(response http.ResponseWriter, request *http.Request) {
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
	webServerMux.HandleFunc("/ws", func(response http.ResponseWriter, request *http.Request) {
		onOpen(response, request)
	})

	// We can serve up static content if it's flagged as headless on command line
	// TODO: Figureout if we need to toggle the proxies on and off
	// start proxies for web based use
	startProxies()

	// static file server
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	staticContent := filepath.Join(path + string(os.PathSeparator) + "web")
	log.Println("Attempting to serve static content from ", staticContent)
	fs := http.FileServer(http.Dir(staticContent))
	webServerMux.Handle("/", fs)

	// TODO: pull ports out into conf files
	log.Println("Attempting to start webservice on localhost:8080")
	if err := http.ListenAndServe(":8080", webServerMux); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
