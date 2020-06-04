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

import "time"

// AlertLevel holds a string that determines the type of alert shown in the UI
type AlertLevel string

const (
	// Info corresponds to a blue alert message in the UI
	Info AlertLevel = "info"
	// Warning corresponds to an orange alert message in the UI
	Warning AlertLevel = "warning"
	// Error corresponds to a red alert message in the UI
	Error AlertLevel = "danger"
	// Success corresponds to a green alert message in the UI
	Success AlertLevel = "success"
)

// Alert basic structure to hold alert messages to pass to the UI
type Alert struct {
	Level   AlertLevel
	Message string
	Fade    bool
}

// Alerts serves as a queue to hold alerts to be sent to the UI,
// which will generally only be needed if any errors are encountered
// during startup before the websocket has been established
var Alerts []Alert

// SendAlert tests for the existence of an established websocket
// and either sends the message over the websocket, or adds it
// to the Alerts queue to be sent later
func SendAlert(lvl AlertLevel, msg string, fade bool) {
	alert := Alert{
		Level:   lvl,
		Message: msg,
		Fade:    fade,
	}

	if ws == nil {
		Alerts = append(Alerts, alert)
	} else {
		sendAlertMessage(alert)
	}
}

func sendAlertMessage(a Alert) {
	if err := ws.WriteJSON(map[string]interface{}{
		"type":      "electron",
		"component": "alert",
		"level":     a.Level,
		"message":   a.Message,
		"fade":      a.Fade,
		"timestamp": time.Now().UnixNano() / 1000000,
	}); err != nil {
		onError(err)
	}
}
