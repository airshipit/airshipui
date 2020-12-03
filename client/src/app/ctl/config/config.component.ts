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
import { Log } from 'src/services/log/log.service';
import { LogMessage } from 'src/services/log/log-message';
import { Context, ManagementConfig, Manifest, EncryptionConfig } from './config.models';
import { WsService } from 'src/services/ws/ws.service';
import { WsMessage, WsReceiver, WsConstants } from 'src/services/ws/ws.models';
import { MatDialog } from '@angular/material/dialog';
import { ConfigNewComponent } from './config-new/config-new.component';

@Component({
  selector: 'app-bare-metal',
  templateUrl: './config.component.html',
})

export class ConfigComponent implements WsReceiver, OnInit {
  className = this.constructor.name;
  type = WsConstants.CTL;
  component = WsConstants.CONFIG;

  airshipConfigPath: string;

  currentContext: string;
  contexts: Context[] = [];
  manifests: Manifest[] = [];
  managementConfigs: ManagementConfig[] = [];
  encryptionConfigs: EncryptionConfig[] = [];

  configs = {};

  constructor(private websocketService: WsService,
              public dialog: MatDialog) {
    this.websocketService.registerFunctions(this);
  }

  ngOnInit(): void {
    this.getConfig();
  }

  async receiver(message: WsMessage): Promise<void> {
    if (message.hasOwnProperty(WsConstants.ERROR)) {
      this.websocketService.printIfToast(message);
    } else {
      switch (message.subComponent) {
        case WsConstants.INIT:
          this.websocketService.printIfToast(message);
          this.getConfig();
          break;
        case WsConstants.SET_AIRSHIP_CONFIG:
          this.websocketService.printIfToast(message);
          this.getConfig();
          break;
        case WsConstants.GET_AIRSHIP_CONFIG_PATH:
          this.airshipConfigPath = message.message;
          break;
        case WsConstants.GET_CURRENT_CONTEXT:
          this.currentContext = message.message;
          break;
        case WsConstants.GET_CONTEXTS:
          Object.assign(this.contexts, message.data);
          break;
        case WsConstants.GET_MANIFESTS:
          this.handleGetManifests(message);
          break;
        case WsConstants.GET_ENCRYPTION_CONFIGS:
          this.handleGetEncryptionConfigs(message);
          break;
        case WsConstants.GET_MANAGEMENT_CONFIGS:
          this.handleGetManagementConfigs(message);
          break;
        case WsConstants.USE_CONTEXT:
          this.getCurrentContext();
          break;
        case WsConstants.SET_CONTEXT:
          this.websocketService.printIfToast(message);
          this.getContexts();
          break;
        case WsConstants.SET_ENCRYPTION_CONFIG:
          this.websocketService.printIfToast(message);
          this.getEncryptionConfigs();
          break;
        case WsConstants.SET_MANIFEST:
          this.websocketService.printIfToast(message);
          this.getManifests();
          break;
        case WsConstants.SET_MANAGEMENT_CONFIG:
          this.websocketService.printIfToast(message);
          this.getManagementConfigs();
          break;
        default:
          Log.Error(new LogMessage('Config message sub component not handled', this.className, message));
          break;
      }
    }
  }

  handleGetManifests(message: WsMessage): void {
    Object.assign(this.manifests, message.data);
    const manifests = 'manifests';
    this.configs[manifests] = [];
    for (const m of this.manifests) {
      this.configs[manifests].push(m.name);
    }
  }

  handleGetEncryptionConfigs(message: WsMessage): void {
    Object.assign(this.encryptionConfigs, message.data);
    const encryption = 'encryption';
    this.configs[encryption] = [];
    for (const e of this.encryptionConfigs) {
      this.configs[encryption].push(e.name);
    }
  }

  handleGetManagementConfigs(message: WsMessage): void {
    Object.assign(this.managementConfigs, message.data);
    const management = 'management';
    this.configs[management] = [];
    for (const m of this.managementConfigs) {
      this.configs[management].push(m.Name);
    }
  }

  getConfig(): void {
    this.getAirshipConfigPath();
    this.getCurrentContext();
    this.getContexts();
    this.getManifests();
    this.getManagementConfigs();
    this.getEncryptionConfigs();
  }

  getAirshipConfigPath(): void {
    this.websocketService.sendMessage(new WsMessage(
      this.type, this.component, WsConstants.GET_AIRSHIP_CONFIG_PATH)
    );
  }

  getCurrentContext(): void {
    this.websocketService.sendMessage(new WsMessage(
      this.type, this.component, WsConstants.GET_CURRENT_CONTEXT)
    );
  }

  getContexts(): void {
    this.websocketService.sendMessage(new WsMessage(
      this.type, this.component, WsConstants.GET_CONTEXTS)
    );
  }

  getManifests(): void {
    this.websocketService.sendMessage(new WsMessage(
      this.type, this.component, WsConstants.GET_MANIFESTS)
    );
  }

  getEncryptionConfigs(): void {
    this.websocketService.sendMessage(new WsMessage(
      this.type, this.component, WsConstants.GET_ENCRYPTION_CONFIGS)
    );
  }

  getManagementConfigs(): void {
    this.websocketService.sendMessage(new WsMessage(
      this.type, this.component, WsConstants.GET_MANAGEMENT_CONFIGS)
    );
  }

  newConfig(configType: string): void {
    const dialogRef = this.dialog.open(ConfigNewComponent, {
      width: '550px',
      height: '650px',
      data: { formName: configType}
    });
  }
}
