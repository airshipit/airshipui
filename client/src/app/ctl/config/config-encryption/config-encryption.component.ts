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

import { Component, OnInit, Input } from '@angular/core';
import { FormControl, FormGroup, Validators, AbstractControl } from '@angular/forms';
import { EncryptionConfig, EncryptionConfigOptions } from '../config.models';
import { WsService } from 'src/services/ws/ws.service';
import { WsMessage, WsConstants } from 'src/services/ws/ws.models';

@Component({
  selector: 'app-config-encryption',
  templateUrl: './config-encryption.component.html',
  styleUrls: ['./config-encryption.component.css']
})
export class ConfigEncryptionComponent implements OnInit {
  @Input() config: EncryptionConfig;
  type = WsConstants.CTL;
  component = WsConstants.CONFIG;

  locked = true;
  group: FormGroup;

  configOptions: AbstractControl[] = [];

  constructor(private websocketService: WsService) {}

  ngOnInit(): void {
    this.group = new FormGroup({
      encryptionKeyPath: new FormControl({value: this.config.encryptionKeyPath, disabled: true},
        Validators.required),
      decryptionKeyPath: new FormControl({value: this.config.decryptionKeyPath, disabled: true},
        Validators.required),
      keySecretName: new FormControl({value: this.config.keySecretName, disabled: true},
        Validators.required),
      keySecretNamespace: new FormControl({value: this.config.keySecretNamespace, disabled: true},
        Validators.required)
      });

    if (this.config.hasOwnProperty('encryptionKeyPath')) {
      this.configOptions.push(this.group.controls.encryptionKeyPath);
      this.configOptions.push(this.group.controls.decryptionKeyPath);
    } else {
      this.configOptions.push(this.group.controls.keySecretName);
      this.configOptions.push(this.group.controls.keySecretNamespace);
    }
  }

  toggleLock(): void {
    for (const control of this.configOptions) {
      if (this.locked) {
        control.enable();
      } else {
        control.disable();
      }
    }

    this.locked = !this.locked;
  }

  setEncryptionConfig(): void {
    const opts: EncryptionConfigOptions = {
      Name: this.config.name,
      EncryptionKeyPath: this.group.controls.encryptionKeyPath.value,
      DecryptionKeyPath: this.group.controls.decryptionKeyPath.value,
      KeySecretName: this.group.controls.keySecretName.value,
      KeySecretNamespace: this.group.controls.keySecretNamespace.value,
    };

    const msg = new WsMessage(this.type, this.component, WsConstants.SET_ENCRYPTION_CONFIG);
    msg.data = JSON.parse(JSON.stringify(opts));
    msg.name = this.config.name;

    this.websocketService.sendMessage(msg);
    this.toggleLock();
  }

}
