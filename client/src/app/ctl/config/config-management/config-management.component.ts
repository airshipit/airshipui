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

import { Component, Input, OnInit } from '@angular/core';
import { ManagementConfig } from '../config.models';
import { FormControl, Validators } from '@angular/forms';
import { WsService } from 'src/services/ws/ws.service';
import { WsMessage, WsConstants } from 'src/services/ws/ws.models';

@Component({
  selector: 'app-config-management',
  templateUrl: './config-management.component.html',
  styleUrls: ['./config-management.component.css']
})
export class ConfigManagementComponent implements OnInit {
  @Input() config: ManagementConfig;
  msgType = WsConstants.CTL;
  component = WsConstants.CONFIG;

  locked = true;

  name = new FormControl({value: '', disabled: true});
  insecure = new FormControl({value: false, disabled: true});
  systemActionRetries = new FormControl({value: '', disabled: true}, Validators.pattern('[0-9]*'));
  systemRebootDelay = new FormControl({value: '', disabled: true}, Validators.pattern('[0-9]*'));
  type = new FormControl({value: '', disabled: true});
  useproxy = new FormControl({value: false, disabled: true});

  controlsArray = [this.name, this.insecure, this.systemRebootDelay, this.systemActionRetries, this.type, this.useproxy];

  constructor(private websocketService: WsService) { }

  ngOnInit(): void {
    this.name.setValue(this.config.name);
    this.insecure.setValue(this.config.insecure);
    this.systemActionRetries.setValue(this.config.systemActionRetries);
    this.systemRebootDelay.setValue(this.config.systemRebootDelay);
    this.type.setValue(this.config.type);
    this.useproxy.setValue(this.config.useproxy);
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

  setManagementConfig(): void {
    const msg = new WsMessage(this.msgType, this.component, WsConstants.SET_MANAGEMENT_CONFIG);
    msg.name = this.name.value;

    const cfg: ManagementConfig = {
      name: this.name.value,
      insecure: this.insecure.value,
      // TODO(mfuller): need to validate these are numerical values in the form
      systemActionRetries: +this.systemActionRetries.value,
      systemRebootDelay: +this.systemRebootDelay.value,
      type: this.type.value,
      useproxy: this.useproxy.value
    };

    msg.data = JSON.parse(JSON.stringify(cfg));
    this.websocketService.sendMessage(msg);
    this.toggleLock();
  }

}
