import { Component, OnDestroy, OnInit } from '@angular/core';
import { NavInterface } from './models/nav.interface';
import { environment } from '../environments/environment';
import { IconService } from '../services/icon/icon.service';
import { WebsocketService } from '../services/websocket/websocket.service';
import { WSReceiver } from '../services/websocket/websocket.models';
import { Dashboard } from '../services/websocket/models/websocket-message/dashboard/dashboard';
import { WebsocketMessage } from 'src/services/websocket/models/websocket-message/websocket-message';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit, WSReceiver {
  type: string = "ui";
  component: string = "any";

  currentYear: number;
  version: string;

  menu: NavInterface [] = [
    {
      displayName: 'Airship',
      iconName: 'airplane',
      children: [
        {
          displayName: 'Bare Metal',
          route: 'ctl/baremetal',
          iconName: 'server'
        }, {
          displayName: 'Documents',
          route: 'ctl/documents',
          iconName: 'doc'
        }]
    }, {
      displayName: 'Dashboards',
      iconName: 'speed',
    }];

  constructor(private iconService: IconService,
              private websocketService: WebsocketService) {
    this.currentYear = new Date().getFullYear();
    this.version = environment.version;
    this.websocketService.registerFunctions(this);
  }

  async receiver(message: WebsocketMessage): Promise<void> {
    if (message.hasOwnProperty("error")) {
      this.websocketService.printIfToast(message);
    } else {
      if (message.hasOwnProperty("dashboards")) {
        this.updateDashboards(message.dashboards);
      } else {
        // TODO (aschiefe): determine what should be notifications and what should be 86ed
        console.log("Message received in app: ", message);
      }
    }
  }

  ngOnInit(): void {
    this.iconService.registerIcons();
  }

  updateDashboards(dashboards: Dashboard[]): void {
    if (this.menu[1].children === undefined) {
      this.menu[1].children = [];
    }
    dashboards.forEach((dashboard) => {
      const navInterface = new NavInterface();
      navInterface.displayName = dashboard.name;
      navInterface.route = dashboard.baseURL;
      navInterface.external = true;
      this.menu[1].children.push(navInterface);
    });
  }

  openLink(url: string): void {
    window.open(url, '_blank');
  }
}
