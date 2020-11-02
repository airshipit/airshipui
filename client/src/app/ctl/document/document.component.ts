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
  selector: 'app-document',
  templateUrl: './document.component.html',
  styleUrls: ['./document.component.css']
})

export class DocumentComponent implements WsReceiver {
  className = this.constructor.name;
  statusMsg: string;
  type = WsConstants.CTL;
  component = WsConstants.DOCUMENT;
  activeLink = 'overview';

  constructor(private websocketService: WsService) {
    this.websocketService.registerFunctions(this);
  }

  public async receiver(message: WsMessage): Promise<void> {
    if (message.hasOwnProperty(WsConstants.ERROR)) {
      this.websocketService.printIfToast(message);
    } else {
      switch (message.subComponent) {
        case WsConstants.PULL:
          this.statusMsg = 'Document pull was a ' + message.message;
          const button = (document.getElementById('DocPullBtn') as HTMLInputElement);
          button.removeAttribute('disabled');
          break;
        default:
          Log.Error(new LogMessage('Document message sub component not handled', this.className, message));
          break;
      }
    }
  }

  documentPull(): void {
    this.statusMsg = '';
    this.websocketService.sendMessage(new WsMessage(this.type, this.component, WsConstants.PULL));
    const button = (document.getElementById('DocPullBtn') as HTMLInputElement);
    button.setAttribute('disabled', 'disabled');
  }
}
