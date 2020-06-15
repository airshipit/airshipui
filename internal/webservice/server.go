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

	"github.com/gorilla/websocket"
	"opendev.org/airship/airshipui/internal/configs"
	"opendev.org/airship/airshipui/internal/integrations/ctl"
)

// gorilla ws specific HTTP upgrade to WebSockets
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// this is a way to allow for arbitrary messages to be processed by the backend
// the message of a specifc component is shunted to that subsystem for further processing
// TODO: make this a dynamic registration of components
var functionMap = map[configs.WsRequestType]map[configs.WsComponentType]func(configs.WsMessage) configs.WsMessage{
	configs.Electron: {
		configs.Keepalive:  keepaliveReply,
		configs.Initialize: clientInit,
	},
	configs.AirshipCTL: {
		configs.CTLConfig: ctl.HandleConfigRequest,
		configs.Baremetal: ctl.HandleBaremetalRequest,
		configs.Document:  ctl.HandleDocumentRequest,
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
		sendAlertMessage(Alerts[0])
		Alerts[0] = configs.WsMessage{}
		Alerts = Alerts[1:]
	}

	go onMessage()
}

// handle messaging to the client
func onMessage() {
	// just in case clean up the websocket
	defer onClose()

	for {
		var request configs.WsMessage
		err := ws.ReadJSON(&request)
		if err != nil {
			onError(err)
			break
		}

		// look through the function map to find the type to handle the request
		if reqType, ok := functionMap[request.Type]; ok {
			// the function map may have a component (function) to process the request
			if component, ok := reqType[request.Component]; ok {
				// get the response and tag the timestamp so it's not repeated across all functions
				response := component(request)
				response.Timestamp = time.Now().UnixNano() / 1000000
				if err = ws.WriteJSON(response); err != nil {
					onError(err)
					break
				}
			} else {
				if err = ws.WriteJSON(requestErrorHelper(fmt.Sprintf("Requested component: %s, not found",
					request.Component), request)); err != nil {
					onError(err)
					break
				}
				log.Printf("Requested component: %s, not found\n", request.Component)
			}
		} else {
			if err = ws.WriteJSON(requestErrorHelper(fmt.Sprintf("Requested type: %s, not found",
				request.Type), request)); err != nil {
				onError(err)
				break
			}
			log.Printf("Requested type: %s, not found\n", request.Type)
		}
	}
}

func requestErrorHelper(err string, request configs.WsMessage) configs.WsMessage {
	return configs.WsMessage{
		Type:      request.Type,
		Component: request.Component,
		Timestamp: time.Now().UnixNano() / 1000000,
		Error:     err,
	}
}

// The keepalive response including a timestamp from the server
// The electron / web app will occasionally ping the server due to the websocket default timeout
func keepaliveReply(configs.WsMessage) configs.WsMessage {
	return configs.WsMessage{
		Type:      configs.Electron,
		Component: configs.Keepalive,
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
	err := ws.WriteJSON(configs.WsMessage{
		Type:      configs.Electron,
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

func clientInit(configs.WsMessage) configs.WsMessage {
	// if no auth method is supplied start with minimal functionality
	if len(configs.UiConfig.AuthMethod.URL) == 0 {
		isAuthenticated = true
	}

	return configs.WsMessage{
		Type:            configs.Electron,
		Component:       configs.Initialize,
		IsAuthenticated: isAuthenticated,
		Dashboards:      configs.UiConfig.Clusters,
		Plugins:         configs.UiConfig.Plugins,
		Authentication:  configs.UiConfig.AuthMethod,
	}
}
