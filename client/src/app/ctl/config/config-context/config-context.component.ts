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

import { Component, OnInit, Input } from '@angular/core';
import { Context, ContextOptions } from '../config.models';
import { WsService } from 'src/services/ws/ws.service';
import { FormControl } from '@angular/forms';
import { WsMessage, WsConstants } from 'src/services/ws/ws.models';

@Component({
  selector: 'app-config-context',
  templateUrl: './config-context.component.html',
  styleUrls: ['./config-context.component.css']
})
export class ConfigContextComponent implements OnInit {
  @Input() context: Context;
  type = WsConstants.CTL;
  component = WsConstants.CONFIG;

  locked = true;

  name = new FormControl({value: '', disabled: true});
  contextKubeconf = new FormControl({value: '', disabled: true});
  manifest = new FormControl({value: '', disabled: true});
  managementConfiguration = new FormControl({value: '', disabled: true});
  encryptionConfig = new FormControl({value: '', disabled: true});

  controlsArray = [this.name, this.contextKubeconf, this.manifest, this.managementConfiguration, this.encryptionConfig];

  constructor(private websocketService: WsService) {}

  ngOnInit(): void {
    this.name.setValue(this.context.name);
    this.contextKubeconf.setValue(this.context.contextKubeconf);
    this.manifest.setValue(this.context.manifest);
    this.encryptionConfig.setValue(this.context.encryptionConfig);
    this.managementConfiguration.setValue(this.context.managementConfiguration);
  }

  toggleLock(): void {
    for (const control of this.controlsArray) {
      if (this.locked) {
        control.enable();
      } else {
        control.disable();
      }
    }

    this.locked = !this.locked;
  }

  setContext(): void {
    const opts: ContextOptions = {
      Name: this.name.value,
      Manifest: this.manifest.value,
      ManagementConfiguration: this.managementConfiguration.value,
      EncryptionConfig: this.encryptionConfig.value,
    };

    const msg = new WsMessage(this.type, this.component, WsConstants.SET_CONTEXT);
    msg.data = JSON.parse(JSON.stringify(opts));
    msg.name = this.name.value;

    this.websocketService.sendMessage(msg);
    this.toggleLock();
  }

  useContext(name: string): void {
    const msg = new WsMessage(this.type, this.component, WsConstants.USE_CONTEXT);
    msg.name = name;
    this.websocketService.sendMessage(msg);
  }
}
