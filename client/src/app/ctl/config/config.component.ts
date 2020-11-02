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

import { Component, OnInit } from '@angular/core';
import { WebsocketService } from '../../../services/websocket/websocket.service';
import { WebsocketMessage, WSReceiver } from '../../../services/websocket/websocket.models';
import { Log } from '../../../services/log/log.service';
import { LogMessage } from '../../../services/log/log-message';
import { Context, ManagementConfig, Manifest, EncryptionConfig } from './config.models';

@Component({
  selector: 'app-bare-metal',
  templateUrl: './config.component.html',
})

export class ConfigComponent implements WSReceiver, OnInit {
  className = this.constructor.name;
  // TODO (aschiefe): extract these strings to constants
  type = 'ctl';
  component = 'config';

  currentContext: string;
  contexts: Context[] = [];
  manifests: Manifest[] = [];
  managementConfigs: ManagementConfig[] = [];
  encryptionConfigs: EncryptionConfig[] = [];

  constructor(private websocketService: WebsocketService) {
    this.websocketService.registerFunctions(this);
  }

  ngOnInit(): void {
    this.getConfig();
  }

  async receiver(message: WebsocketMessage): Promise<void> {
    if (message.hasOwnProperty('error')) {
      this.websocketService.printIfToast(message);
    } else {
      switch (message.subComponent) {
        case 'getCurrentContext':
          this.currentContext = message.message;
          break;
        case 'getContexts':
          Object.assign(this.contexts, message.data);
          break;
        case 'getManifests':
          Object.assign(this.manifests, message.data);
          break;
        case 'getEncryptionConfigs':
          Object.assign(this.encryptionConfigs, message.data);
          break;
        case 'getManagementConfigs':
          Object.assign(this.managementConfigs, message.data);
          break;
        case 'useContext':
          this.getCurrentContext();
          break;
        case 'setContext':
          this.websocketService.printIfToast(message);
          break;
        case 'setEncryptionConfig':
          this.websocketService.printIfToast(message);
          break;
        case 'setManifest':
          this.websocketService.printIfToast(message);
          break;
        case 'setManagementConfig':
          this.websocketService.printIfToast(message);
          break;
        default:
          Log.Error(new LogMessage('Config message sub component not handled', this.className, message));
          break;
      }
    }
  }

  getConfig(): void {
    this.getCurrentContext();
    this.getContexts();
    this.getManifests();
    this.getManagementConfigs();
    this.getEncryptionConfigs();
  }

  getCurrentContext(): void {
    this.websocketService.sendMessage(new WebsocketMessage(
      this.type, this.component, 'getCurrentContext')
    );
  }

  getContexts(): void {
    this.websocketService.sendMessage(new WebsocketMessage(
      this.type, this.component, 'getContexts')
    );
  }

  getManifests(): void {
    this.websocketService.sendMessage(new WebsocketMessage(
      this.type, this.component, 'getManifests')
    );
  }

  getEncryptionConfigs(): void {
    this.websocketService.sendMessage(new WebsocketMessage(
      this.type, this.component, 'getEncryptionConfigs')
    );
  }

  getManagementConfigs(): void {
    this.websocketService.sendMessage(new WebsocketMessage(
      this.type, this.component, 'getManagementConfigs')
    );
  }
}
