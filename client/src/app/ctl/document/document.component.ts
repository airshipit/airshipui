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

import {Component} from '@angular/core';
import {WebsocketService} from '../../../services/websocket/websocket.service';
import {WebsocketMessage, WSReceiver} from '../../../services/websocket/websocket.models';
import {Log} from '../../../services/log/log.service';
import {LogMessage} from '../../../services/log/log-message';

@Component({
  selector: 'app-document',
  templateUrl: './document.component.html',
  styleUrls: ['./document.component.css']
})

export class DocumentComponent implements WSReceiver {
  className = this.constructor.name;
  statusMsg: string;

  type = 'ctl';
  component = 'document';
  activeLink = 'overview';

  constructor(private websocketService: WebsocketService) {
    this.websocketService.registerFunctions(this);
  }

  public async receiver(message: WebsocketMessage): Promise<void> {
    if (message.hasOwnProperty('error')) {
      this.websocketService.printIfToast(message);
    } else {
      switch (message.subComponent) {
        case 'pull':
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
    this.websocketService.sendMessage(new WebsocketMessage(this.type, this.component, 'pull'));
    const button = (document.getElementById('DocPullBtn') as HTMLInputElement);
    button.setAttribute('disabled', 'disabled');
  }
}
