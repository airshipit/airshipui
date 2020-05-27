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
	"fmt"
	"log"
	"net/http"
	"time"

	"opendev.org/airship/airshipui/internal/configs"

	"github.com/gorilla/websocket"
)

// just a base structure to return from the web service
type wsRequest struct {
	Type      string                 `json:"type,omitempty"`
	Component string                 `json:"component,omitempty"`
	Error     string                 `json:"error"`
	Data      map[string]interface{} `json:"data"`
}

// Alert basic structure to hold alert messages to pass to the UI
type Alert struct {
	Level   string
	Message string
}

var Alerts []Alert

// gorilla ws specific HTTP upgrade to WebSockets
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// this is a way to allow for arbitrary messages to be processed by the backend
// most likely we will need to have sub components register with the system
// TODO: make this a dynamic registration of components
var functionMap = map[string]map[string]func() map[string]interface{}{
	"electron": {
		"keepalive":  keepaliveReply,
		"initialize": clientInit,
	},
}

// websocket that'll be reused by several places
var ws *websocket.Conn

// semaphore to signal the UI to authenticate
var isAuthenticated bool

// handle the origin request & upgrade to websocket
func onOpen(w http.ResponseWriter, r *http.Request) {
	// gorilla ws will give a 403 on a cross origin request, so we silence its complaints
	// This happens with electron because it's sending an origin of 'file://' instead of 'localhost:8080'
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade to websocket protocol over http
	log.Printf("Establishing the websocket")
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not open websocket connection from: %s\n", r.Host)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	ws = wsConn
	log.Printf("WebSocket established with %s\n", ws.RemoteAddr().String())

	// send any initialization alerts to UI and clear the queue
	for len(Alerts) > 0 {
		sendAlert(Alerts[0].Level, Alerts[0].Message)
		Alerts[0] = Alert{}
		Alerts = Alerts[1:]
	}

	go onMessage()
}

// handle messaging to the client
func onMessage() {
	// just in case clean up the websocket
	defer onClose()

	for {
		var request wsRequest
		err := ws.ReadJSON(&request)
		if err != nil {
			onError(err)
			break
		}

		// look through the function map to find the type to handle the request
		if reqType, ok := functionMap[request.Type]; ok {
			// the function map may have a component (function) to process the request
			if component, ok := reqType[request.Component]; ok {
				if err = ws.WriteJSON(component()); err != nil {
					onError(err)
					break
				}
			} else {
				request.Error = fmt.Sprintf("Requested component: %s, not found", request.Component)
				if err = ws.WriteJSON(request); err != nil {
					onError(err)
					break
				}
				log.Printf("Requested component: %s, not found\n", request.Component)
			}
		} else {
			request.Error = fmt.Sprintf("Requested type: %s, not found", request.Type)
			if err = ws.WriteJSON(request); err != nil {
				onError(err)
				break
			}
			log.Printf("Requested type: %s, not found\n", request.Type)
		}
	}
}

// The keepalive response including a timestamp from the server
// The electron / web app will occasionally ping the server due to the websocket default timeout
func keepaliveReply() map[string]interface{} {
	return map[string]interface{}{
		"type":      "electron",
		"component": "keepalive",
		"timestamp": time.Now().UnixNano() / 1000000,
	}
}

// common websocket close with logging
func onClose() {
	log.Printf("Closing websocket")
	// ws.Close()
}

// common websocket error handling with logging
func onError(err error) {
	log.Printf("Error receiving / sending message: %s\n", err)
}

// handle an auth complete attempt
func handleAuth(w http.ResponseWriter, r *http.Request) {
	// TODO: handle the response body to capture the credentials
	err := ws.WriteJSON(map[string]interface{}{
		"type":      "electron",
		"component": "authcomplete",
		"timestamp": time.Now().UnixNano() / 1000000,
	})

	// error sending the websocket request
	if err != nil {
		onError(err)
	} else {
		isAuthenticated = true
	}
}

// SendAlert sends an alert message to the frontend handler
// to display alerts in the UI itself
func sendAlert(alertLvl string, msg string) {
	if err := ws.WriteJSON(map[string]interface{}{
		"type":      "electron",
		"component": "alert",
		"level":     alertLvl,
		"message":   msg,
		"timestamp": time.Now().UnixNano() / 1000000,
	}); err != nil {
		onError(err)
	}
}

// WebServer will run the handler functions for WebSockets
// TODO: potentially add in the ability to serve static content
func WebServer() {
	// some things may need a redirect so we'll give them a url to do that with
	http.HandleFunc("/auth", handleAuth)

	// hand off the websocket upgrade over http
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		onOpen(w, r)
	})

	log.Println("Attempting to start webservice on localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func clientInit() map[string]interface{} {
	// if no auth method is supplied start with minimal functionality
	if len(configs.UiConfig.AuthMethod.URL) == 0 {
		isAuthenticated = true
	}

	return map[string]interface{}{
		"type":            "electron",
		"component":       "initialize",
		"timestamp":       time.Now().UnixNano() / 1000000,
		"isAuthenticated": isAuthenticated,
		"dashboards":      configs.UiConfig.Clusters,
		"plugins":         configs.UiConfig.Plugins,
		"authentication":  configs.UiConfig.AuthMethod,
	}
}
