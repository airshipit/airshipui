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
}
