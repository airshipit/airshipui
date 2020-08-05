export interface WSReceiver {
    // the holy trinity of the websocket messages, a triumvirate if you will, which is how all are routed
    type: string;
    component: string;

    // This is the method which will need to be implemented in the component to handle the messages
    receiver(message: WebsocketMessage): Promise<void>;
}

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
  constructor(type?: string | undefined, component?: string | undefined, subComponent?: string | undefined) {
    this.type = type;
    this.component = component;
    this.subComponent = subComponent;
  }
}

export class Dashboard {
  name: string;
  baseURL: string;
  path: string;
  isProxied: boolean;
}
