import {Dashboard} from './dashboard/dashboard';

export class WebsocketMessage {
  type: string;
  component: string;
  subComponent: string;
  timestamp: number;
  dashboards: Dashboard[];
  error: string;
  fade: boolean;
  html: string;
  name: string;
  isAuthenticated: boolean;
  message: string;
  data: JSON;
  yaml: string;

  // this constructor looks like this in case anyone decides they want just a raw message with no data predefined
  // or an easy way to specify the defaults
  constructor (type?: string | undefined, component?: string | undefined, subComponent?: string | undefined) {
    this.type = type;
    this.component = component;
    this.subComponent = subComponent;
  }
}
