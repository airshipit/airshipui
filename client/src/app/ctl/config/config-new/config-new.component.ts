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

import { Component, Inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { WsService } from 'src/services/ws/ws.service';
import { FormControl } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { WsConstants, WsMessage } from 'src/services/ws/ws.models';

@Component({
  selector: 'app-config-new',
  templateUrl: './config-new.component.html',
  styleUrls: ['./config-new.component.css']
})
export class ConfigNewComponent implements OnInit {
  group: FormGroup;
  subComponent: string;
  encryptionType: string;
  checkoutTypes = ['Branch', 'Tag', 'CommitHash'];
  checkoutType = 'Branch';

  constructor(private websocketService: WsService,
              private fb: FormBuilder,
              @Inject(MAT_DIALOG_DATA) public data: {
                formName: string,
                configs: {}
              },
              public dialogRef: MatDialogRef<ConfigNewComponent>) { }

  ngOnInit(): void {
    switch (this.data.formName) {
      case 'context':
        this.group = this.fb.group({
          Name: new FormControl('', Validators.required),
          Manifest: new FormControl(''),
          EncryptionConfig: new FormControl(''),
          ManagementConfiguration: new FormControl('')
        });
        this.subComponent = WsConstants.SET_CONTEXT;
        break;
      case 'manifest':
        this.group = this.fb.group({
          Name: new FormControl('', Validators.required),
          TargetPath: new FormControl('', Validators.required),
          MetadataPath: new FormControl('', Validators.required),
          // new manifests seem to get an auto-generated repo named 'primary'
          // that won't get configured properly unless it's done here, so
          // don't let users modify this field
          RepoName: new FormControl({value: 'primary', disabled: true}),
          URL: new FormControl('', Validators.required),
          Tag: new FormControl(''),
          CommitHash: new FormControl(''),
          Branch: new FormControl(''),
          IsPhase: new FormControl(false),
          Force: new FormControl(false)
        });
        this.subComponent = WsConstants.SET_MANIFEST;
        break;
      case 'encryption':
        this.group = this.fb.group({
          Name: new FormControl('', Validators.required),
          EncryptionKeyPath: new FormControl('', Validators.required),
          DecryptionKeyPath: new FormControl('', Validators.required),
          KeySecretName: new FormControl('', Validators.required),
          KeySecretNamespace: new FormControl('', Validators.required),
        });
        this.subComponent = WsConstants.SET_ENCRYPTION_CONFIG;
        break;
      case 'management':
        // NOTE: capitalizations are different for management config due to
        // inconsistent json definitions in airshipctl
        this.group = this.fb.group({
          Name: new FormControl('', Validators.required),
          type: new FormControl(''),
          insecure: new FormControl(false),
          useproxy: new FormControl(false),
          systemActionRetries: new FormControl(0, Validators.pattern('^[0-9]*$')),
          systemRebootDelay: new FormControl(0, Validators.pattern('^[0-9]*$'))
        });
        this.subComponent = WsConstants.SET_MANAGEMENT_CONFIG;
        break;
    }
  }

  setConfig(): void {
    const msg = new WsMessage(WsConstants.CTL, WsConstants.CONFIG, this.subComponent);
    const opts = {};
    for (const [key, val] of Object.entries(this.group.controls)) {
      if (key === 'systemActionRetries' || key === 'systemRebootDelay') {
        opts[key] = +val.value;
      } else {
        opts[key] = val.value;
      }
    }
    const name = 'Name';
    msg.name = opts[name];
    msg.data = JSON.parse(JSON.stringify(opts));

    this.websocketService.sendMessage(msg);
    this.closeDialog();
  }

  closeDialog(): void {
    this.dialogRef.close();
  }

  onEncryptionChange(event: any): void {
    if (this.encryptionType === 'encryption') {
      this.group.controls.EncryptionKeyPath.enable();
      this.group.controls.DecryptionKeyPath.enable();
      this.group.controls.KeySecretName.disable();
      this.group.controls.KeySecretNamespace.disable();
    } else {
      this.group.controls.EncryptionKeyPath.disable();
      this.group.controls.DecryptionKeyPath.disable();
      this.group.controls.KeySecretName.enable();
      this.group.controls.KeySecretNamespace.enable();
    }
  }

}
