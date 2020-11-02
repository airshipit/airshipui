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
  selector: 'app-cluster',
  templateUrl: './cluster.component.html',
  styleUrls: ['./cluster.component.css']
})

export class ClusterComponent implements WsReceiver {
  className = this.constructor.name;
  type = WsConstants.CTL;
  component = WsConstants.CLUSTER;

  constructor(private websocketService: WsService) {
    this.websocketService.registerFunctions(this);
  }

  async receiver(message: WsMessage): Promise<void> {
    if (message.hasOwnProperty(WsConstants.ERROR)) {
      this.websocketService.printIfToast(message);
    } else {
      switch (message.subComponent) {
        default:
          Log.Error(new LogMessage('Cluster message sub component not handled', this.className, message));
          break;
      }
    }
  }
}
