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
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/ctl"
)

// gorilla ws specific HTTP upgrade to WebSockets
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// websocket that'll be reused by several places
var ws *websocket.Conn
var writeMutex sync.Mutex

// this is a way to allow for arbitrary messages to be processed by the backend
// the message of a specifc component is shunted to that subsystem for further processing
var functionMap = map[configs.WsRequestType]map[configs.WsComponentType]func(configs.WsMessage) configs.WsMessage{
	configs.UI: {
		configs.Keepalive:  keepaliveReply,
		configs.Initialize: clientInit,
	},
	configs.CTL: ctl.CTLFunctionMap,
}

// handle the origin request & upgrade to websocket
func onOpen(response http.ResponseWriter, request *http.Request) {
	// gorilla ws will give a 403 on a cross origin request, so to silence its complaints
	// upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade to websocket protocol over http
	log.Printf("Establishing the websocket")
	wsConn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Printf("Could not open websocket connection from: %s\n", request.Host)
		http.Error(response, "Could not open websocket connection", http.StatusBadRequest)
	}

	ws = wsConn
	log.Printf("WebSocket established with %s\n", ws.RemoteAddr().String())

	go onMessage()
	sendInit()
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

		// this has to be a go routine otherwise it will block any incoming messages waiting for a command return
		go func() {
			// look through the function map to find the type to handle the request
			if reqType, ok := functionMap[request.Type]; ok {
				// the function map may have a component (function) to process the request
				if component, ok := reqType[request.Component]; ok {
					// get the response and tag the timestamp so it's not repeated across all functions

					response := component(request)
					response.Timestamp = time.Now().UnixNano() / 1000000
					if err = WebSocketSend(response); err != nil {
						onError(err)
					}
				} else {
					if err = WebSocketSend(requestErrorHelper(fmt.Sprintf("Requested component: %s, not found",
						request.Component), request)); err != nil {
						onError(err)
					}
					log.Printf("Requested component: %s, not found\n", request.Component)
				}
			} else {
				if err = WebSocketSend(requestErrorHelper(fmt.Sprintf("Requested type: %s, not found",
					request.Type), request)); err != nil {
					onError(err)
				}
				log.Printf("Requested type: %s, not found\n", request.Type)
			}
		}()
	}
}

// common websocket close with logging
func onClose() {
	log.Printf("Closing websocket")
}

// common websocket error handling with logging
func onError(err error) {
	log.Printf("Error receiving / sending message: %s\n", err)
}

// WebSocketSend allows for the sender to be thread safe, we cannot write to the websocket at the same time
func WebSocketSend(response configs.WsMessage) error {
	writeMutex.Lock()
	defer writeMutex.Unlock()

	return ws.WriteJSON(response)
}

// The keepalive response including a timestamp from the server
// The UI will occasionally ping the server due to the websocket default timeout
func keepaliveReply(configs.WsMessage) configs.WsMessage {
	return configs.WsMessage{
		Type:      configs.UI,
		Component: configs.Keepalive,
	}
}

// formats an error response in the way that we're expecting on the UI
func requestErrorHelper(err string, request configs.WsMessage) configs.WsMessage {
	return configs.WsMessage{
		Type:      request.Type,
		Component: request.Component,
		Timestamp: time.Now().UnixNano() / 1000000,
		Error:     err,
	}
}

// sendInit is generated on the onOpen event and sends the information the UI needs to startup
func sendInit() {
	response := clientInit(configs.WsMessage{
		Timestamp: time.Now().UnixNano() / 1000000,
	})

	if err := WebSocketSend(response); err != nil {
		onError(err)
	}
}

// clientInit is in the function map if the client requests an init message this is the handler
// TODO (asciefe): determine if this is still necessary
func clientInit(configs.WsMessage) configs.WsMessage {
	// if no auth method is supplied start with minimal functionality
	if configs.UIConfig.AuthMethod == nil {
		isAuthenticated = true
	}

	return configs.WsMessage{
		Type:            configs.UI,
		Component:       configs.Initialize,
		IsAuthenticated: isAuthenticated,
		Dashboards:      configs.UIConfig.Dashboards,
		Authentication:  configs.UIConfig.AuthMethod,
	}
}
