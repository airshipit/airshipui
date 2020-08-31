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

export interface WSReceiver {
    // the holy trinity of the websocket messages, a triumvirate if you will, which is how all are routed
    type: string;
    component: string;

    // This is the method which will need to be implemented in the component to handle the messages
    receiver(message: WebsocketMessage): Promise<void>;
}

// WebsocketMessage is the structure for the json that is used to talk to the backend
export class WebsocketMessage {
  sessionID: string;
  type: string;
  component: string;
  subComponent: string;
  timestamp: number;
  dashboards: Dashboard[];
  error: string;
  html: string;
  name: string;
  details: string;
  id: string;
  isAuthenticated: boolean;
  message: string;
  token: string;
  data: JSON;
  yaml: string;
  authentication: Authentication;

  // this constructor looks like this in case anyone decides they want just a raw message with no data predefined
  // or an easy way to specify the defaults
  constructor(type?: string | undefined, component?: string | undefined, subComponent?: string | undefined) {
    this.type = type;
    this.component = component;
    this.subComponent = subComponent;
  }
}

// Dashboard has the urls of the links that will pop out new dashboard tabs on the left hand side
export class Dashboard {
  name: string;
  baseURL: string;
  path: string;
  isProxied: boolean;
}

// AuthMessage is used to send and auth request and hold the token if it's authenticated
export class Authentication {
  id: string;
  password: string;

  constructor(id?: string | undefined, password?: string | undefined) {
    this.id = id;
    this.password = password;
  }
}
