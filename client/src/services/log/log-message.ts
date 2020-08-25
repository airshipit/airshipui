import { WebsocketMessage } from '../websocket/websocket.models';

export class LogMessage {
    // the holy trinity of the websocket messages, a triumvirate if you will, which is how all are routed
    message: string;
    className: string;
    logMessage: string | WebsocketMessage;

    constructor(message?: string | undefined, className?: string | undefined, logMessage?: string | WebsocketMessage | undefined) {
        this.message = message;
        this.className = className;
        this.logMessage = logMessage;
    }
}
