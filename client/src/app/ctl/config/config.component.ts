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
import { WebsocketService } from '../../../services/websocket/websocket.service';
import { WebsocketMessage, WSReceiver } from '../../../services/websocket/websocket.models';
import { Log } from '../../../services/log/log.service';
import { LogMessage } from '../../../services/log/log-message';

@Component({
  selector: 'app-bare-metal',
  templateUrl: './config.component.html',
})

export class ConfigComponent implements WSReceiver {
  className = this.constructor.name;
  // TODO (aschiefe): extract these strings to constants
  type = 'ctl';
  component = 'config';

  constructor(private websocketService: WebsocketService) {
    this.websocketService.registerFunctions(this);
  }

  async receiver(message: WebsocketMessage): Promise<void> {
    if (message.hasOwnProperty('error')) {
      this.websocketService.printIfToast(message);
    } else {
      // TODO (aschiefe): determine what should be notifications and what should be 86ed
      Log.Debug(new LogMessage('Message received in config', this.className, message));
    }
  }
}
