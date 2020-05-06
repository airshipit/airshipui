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
var ws = null;
var timeout = null;

// establish a session when browser is open
if (document.addEventListener) {
    document.addEventListener("DOMContentLoaded", function () {
        // register the webservice so it's available the entire time of the process
        register();
    }, false);
}

// listen for the unload event and close the web socket if open
document.addEventListener("unload", function () {
    if (ws !== null) {
        if (ws.readyState !== ws.CLOSED) { ws.close(); }
    }
});

function register() {
    if (ws !== null) {
        ws.close();
        ws = null;
    }

    ws = new WebSocket("ws://localhost:8080/ws");

    ws.onmessage = function (event) {
        handleMessages(event);
    }

    ws.onerror = function (event) {
        console.log("Web Socket received an error: ", event.code);
        close(event.code);
    }

    ws.onopen = function () {
        open();
    }

    ws.onclose = function (event) {
        close(event.code);
    }
}

function handleMessages(message) {
    var json = JSON.parse(message.data);
    // keepalives and inits aren't interesting to other pages
    if (json["type"] === "electron") {
        console.log(json);
        if (json["component"] === "initialize") {
            if (!json["isAuthenticated"]) {
                authenticate(json["authentication"]);
            } else {
                authComplete();
            }
            addPlugins(json["plugins"]);
        } else if (json["component"] === "authcomplete") {
            authComplete();
        }
    } else {
        // TODO: determine if we're dispatching events or just doing function calls
        // events based on the type are interesting to other pages
        // create and dispatch an event based on the data received
        document.dispatchEvent(new CustomEvent(json["type"], { detail: json }));
    }

    // TODO: Determine if these should be suppressed or only allowed in specific cases
    console.log("Received message: " + message.data);
}

function open() {
    console.log("Websocket established");
    var json = { "type": "electron", "component": "initialize" };
    ws.send(JSON.stringify(json));
    // start up the keepalive so the websocket stays open
    keepAlive();
}

function close(code) {
    switch (code) {
        case 1000: console.log("Web Socket Closed: Normal closure: ", code); break;
        case 1001: console.log("Web Socket Closed: An endpoint is \"going away\", such as a server going down or a browser having navigated away from a page:", code); break;
        case 1002: console.log("Web Socket Closed: terminating the connection due to a protocol error: ", code); break;
        case 1003: console.log("Web Socket Closed: terminating the connection because it has received a type of data it cannot accept: ", code); break;
        case 1004: console.log("Web Socket Closed: Reserved. The specific meaning might be defined in the futur: ", code); break;
        case 1005: console.log("Web Socket Closed: No status code was actually present: ", code); break;
        case 1006: console.log("Web Socket Closed: The connection was closed abnormally: ", code); break;
        case 1007: console.log("Web Socket Closed: terminating the connection because it has received data within a message that was not consistent with the type of the message: ", code); break;
        case 1008: console.log("Web Socket Closed: terminating the connection because it has received a message that \"violates its policy\": ", code); break;
        case 1009: console.log("Web Socket Closed: terminating the connection because it has received a message that is too big for it to process: ", code); break;
        case 1010: console.log("Web Socket Closed: client is terminating the connection because it has expected the server to negotiate one or more extension, but the server didn't return them in the response message of the WebSocket handshake: ", code); break;
        case 1011: console.log("Web Socket Closed: server is terminating the connection because it encountered an unexpected condition that prevented it from fulfilling the request: ", code); break;
        case 1015: console.log("Web Socket Closed: closed due to a failure to perform a TLS handshake (e.g., the server certificate can't be verified): ", code); break;
        default: console.log("Web Socket Closed: unknown error code: ", code); break;
    }

    ws = null;
}

function authComplete() {
    document.getElementById("HeaderDiv").style.display = "";
    document.getElementById("MainDiv").style.display = "";
    document.getElementById("DashView").style.display = "none";
    document.getElementById("FooterDiv").style.display = "";
}

function keepAlive() {
    if (ws !== null) {
        if (ws.readyState !== ws.CLOSED) {
            // clear the previously set timeout
            window.clearTimeout(timeout);
            window.clearInterval(timeout);
            var json = { "id": "poc", "type": "electron", "component": "keepalive" };
            ws.send(JSON.stringify(json));
            timeout = window.setTimeout(keepAlive, 60000);
        }
    }
}

function sendMessage(json) { // eslint-disable-line no-unused-vars
    if (ws.readyState === WebSocket.CLOSED) {
        register();
    }
    console.log("Attempting to send: ", json);
    ws.send(JSON.stringify(json));
}