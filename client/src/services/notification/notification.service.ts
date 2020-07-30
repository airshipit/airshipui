import { Injectable } from '@angular/core';
import { ToastrService } from 'ngx-toastr';
import { WebsocketService } from '../websocket/websocket.service';
import { WebsocketMessage } from '../websocket/models/websocket-message/websocket-message';

@Injectable({
  providedIn: 'root'
})
export class NotificationService {

  constructor(private toastrService: ToastrService,
              private websocketService: WebsocketService) {
    this.websocketService.subject.subscribe(message => {
        this.printIfToast(message);
      }
    );
  }

  printIfToast(message: WebsocketMessage): void {
    if (message.error !== undefined && message.error !== null) {
      this.toastrService.error(message.error);
    }
  }
}
