import { Component, OnDestroy, OnInit } from '@angular/core';
import { NavInterface } from './models/nav.interface';
import { environment } from '../environments/environment';
import { IconService } from '../services/icon/icon.service';
import { NotificationService } from '../services/notification/notification.service';
import {WebsocketService} from '../services/websocket/websocket.service';
import {Dashboard} from '../services/websocket/models/websocket-message/dashboard/dashboard';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnDestroy, OnInit {

  currentYear: number;
  version: string;

  menu: NavInterface [] = [
    {
      displayName: 'Airship',
      iconName: 'airplane',
      children: [
        {
          displayName: 'Bare Metal',
          route: 'airship/bare-metal',
          iconName: 'server'
        }, {
          displayName: 'Documents',
          route: 'airship/documents/overview',
          iconName: 'doc'
        }]
    }, {
      displayName: 'Dashboards',
      iconName: 'speed',
    }];

  constructor(private iconService: IconService,
              private notificationService: NotificationService,
              private websocketService: WebsocketService) {
    this.currentYear = new Date().getFullYear();
    this.version = environment.version;
    this.websocketService.subject.subscribe(message => {
      if (message.type === 'airshipui' && message.component === 'initialize' && message.dashboards !== undefined) {
        this.updateDashboards(message.dashboards);
      }
    });
  }

  ngOnDestroy(): void {
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
