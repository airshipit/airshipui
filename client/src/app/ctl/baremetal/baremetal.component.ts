import { Component, OnInit } from '@angular/core';
import {WebsocketMessage} from '../../../services/websocket/models/websocket-message/websocket-message';
import {WebsocketService} from '../../../services/websocket/websocket.service';
import { WSReceiver } from '../../../services/websocket//websocket.models';

@Component({
  selector: 'app-bare-metal',
  templateUrl: './baremetal.component.html',
  styleUrls: ['./baremetal.component.css']
})

export class BareMetalComponent implements WSReceiver {
  // TODO (aschiefe): extract these strings to constants
  type: string = "ctl";
  component: string = "baremetal";

  constructor(private websocketService: WebsocketService) {
    this.websocketService.registerFunctions(this);
  }

  async receiver(message: WebsocketMessage): Promise<void> {
    if (message.hasOwnProperty("error")) {
      this.websocketService.printIfToast(message);
    } else {
      // TODO (aschiefe): determine what should be notifications and what should be 86ed
      console.log("Message received in baremetal: ", message);
    }
  }

  generateIso(): void {
    this.websocketService.sendMessage(new WebsocketMessage(this.type, this.component, "generateISO"));
  }
}
