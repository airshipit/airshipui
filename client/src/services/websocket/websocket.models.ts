import { WebsocketMessage } from './models/websocket-message/websocket-message';

export interface WSReceiver {
    // the holy trinity of the websocket messages, a triumvirate if you will, which is how all are routed
    type: string;
    component: string;

    // This is the method which will need to be implemented in the component to handle the messages
    receiver(message: WebsocketMessage): Promise<void>;
}