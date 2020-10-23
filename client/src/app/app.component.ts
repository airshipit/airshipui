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

import { Component, OnInit, ViewChild } from '@angular/core';
import { MatAccordion } from '@angular/material/expansion';
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
  @ViewChild(MatAccordion) accordion: MatAccordion;

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
          iconName: 'devices'
        }, {
          displayName: 'Cluster',
          route: 'ctl/cluster',
          iconName: 'server'
        } , {
          displayName: 'Config',
          route: 'ctl/config',
          iconName: 'settings'
        }, {
          displayName: 'Documents',
          route: 'ctl/documents',
          iconName: 'doc'
        }, {
          displayName: 'Image',
          route: 'ctl/image',
          iconName: 'camera'
        }, {
          displayName: 'Phase',
          route: 'ctl/phase',
          iconName: 'list_alt'
        }, {
          displayName: 'Secret',
          route: 'ctl/secret',
          iconName: 'security'
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
      switch (message.component) {
        case 'log':
          Log.Debug(new LogMessage('Log message received in app', this.className, message));
          const panel = document.getElementById('logPanel');
          panel.appendChild(document.createTextNode(message.message));
          panel.appendChild(document.createElement('br'));
          break;
        case 'initialize':
          Log.Debug(new LogMessage('Initialize message received in app', this.className, message));
          if (message.hasOwnProperty('dashboards')) {
            this.updateDashboards(message.dashboards);
          }
          break;
        case 'keepalive':
          Log.Debug(new LogMessage('Keepalive message received in app', this.className, message));
          break;
        default:
          Log.Debug(new LogMessage('Uncategorized message received in app', this.className, message));
          break;
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
