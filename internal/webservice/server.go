/*
 Copyright (c) 2020 AT&T. All Rights Reserved.

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

	"github.com/gorilla/websocket"
	"opendev.org/airship/airshipui/internal/plugin/airshipopenstack"
)

// just a base structure to return from the web service
type Message struct {
	ID      int    `json:"id,omitempty"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type wsRequest struct {
	ID        string              `json:"id,omitempty"`
	Type      string              `json:"type,omitempty"`
	Component string              `json:"component,omitempty"`
	Error     string              `json:"error"`
	Data      []map[string]string `json:"data"`
}

// gorilla ws specific HTTP upgrade to WebSockets
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// this is a way to allow for arbitrary messages to be processed by the backend
// most likely we will need to have sub components register with the system
// TODO: make this a dynamic registration of components
var functionMap = map[string]map[string]func() []map[string]string{
	"openstack": {
		"getFlavors":  airshipopenstack.GetFlavors,
		"getImages":   airshipopenstack.GetImages,
		"getVMs":      airshipopenstack.GetVMs,
		"getDomains":  airshipopenstack.GetDomains,
		"getProjects": airshipopenstack.GetProjects,
		"getUsers":    airshipopenstack.GetUsers,
		"getNetworks": airshipopenstack.GetNetworks,
		"getSubnets":  airshipopenstack.GetSubnets,
		"getVolumes":  airshipopenstack.GetVolumes,
	},
	"electron": {
		"keepalive": basicReply,
		"getID":     basicReply,
	},
}

var ws *websocket.Conn

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
				request.Data = component()
				if err = ws.WriteJSON(request); err != nil {
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

func basicReply() []map[string]string {
	m := make([]map[string]string, 0)
	m = append(m, map[string]string{
		"ID":        "foo",
		"type":      "bar",
		"component": "glitch",
	})
	return m
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

// WebServer will run the handler functions for both normal REST requests and WebSockets
func WebServer() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		onOpen(w, r)
	})

	log.Println("Attempting to start webservice on localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
