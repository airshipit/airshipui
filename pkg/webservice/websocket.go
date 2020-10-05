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
	"opendev.org/airship/airshipui/pkg/log"
	"opendev.org/airship/airshipui/pkg/statistics"
)

// Session is a struct to hold information about a given session
type session struct {
	sessionID  string
	jwt        string
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
var funcMap = map[configs.WsRequestType]map[configs.WsComponentType]func(*string, configs.WsMessage) configs.WsMessage{
	configs.UI: {
		configs.Keepalive: keepaliveReply,
		configs.Auth:      handleAuth,
	},
}

// AppendToFunctionMap allows us to break up the circular reference from the other packages
// It does however require them to implement an init function to append them
// TODO: maybe some form of an interface to enforce this may be necessary?
func AppendToFunctionMap(requestType configs.WsRequestType,
	functions map[configs.WsComponentType]func(*string, configs.WsMessage) configs.WsMessage) {
	funcMap[requestType] = functions
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
	log.Debugf("WebSocket session %s established with %s\n", session.sessionID, session.ws.RemoteAddr().String())

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
			// test the auth token for request validity on non auth requests
			// TODO (aschiefe): this will need to be amended when refresh tokens are implemented
			var user *string
			if request.Type != configs.UI && request.Component != configs.Auth && request.SubComponent != configs.Authenticate {
				if request.Token != nil {
					user, err = validateToken(*request.Token)
				} else {
					err = errors.New("No authentication token found")
				}
			}
			if err != nil {
				// deny the request if we get a bad token, this will force the UI to a login screen
				e := "Invalid token, authentication denied"
				response := configs.WsMessage{
					Type:         configs.UI,
					Component:    configs.Auth,
					SubComponent: configs.Denied,
					Error:        &e,
				}
				if err = session.webSocketSend(response); err != nil {
					session.onError(err)
				}
			} else {
				// This is the middleware to be able to record when a transaction starts and ends for the statistics recorder
				// It is possible for the backend to send messages without a valid user
				transaction := statistics.NewTransaction(user, request)

				// look through the function map to find the type to handle the request
				if reqType, ok := funcMap[request.Type]; ok {
					// the function map may have a component (function) to process the request
					if component, ok := reqType[request.Component]; ok {
						response := component(user, request)
						if err = session.webSocketSend(response); err != nil {
							session.onError(err)
						}
						go transaction.Complete(response.Error == nil)
					} else {
						if err = session.webSocketSend(requestErrorHelper(fmt.Sprintf("Requested component: %s, not found",
							request.Component), request)); err != nil {
							session.onError(err)
						}
						log.Errorf("Requested component: %s, not found\n", request.Component)
						go transaction.Complete(false)
					}
				} else {
					if err = session.webSocketSend(requestErrorHelper(fmt.Sprintf("Requested type: %s, not found",
						request.Type), request)); err != nil {
						session.onError(err)
					}
					log.Errorf("Requested type: %s, not found\n", request.Type)
					go transaction.Complete(false)
				}
			}
		}()
	}
}

// common websocket close with logging
func (session *session) onClose() {
	log.Debugf("Closing websocket for session %s", session.sessionID)
	session.ws.Close()
	delete(sessions, session.sessionID)
}

// common websocket error handling with logging
func (session *session) onError(err error) {
	log.Errorf("Error receiving / sending message: %s\n", err)
}

// The UI will occasionally ping the server due to the websocket default timeout
func keepaliveReply(*string, configs.WsMessage) configs.WsMessage {
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
		Error:     &err,
	}
}

// newSession generates a new session
func newSession(ws *websocket.Conn) *session {
	id := uuid.New().String()

	session := &session{
		sessionID: id,
		ws:        ws,
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
	response.SessionID = session.sessionID

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
		Type:       configs.UI,
		Component:  configs.Initialize,
		Dashboards: configs.UIConfig.Dashboards,
		AuthMethod: configs.UIConfig.AuthMethod,
	}); err != nil {
		log.Errorf("Error receiving / sending init to session %s: %s\n", session.sessionID, err)
	}
}

// CloseAllSessions is called when the system is exiting to cleanly close all the current connections
func CloseAllSessions() {
	for _, session := range sessions {
		session.onClose()
	}
}
