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
import { WsService } from 'src/services/ws/ws.service';
import { WsMessage, WsReceiver, WsConstants } from 'src/services/ws/ws.models';
import { Log } from 'src/services/log/log.service';
import { LogMessage } from 'src/services/log/log-message';
import { MatStepper } from '@angular/material/stepper';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { STEPPER_GLOBAL_OPTIONS, StepperSelectionEvent } from '@angular/cdk/stepper';

@Component({
  selector: 'app-secret',
  templateUrl: './secret.component.html',
  styleUrls: ['./secret.component.css'],
  providers: [{
    provide: STEPPER_GLOBAL_OPTIONS, useValue: {showError: true}
  }]
})

export class SecretComponent implements WsReceiver, OnInit {
  className = this.constructor.name;
  type = WsConstants.CTL;
  component = WsConstants.SECRET;

  // form groups, these control the stepper and the validators for the inputs
  decryptSrcFG: FormGroup;
  decryptDestFG: FormGroup;
  encryptSrcFG: FormGroup;
  encryptDestFG: FormGroup;

  @ViewChild('decryptStepper', { static: false }) decryptStepper: MatStepper;
  @ViewChild('encryptStepper', { static: false }) encryptStepper: MatStepper;

  constructor(private websocketService: WsService, private formBuilder: FormBuilder) {
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

  async receiver(message: WsMessage): Promise<void> {
    if (message.hasOwnProperty(WsConstants.ERROR)) {
      this.websocketService.printIfToast(message);
    } else {
      switch (message.subComponent) {
        case WsConstants.GENERATE:
          document.getElementById('GenerateOutputDiv').innerHTML = message.message;
          break;
        default:
          Log.Error(new LogMessage('Secret message sub component not handled', this.className, message));
          break;
      }
    }
  }

  decrypt(): void {
    const message = new WsMessage(this.type, this.component, WsConstants.DECRYPT);
    message.message = 'Decrypt is currently not implemented in CTL';
    this.websocketService.printIfToast(message);
    this.decryptStepper.reset();
  }

  encrypt(): void {
    const message = new WsMessage(this.type, this.component, WsConstants.ENCRYPT);
    message.message = 'Encrypt is currently not implemented in CTL';
    this.websocketService.printIfToast(message);
    this.encryptStepper.reset();
  }

  generateSecret(): void {
    const message = new WsMessage(this.type, this.component, WsConstants.GENERATE);
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
