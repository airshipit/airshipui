import { Component, OnInit } from '@angular/core';
import { WebsocketService } from '../../../../services/websocket/websocket.service';
import { WebsocketMessage } from '../../../../services/websocket/models/websocket-message/websocket-message';

@Component({
  selector: 'app-document-pull',
  templateUrl: './document-pull.component.html',
  styleUrls: ['./document-pull.component.css']
})
export class DocumentPullComponent implements OnInit {

  obby: string;

  constructor(private websocketService: WebsocketService) {

  }

  ngOnInit(): void {
    this.websocketService.subject.subscribe(message => {
        if (message.type === 'airshipctl' && message.component === 'document' && message.subComponent === 'docPull') {
          this.obby = JSON.stringify(message);
        }
      }
    );
  }

  documentPull(): void {
    const websocketMessage = new WebsocketMessage();
    websocketMessage.type = 'airshipctl';
    websocketMessage.component = 'document';
    websocketMessage.subComponent = 'docPull';
    this.websocketService.sendMessage(websocketMessage);
  }

}
