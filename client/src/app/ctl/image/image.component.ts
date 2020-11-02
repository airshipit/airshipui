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
import { WsService } from 'src/services/ws/ws.service';
import { WsMessage, WsReceiver, WsConstants } from 'src/services/ws/ws.models';
import { Log } from 'src/services/log/log.service';
import { LogMessage } from 'src/services/log/log-message';

@Component({
  selector: 'app-bare-metal',
  templateUrl: './image.component.html',
  styleUrls: ['./image.component.css']
})

export class ImageComponent implements WsReceiver {
  className = this.constructor.name;
  type = WsConstants.CTL;
  component = WsConstants.IMAGE;
  statusMsg: string;

  constructor(private websocketService: WsService) {
    this.websocketService.registerFunctions(this);
  }

  async receiver(message: WsMessage): Promise<void> {
    if (message.hasOwnProperty(WsConstants.ERROR)) {
      this.websocketService.printIfToast(message);
    } else {
      // TODO (aschiefe): determine what should be notifications and what should be 86ed
      Log.Debug(new LogMessage('Message received in image', this.className, message));
    }
  }

  generateIso(): void {
    this.websocketService.sendMessage(new WsMessage(this.type, this.component, WsConstants.GENERATE));
  }
}
