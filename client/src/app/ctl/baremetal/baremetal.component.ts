import {Component} from '@angular/core';
import {WebsocketService} from '../../../services/websocket/websocket.service';
import { WebsocketMessage, WSReceiver } from '../../../services/websocket/websocket.models';
import { Log } from '../../../services/log/log.service';
import { LogMessage } from '../../../services/log/log-message';

@Component({
  selector: 'app-bare-metal',
  templateUrl: './baremetal.component.html',
})

export class BaremetalComponent implements WSReceiver {
  className = this.constructor.name;
  // TODO (aschiefe): extract these strings to constants
  type = 'ctl';
  component = 'baremetal';

  constructor(private websocketService: WebsocketService) {
    this.websocketService.registerFunctions(this);
  }

  async receiver(message: WebsocketMessage): Promise<void> {
    if (message.hasOwnProperty('error')) {
      this.websocketService.printIfToast(message);
    } else {
      // TODO (aschiefe): determine what should be notifications and what should be 86ed
      Log.Debug(new LogMessage('Message received in baremetal', this.className, message));
    }
  }

  generateIso(): void {
    this.websocketService.sendMessage(new WebsocketMessage(this.type, this.component, 'generateISO'));
  }
}
