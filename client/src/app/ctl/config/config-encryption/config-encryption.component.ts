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
import { FormControl } from '@angular/forms';
import { EncryptionConfig, EncryptionConfigOptions } from '../config.models';
import { WebsocketService } from '../../../../services/websocket/websocket.service';
import { WebsocketMessage } from '../../../../services/websocket/websocket.models';

@Component({
  selector: 'app-config-encryption',
  templateUrl: './config-encryption.component.html',
  styleUrls: ['./config-encryption.component.css']
})
export class ConfigEncryptionComponent implements OnInit {
  @Input() config: EncryptionConfig;
  type = 'ctl';
  component = 'config';

  locked = true;
  name = new FormControl({value: '', disabled: true});
  encryptionKeyPath = new FormControl({value: '', disabled: true});
  decryptionKeyPath = new FormControl({value: '', disabled: true});
  keySecretName = new FormControl({value: '', disabled: true});
  keySecretNamespace = new FormControl({value: '', disabled: true});

  controlsArray = [this.encryptionKeyPath, this.decryptionKeyPath, this.keySecretName, this.keySecretNamespace];

  constructor(private websocketService: WebsocketService) {}

  ngOnInit(): void {
    this.name.setValue(this.config.name);
    this.encryptionKeyPath.setValue(this.config.encryptionKeyPath);
    this.decryptionKeyPath.setValue(this.config.decryptionKeyPath);
    this.keySecretName.setValue(this.config.keySecretName);
    this.keySecretNamespace.setValue(this.config.keySecretNamespace);
  }

  toggleLock(): void {
    for (const control of this.controlsArray) {
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
      Name: this.name.value,
      EncryptionKeyPath: this.encryptionKeyPath.value,
      DecryptionKeyPath: this.decryptionKeyPath.value,
      KeySecretName: this.keySecretName.value,
      KeySecretNamespace: this.keySecretNamespace.value,
    };

    const msg = new WebsocketMessage(this.type, this.component, 'setEncryptionConfig');
    msg.data = JSON.parse(JSON.stringify(opts));
    msg.name = this.name.value;

    this.websocketService.sendMessage(msg);
    this.toggleLock();
  }

}
