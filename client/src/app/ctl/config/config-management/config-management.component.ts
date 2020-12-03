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
import { FormControl, FormGroup, Validators, AbstractControl } from '@angular/forms';
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

  group: FormGroup;

  constructor(private websocketService: WsService) { }

  ngOnInit(): void {
    this.group = new FormGroup({
      name: new FormControl({value: this.config.Name, disabled: true}),
      insecure: new FormControl({value: this.config.insecure, disabled: true}),
      systemActionRetries: new FormControl({value: this.config.systemActionRetries, disabled: true},
        Validators.pattern('^[0-9]*$')),
      systemRebootDelay: new FormControl({value: this.config.systemRebootDelay, disabled: true},
        Validators.pattern('^[0-9]*$')),
      type: new FormControl({value: this.config.type, disabled: true}),
      useproxy: new FormControl({value: this.config.useproxy, disabled: true})
    });
  }

  toggleLock(): void {
    if (this.group.disabled) {
      this.group.enable();
    } else {
      this.group.disable();
    }
    this.locked = !this.locked;
  }

  setManagementConfig(): void {
    const msg = new WsMessage(this.msgType, this.component, WsConstants.SET_MANAGEMENT_CONFIG);
    msg.name = this.group.controls.name.value;

    const cfg: ManagementConfig = {
      Name: this.group.controls.name.value,
      insecure: this.group.controls.insecure.value,
      systemActionRetries: +this.group.controls.systemActionRetries.value,
      systemRebootDelay: +this.group.controls.systemRebootDelay.value,
      type: this.group.controls.type.value,
      useproxy: this.group.controls.useproxy.value
    };

    msg.data = JSON.parse(JSON.stringify(cfg));
    this.websocketService.sendMessage(msg);
    this.toggleLock();
  }

}
