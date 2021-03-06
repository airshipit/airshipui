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

import { WsMessage } from '../ws/ws.models';

export class LogMessage {
    // the holy trinity of the websocket messages, a triumvirate if you will, which is how all are routed
    message: string;
    className: string;
    logMessage: string | WsMessage;

    constructor(message?: string | undefined, className?: string | undefined, logMessage?: string | WsMessage | undefined) {
        this.message = message;
        this.className = className;
        this.logMessage = logMessage;
    }
}
