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
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/ctl"
	"opendev.org/airship/airshipui/pkg/log"
)

// session is a struct to hold information about a given session
type session struct {
	id         string
	writeMutex sync.Mutex
	ws         *websocket.Conn
}

// sessions keeps track of open websocket sessions
var sessions = map[string]*session{}

// gorilla ws specific HTTP upgrade to WebSockets
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// this is a way to allow for arbitrary messages to be processed by the backend
// the message of a specifc component is shunted to that subsystem for further processing
var functionMap = map[configs.WsRequestType]map[configs.WsComponentType]func(configs.WsMessage) configs.WsMessage{
	configs.UI: {
		configs.Keepalive: keepaliveReply,
	},
	configs.CTL: ctl.CTLFunctionMap,
}

// handle the origin request & upgrade to websocket
func onOpen(response http.ResponseWriter, request *http.Request) {
	// gorilla ws will give a 403 on a cross origin request, so to silence its complaints
	// upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade to websocket protocol over http
	wsConn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Errorf("Could not open websocket connection from: %s\n", request.Host)
		http.Error(response, "Could not open websocket connection", http.StatusBadRequest)
	}

	session := newSession(wsConn)
	log.Debugf("WebSocket session %s established with %s\n", session.id, session.ws.RemoteAddr().String())

	go session.onMessage()
}

// handle messaging to the client
func (session *session) onMessage() {
	// just in case clean up the websocket
	defer session.onClose()

	for {
		var request configs.WsMessage
		err := session.ws.ReadJSON(&request)
		if err != nil {
			session.onError(err)
			break
		}

		// this has to be a go routine otherwise it will block any incoming messages waiting for a command return
		go func() {
			// look through the function map to find the type to handle the request
			if reqType, ok := functionMap[request.Type]; ok {
				// the function map may have a component (function) to process the request
				if component, ok := reqType[request.Component]; ok {
					response := component(request)
					if err = session.webSocketSend(response); err != nil {
						session.onError(err)
					}
				} else {
					if err = session.webSocketSend(requestErrorHelper(fmt.Sprintf("Requested component: %s, not found",
						request.Component), request)); err != nil {
						session.onError(err)
					}
					log.Errorf("Requested component: %s, not found\n", request.Component)
				}
			} else {
				if err = session.webSocketSend(requestErrorHelper(fmt.Sprintf("Requested type: %s, not found",
					request.Type), request)); err != nil {
					session.onError(err)
				}
				log.Errorf("Requested type: %s, not found\n", request.Type)
			}
		}()
	}
}

// common websocket close with logging
func (session *session) onClose() {
	log.Debugf("Closing websocket for session %s", session.id)
	session.ws.Close()
	delete(sessions, session.id)
}

// common websocket error handling with logging
func (session *session) onError(err error) {
	log.Errorf("Error receiving / sending message: %s\n", err)
}

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
		Error:     err,
	}
}

// newSession generates a new session
func newSession(ws *websocket.Conn) *session {
	id := uuid.New().String()

	session := &session{
		id: id,
		ws: ws,
	}

	// keep track of the session
	sessions[id] = session

	// send the init message to the client
	go session.sendInit()

	return session
}

// webSocketSend allows for the sender to be thread safe, we cannot write to the websocket at the same time
func (session *session) webSocketSend(response configs.WsMessage) error {
	session.writeMutex.Lock()
	defer session.writeMutex.Unlock()
	response.Timestamp = time.Now().UnixNano() / 1000000
	response.SessionID = session.id

	return session.ws.WriteJSON(response)
}

// WebSocketSend allows of other packages to send a request for the websocket
func WebSocketSend(response configs.WsMessage) error {
	if session, ok := sessions[response.SessionID]; ok {
		return session.webSocketSend(response)
	}

	return errors.New("session id " + response.SessionID + "not found")
}

// sendInit is generated on the onOpen event and sends the information the UI needs to startup
func (session *session) sendInit() {
	if err := session.webSocketSend(configs.WsMessage{
		Type:            configs.UI,
		Component:       configs.Initialize,
		IsAuthenticated: true,
		Dashboards:      configs.UIConfig.Dashboards,
		Authentication:  configs.UIConfig.AuthMethod,
	}); err != nil {
		log.Errorf("Error receiving / sending init to session %s: %s\n", session.id, err)
	}
}

// CloseAllSessions is called when the system is exiting to cleanly close all the current connections
func CloseAllSessions() {
	for _, session := range sessions {
		session.onClose()
	}
}
