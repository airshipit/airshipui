import { WebsocketMessage } from '../websocket/websocket.models';

export class LogMessage {
    // the holy trinity of the websocket messages, a triumvirate if you will, which is how all are routed
    message: string;
    className: string;
    wsMessage: WebsocketMessage;

    constructor(message?: string | undefined, className?: string | undefined, wsMessage?: WebsocketMessage | undefined) {
        this.message = message;
        this.className = className;
        this.wsMessage = wsMessage;
    }
}
