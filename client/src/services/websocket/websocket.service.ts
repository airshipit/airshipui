/*
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/

import { Injectable, OnDestroy } from '@angular/core';
import { WebsocketMessage, WSReceiver, Authentication } from './websocket.models';
import { ToastrService } from 'ngx-toastr';
import 'reflect-metadata';

@Injectable({
  providedIn: 'root'
})

export class WebsocketService implements OnDestroy {
  // to avoid circular includes this has to go here
  public static token: string;
  public static tokenExpiration: number;

  private ws: WebSocket;
  private timeout: any;
  private sessionID: string;

  // functionMap is how we know where to send the direct messages
  // the structure of this map is: type -> component -> receiver
  private functionMap = new Map<string, Map<string, WSReceiver>>();

  // messageToObject unmarshalls the incoming message into a WebsocketMessage object
  private static messageToObject(incomingMessage: string): WebsocketMessage {
    const json = JSON.parse(incomingMessage);
    const obj = new WebsocketMessage();
    Object.assign(obj, json);
    return obj;
  }

  // when the WebsocketService is created the toast message is initialized and a websocket is registered
  constructor(private toastrService: ToastrService) {
    this.register();
  }

  // catch the page destroy and shut down the websocket connection normally
  ngOnDestroy(): void {
    this.ws.close();
  }

  // sendMessage will relay a WebsocketMessage to the go backend
  public async sendMessage(message: WebsocketMessage): Promise<void> {
    try {
      message.sessionID = this.sessionID;
      message.timestamp = new Date().getTime();
      if (WebsocketService.token !== undefined) { message.token = WebsocketService.token; }
      // TODO (aschiefe): determine if this debug statement is a good thing (tm)
      // Log.Debug(new LogMessage('Sending WebSocket Message', this.className, message));
      this.ws.send(JSON.stringify(message));
    } catch (err) {
      // on a refresh it may fire a request before the backend is ready so give it ye'ol retry
      // TODO (aschiefe): determine if there's a limit on retries
      return new Promise(() => setTimeout(() => { this.sendMessage(message); }, 100));
    }
  }

  // register initializes the websocket communication with the go backend
  private register(): void {
    if (this.ws !== undefined && this.ws !== null) {
      this.ws.close();
    }

    this.ws = new WebSocket('wss://localhost:10443/ws');

    this.ws.onmessage = (event) => {
      this.messageHandler(WebsocketService.messageToObject(event.data));
    };

    this.ws.onerror = (event) => {
      console.log('Web Socket received an error: ', event);
    };

    this.ws.onopen = () => {
      console.log('Websocket established');
      // start up the keepalive so the websocket-message stays open
      this.timeout = setTimeout(() => { this.keepAlive(); }, 60000);
    };

    this.ws.onclose = (event) => {
      this.close(event.code);
    };
  }

  private close(code): void {
    switch (code) {
      case 1000:
        console.log('Web Socket Closed: Normal closure: ', code);
        break;
      case 1001:
        console.log('Web Socket Closed: An endpoint is "going away", such as a server going down or a browser having navigated away from a page:', code);
        break;
      case 1002:
        console.log('Web Socket Closed: terminating the connection due to a protocol error: ', code);
        break;
      case 1003:
        console.log('Web Socket Closed: terminating the connection because it has received a type of data it cannot accept: ', code);
        break;
      case 1004:
        console.log('Web Socket Closed: Reserved. The specific meaning might be defined in the futur: ', code);
        break;
      case 1005:
        console.log('Web Socket Closed: No status code was actually present: ', code);
        break;
      case 1006:
        console.log('Web Socket Closed: The connection was closed abnormally: ', code);
        break;
      case 1007:
        console.log('Web Socket Closed: terminating the connection because it has received data within a message that was not ' +
          'consistent with the type of the message: ', code);
        break;
      case 1008:
        console.log('Web Socket Closed: terminating the connection because it has received a message that "violates its policy": ', code);
        break;
      case 1009:
        console.log('Web Socket Closed: terminating the connection because it has received a message that is too big for it to ' +
          'process: ', code);
        break;
      case 1010:
        console.log('Web Socket Closed: client is terminating the connection because it has expected the server to negotiate ' +
          'one or more extension, but the server didn\'t return them in the response message of the WebSocket handshake: ', code);
        break;
      case 1011:
        console.log('Web Socket Closed: server is terminating the connection because it encountered an unexpected condition that' +
          ' prevented it from fulfilling the request: ', code);
        break;
      case 1015:
        console.log('Web Socket Closed: closed due to a failure to perform a TLS handshake (e.g., the server certificate can\'t be' +
          ' verified): ', code);
        break;
      default:
        console.log('Web Socket Closed: unknown error code: ', code);
        break;
    }

    this.ws = null;
  }

  // Takes the WebsocketMessage and iterates through the function map to send a directed message when it shows up
  private async messageHandler(message: WebsocketMessage): Promise<void> {
    if (this.sessionID === undefined && message.hasOwnProperty('sessionID')) {
      this.sessionID = message.sessionID;
    }

    switch (message.type) {
      case 'alert': this.toastrService.warning(message.message); break; // TODO (aschiefe): improve alert handling
      default:  if (this.functionMap.hasOwnProperty(message.type)) {
                  if (this.functionMap[message.type].hasOwnProperty(message.component)){
                    this.functionMap[message.type][message.component].receiver(message);
                  } else {
                    // special case where we want to handle all top level messages at a specific component
                    if (this.functionMap[message.type].hasOwnProperty('any')) {
                      this.functionMap[message.type].any.receiver(message);
                    } else {
                      this.printIfToast(message);
                    }
                  }
                } else {
                  this.toastrService.info(message.message);
                }
                break;
    }
  }

  // websockets time out after 5 minutes of inactivity, this keeps the backend engaged so it doesn't time
  private keepAlive(): void {
    // clear the previously set timeout
    window.clearTimeout(this.timeout);
    window.clearInterval(this.timeout);
    if (this.ws !== undefined && this.ws !== null && this.ws.readyState !== this.ws.CLOSED) {
      this.sendMessage(new WebsocketMessage('ui', 'keepalive', null));
    }
    this.timeout = setTimeout(() => { this.keepAlive(); }, 60000);
  }

  // registerFunctions is a is called out of the target's constructor so it can auto populate the function map
  public registerFunctions(target: WSReceiver): void {
    const type = target.type;
    const component = target.component;
    if (this.functionMap.hasOwnProperty(type)) {
      this.functionMap[type][component] = target;
    } else {
      const components = new Map<string, WSReceiver>();
      components[component] = target;
      this.functionMap[type] = components;
    }
  }

  // printIfToast puts up the toast popup message on the UI
  printIfToast(message: WebsocketMessage): void {
    if (message.error !== undefined && message.error !== null) {
      this.toastrService.error(message.error);
    } else {
      console.log(message);
      this.toastrService.info(message.message);
    }
  }
}
