import { Component, OnInit } from '@angular/core';
import {WebsocketMessage} from '../../../services/websocket/models/websocket-message/websocket-message';
import {WebsocketService} from '../../../services/websocket/websocket.service';

@Component({
  selector: 'app-bare-metal',
  templateUrl: './bare-metal.component.html',
  styleUrls: ['./bare-metal.component.css']
})
export class BareMetalComponent implements OnInit {

  private message: WebsocketMessage;

  constructor(private websocketService: WebsocketService) {
  }

  ngOnInit(): void { }

  generateIso(): void {
    this.message = new WebsocketMessage();
    this.message.type = 'airshipctl';
    this.message.component = 'baremetal';
    this.message.subComponent = 'generateISO';
    this.websocketService.sendMessage(this.message);
  }
}
