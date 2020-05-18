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
	"time"

	"opendev.org/airship/airshipui/internal/configs"
)

// Alerts serves as a queue to hold alerts to be sent to the UI,
// which will generally only be needed if any errors are encountered
// during startup before the websocket has been established
var Alerts []configs.WsMessage

// SendAlert tests for the existence of an established websocket
// and either sends the message over the websocket, or adds it
// to the Alerts queue to be sent later
func SendAlert(lvl configs.WsComponentType, msg string, fade bool) {
	alert := configs.WsMessage{
		Type:      configs.Alert,
		Component: lvl,
		Message:   msg,
		Fade:      fade,
		Timestamp: time.Now().UnixNano() / 1000000,
	}

	if ws == nil {
		Alerts = append(Alerts, alert)
	} else {
		sendAlertMessage(alert)
	}
}

func sendAlertMessage(a configs.WsMessage) {
	if err := ws.WriteJSON(a); err != nil {
		onError(err)
	}
}
