import { Injectable } from '@angular/core';
import { WebsocketMessage } from './models/websocket-message/websocket-message';
import { Subject } from 'rxjs';
import {Dashboard} from './models/websocket-message/dashboard/dashboard';
import {Executable} from './models/websocket-message/dashboard/executable/executable';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {

  public subject = new Subject<WebsocketMessage>();
  private ws: WebSocket;
  private timeout: number;

  private static convertIncomingMessageJsonToObject(incomingMessage: string): WebsocketMessage {
    const json = JSON.parse(incomingMessage);
    const messageTransform = new WebsocketMessage();
    if (typeof json.type === 'string') {
      messageTransform.type = json.type;
    }
    if (typeof json.component === 'string') {
      messageTransform.component = json.component;
    }
    if (typeof json.subComponent === 'string') {
      messageTransform.subComponent = json.subComponent;
    }
    if (typeof json.timestamp === 'number') {
      messageTransform.timestamp = json.timestamp;
    }
    if (typeof json.error === 'string') {
      messageTransform.error = json.error;
    }
    if (typeof json.fade === 'boolean') {
      messageTransform.fade = json.fade;
    }
    if (typeof json.html === 'string') {
      messageTransform.html = json.html;
    }
    if (typeof json.isAuthenticated === 'boolean') {
      messageTransform.isAuthenticated = json.isAuthenticated;
    }
    if (typeof json.message === 'string') {
      messageTransform.message = json.message;
    }
    if (typeof json.data === 'string' && JSON.parse(json.data)) {
      messageTransform.data = JSON.parse(json.data);
    }
    if (typeof json.yaml === 'string') {
      messageTransform.yaml = json.yaml;
    }
    if (typeof json.dashboards !== undefined && Array.isArray(json.dashboards)) {
      json.dashboards.forEach(dashboard => {
        const dashboardTransform = new Dashboard();
        if (typeof dashboard.name === 'string') {
          dashboardTransform.name = dashboard.name;
        }
        if (typeof dashboard.baseURL === 'string') {
          dashboardTransform.baseURL = dashboard.baseURL;
        }
        if (typeof dashboard.path === 'string') {
          dashboardTransform.path = dashboard.path;
        }
        if (typeof dashboard.isProxied === 'boolean') {
          dashboardTransform.isProxied = dashboard.isProxied;
        }
        if (typeof dashboard.executable === 'object') {
          const executableTransform = new Executable();
          if (typeof dashboard.executable.autoStart === 'boolean') {
            executableTransform.autoStart = dashboard.executable.autoStart;
          }
          if (typeof dashboard.executable.filePath === 'string') {
            executableTransform.filePath = dashboard.executable.filePath;
          }
          if (typeof dashboard.executable.args !== undefined && Array.isArray(typeof dashboard.executable.args)) {
            dashboard.executable.args.forEach(arg => {
              if (typeof arg === 'string') {
                executableTransform.args.push(arg);
              }
            });
          }
        }
        if (messageTransform.dashboards === undefined) {
          messageTransform.dashboards = [];
        }
        messageTransform.dashboards.push(dashboardTransform);
      });
    }
    return messageTransform;
  }

  constructor() {
    this.register();
  }

  public sendMessage(message: WebsocketMessage): void {
    message.timestamp = new Date().getTime();
    this.ws.send(JSON.stringify(message));
  }

  private register(): void {
    if (this.ws !== undefined && this.ws !== null) {
      this.ws.close();
      this.ws = null;
    }

    this.ws = new WebSocket('ws://localhost:8080/ws');

    this.ws.onmessage = (event) => {
      this.subject.next(WebsocketService.convertIncomingMessageJsonToObject(event.data));
    };

    this.ws.onerror = (event) => {
      console.log('Web Socket received an error: ', event);
    };

    this.ws.onopen = () => {
      console.log('Websocket established');
      const json = { type: 'airshipui', component: 'initialize' };
      this.ws.send(JSON.stringify(json));
      // start up the keepalive so the websocket-message stays open
      this.keepAlive();
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

  private keepAlive(): void {
    if (this.ws !== undefined && this.ws !== null && this.ws.readyState !== this.ws.CLOSED) {
      // clear the previously set timeout
      window.clearTimeout(this.timeout);
      window.clearInterval(this.timeout);
      const json = { type: 'airshipui', component: 'keepalive' };
      this.ws.send(JSON.stringify(json));
      this.timeout = window.setTimeout(this.keepAlive, 60000);
    }
  }
}
