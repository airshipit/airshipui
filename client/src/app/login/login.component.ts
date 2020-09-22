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

import { Component, OnInit } from '@angular/core';
import { WebsocketService } from 'src/services/websocket/websocket.service';
import { WSReceiver, WebsocketMessage, Authentication } from 'src/services/websocket/websocket.models';

@Component({
    styleUrls: ['login.component.css'],
    templateUrl: 'login.component.html',
})

export class LoginComponent implements WSReceiver, OnInit {
    className = this.constructor.name;
    type = 'ui'; // needed to have the websocket service in the constructor
    component = 'auth'; // needed to have the websocket service in the constructor

    constructor(private websocketService: WebsocketService) { }

    ngOnInit(): void {
        // bind the enter key to the submit button on the page
        document.getElementById('passwd')
            .addEventListener('keyup', (event) => {
                event.preventDefault();
                if (event.key === 'Enter') {
                    document.getElementById('loginSubmit').click();
                }
            });
    }

    // This will always throw an error but should never be called because we did not register a receiver
    // The auth guard will take care of the auth messages since it's dealing with the tokens
    receiver(message: WebsocketMessage): Promise<void> {
        throw new Error('Method not implemented.');
    }

    // formSubmit sends the auth request to the backend
    public formSubmit(id, passwd): void {
        const message = new WebsocketMessage(this.type, this.component, 'authenticate');
        message.authentication = new Authentication(id, passwd);
        this.websocketService.sendMessage(message);
    }
}
