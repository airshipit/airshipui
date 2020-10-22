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

import { Component, OnInit, ViewChild } from '@angular/core';
import { WebsocketService } from 'src/services/websocket/websocket.service';
import { WebsocketMessage, WSReceiver } from 'src/services/websocket/websocket.models';
import { Log } from 'src/services/log/log.service';
import { LogMessage } from 'src/services/log/log-message';
import { MatStepper } from '@angular/material/stepper';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import {STEPPER_GLOBAL_OPTIONS, StepperSelectionEvent } from '@angular/cdk/stepper';

@Component({
  selector: 'app-secret',
  templateUrl: './secret.component.html',
  styleUrls: ['./secret.component.css'],
  providers: [{
    provide: STEPPER_GLOBAL_OPTIONS, useValue: {showError: true}
  }]
})

export class SecretComponent implements WSReceiver, OnInit {
  className = this.constructor.name;
  // TODO (aschiefe): extract these strings to constants
  type = 'ctl';
  component = 'secret';

  // form groups, these control the stepper and the validators for the inputs
  decryptSrcFG: FormGroup;
  decryptDestFG: FormGroup;
  encryptSrcFG: FormGroup;
  encryptDestFG: FormGroup;

  @ViewChild('decryptStepper', { static: false }) decryptStepper: MatStepper;
  @ViewChild('encryptStepper', { static: false }) encryptStepper: MatStepper;

  constructor(private websocketService: WebsocketService, private formBuilder: FormBuilder) {
    this.websocketService.registerFunctions(this);
  }

  ngOnInit(): void {
    this.decryptSrcFG = this.formBuilder.group({
      decryptSrcCtrl: ['', Validators.required]
    });
    this.decryptDestFG = this.formBuilder.group({
      decryptDestCtrl: ['', Validators.required]
    });
    this.encryptSrcFG = this.formBuilder.group({
      encryptSrcCtrl: ['', Validators.required]
    });
    this.encryptDestFG = this.formBuilder.group({
      encryptDestCtrl: ['', Validators.required]
    });
  }

  async receiver(message: WebsocketMessage): Promise<void> {
    if (message.hasOwnProperty('error')) {
      this.websocketService.printIfToast(message);
    } else {
      switch (message.subComponent) {
        case 'generate':
          document.getElementById('GenerateOutputDiv').innerHTML = message.message;
          break;
        default:
          Log.Error(new LogMessage('Secret message sub component not handled', this.className, message));
          break;
      }
    }
  }

  decrypt(): void {
    const message = new WebsocketMessage(this.type, this.component, 'decrypt');
    message.message = 'Decrypt is currently not implemented in CTL';
    this.websocketService.printIfToast(message);
    this.decryptStepper.reset();
  }

  encrypt(): void {
    const message = new WebsocketMessage(this.type, this.component, 'encrypt');
    message.message = 'Encrypt is currently not implemented in CTL';
    this.websocketService.printIfToast(message);
    this.encryptStepper.reset();
  }

  generateSecret(): void {
    const message = new WebsocketMessage(this.type, this.component, 'generate');
    Log.Debug(new LogMessage('Attempting to generate secret', this.className, message));
    this.websocketService.sendMessage(message);
  }

  // replace the destination with a pre populated value based on the src
  decryptChange(event: StepperSelectionEvent): void {
    if (event.selectedIndex === 1) {
      const src = (document.getElementById('decryptSrc') as HTMLInputElement).value;
      const filename = src.split('/').pop();
      const newFilename = filename.replace('encrypted-', '');
      this.decryptDestFG.controls.decryptDestCtrl.setValue(src.replace(filename, newFilename));
    } else if (event.selectedIndex === 2) {
      const src = (document.getElementById('decryptSrc') as HTMLInputElement).value;
      const dest = (document.getElementById('decryptSrc') as HTMLInputElement).value;
      document.getElementById('decryptSrcTd').innerHTML = src;
      document.getElementById('decryptDestTd').innerHTML = dest;
    }
  }

  // replace the destination with a pre populated value based on the src
  encryptChange(event: StepperSelectionEvent): void {
    if (event.selectedIndex === 1) {
      const src = (document.getElementById('encryptSrc') as HTMLInputElement).value;
      const filename = src.split('/').pop();
      const newFilename = 'encrypted-' + filename;
      this.encryptDestFG.controls.encryptDestCtrl.setValue(src.replace(filename, newFilename));
    } else if (event.selectedIndex === 2) {
      const src = (document.getElementById('encryptSrc') as HTMLInputElement).value;
      const dest = (document.getElementById('encryptSrc') as HTMLInputElement).value;
      document.getElementById('encryptSrcTd').innerHTML = src;
      document.getElementById('encryptDestTd').innerHTML = dest;
    }
  }
}
