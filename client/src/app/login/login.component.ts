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
import { WsService } from 'src/services/ws/ws.service';
import { WsReceiver, WsMessage, Authentication, WsConstants } from 'src/services/ws/ws.models';

@Component({
    styleUrls: ['login.component.css'],
    templateUrl: 'login.component.html',
})

export class LoginComponent implements WsReceiver, OnInit {
    className = this.constructor.name;
    type = WsConstants.UI;
    component = WsConstants.LOGIN;

    constructor(private websocketService: WsService) { }

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
    receiver(message: WsMessage): Promise<void> {
        throw new Error('Method not implemented.');
    }

    // formSubmit sends the auth request to the backend
    public formSubmit(id, passwd): void {
        const message = new WsMessage(this.type, WsConstants.AUTH, WsConstants.AUTHENTICATE);
        message.authentication = new Authentication(id, passwd);
        this.websocketService.sendMessage(message);
    }
}
