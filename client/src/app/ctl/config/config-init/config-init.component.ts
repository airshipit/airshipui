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

import { Component } from '@angular/core';
import { FormControl } from '@angular/forms';
import { WsService } from 'src/services/ws/ws.service';
import { WsMessage, WsConstants } from 'src/services/ws/ws.models';

@Component({
  selector: 'app-config-init',
  templateUrl: './config-init.component.html',
  styleUrls: ['./config-init.component.css']
})
export class ConfigInitComponent {
  type = WsConstants.CTL;
  component = WsConstants.CONFIG;

  initValue = new FormControl('');
  specifyValue = new FormControl('');

  constructor(private websocketService: WsService) {}

  initAirshipConfig(): void {
    const msg = new WsMessage(this.type, this.component, WsConstants.INIT);
    msg.message = this.initValue.value;
    this.websocketService.sendMessage(msg);
  }

  setAirshipConfig(): void {
    const msg = new WsMessage(this.type, this.component, WsConstants.SET_AIRSHIP_CONFIG);
    msg.message = this.specifyValue.value;
    this.websocketService.sendMessage(msg);
  }

}
