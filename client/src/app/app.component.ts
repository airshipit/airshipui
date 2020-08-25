import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';
import { IconService } from 'src/services/icon/icon.service';
import { WebsocketService } from 'src/services/websocket/websocket.service';
import { Log } from 'src/services/log/log.service';
import { LogMessage } from 'src/services/log/log-message';
import { Dashboard, WSReceiver, WebsocketMessage } from 'src/services/websocket/websocket.models';
import { Nav } from './app.models';
import { AuthGuard } from 'src/services/auth-guard/auth-guard.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})

export class AppComponent implements OnInit, WSReceiver {
  className = this.constructor.name;
  type = 'ui';
  component = 'any';

  currentYear: number;
  version: string;

  menu: Nav [] = [
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
    if (message.hasOwnProperty('error')) {
      this.websocketService.printIfToast(message);
    } else {
      if (message.hasOwnProperty('dashboards')) {
        this.updateDashboards(message.dashboards);
      } else {
        // TODO (aschiefe): determine what should be notifications and what should be 86ed
        Log.Debug(new LogMessage('Message received in app', this.className, message));
      }
    }
  }

  public authToggle(): void {
    const button = document.getElementById('loginButton');

    if (button.innerText === 'Logout') {
      AuthGuard.logout();
      button.innerText = 'Login';
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
      const navInterface = new Nav();
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
