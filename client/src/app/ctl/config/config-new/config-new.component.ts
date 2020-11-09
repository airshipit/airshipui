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
import { FormBuilder, FormGroup } from '@angular/forms';
import { WsService } from 'src/services/ws/ws.service';
import { FormControl } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ContextOptions, EncryptionConfigOptions, ManagementConfig, ManifestOptions } from '../config.models';
import { WsConstants, WsMessage } from 'src/services/ws/ws.models';

@Component({
  selector: 'app-config-new',
  templateUrl: './config-new.component.html',
  styleUrls: ['./config-new.component.css']
})
export class ConfigNewComponent implements OnInit {
  group: FormGroup;

  dataObj: any;
  keys: string[] = [];

  dataObjs = {
    context: new ContextOptions(),
    manifest: new ManifestOptions(),
    encryption: new EncryptionConfigOptions(),
    management: new ManagementConfig()
  };

  constructor(private websocketService: WsService,
              private fb: FormBuilder,
              @Inject(MAT_DIALOG_DATA) public data: {formName: string},
              public dialogRef: MatDialogRef<ConfigNewComponent>) { }

  ngOnInit(): void {
    const grp = {};
    this.dataObj = this.dataObjs[this.data.formName];

    for (const [key, val] of Object.entries(this.dataObj)) {
      this.keys.push(key);
      grp[key] = new FormControl(val);
    }

    this.group = new FormGroup(grp);
  }

  setConfig(type: string): void {
    let subComponent = '';

    switch (type) {
      case 'context':
        subComponent = WsConstants.SET_CONTEXT;
        break;
      case 'manifest':
        subComponent = WsConstants.SET_MANIFEST;
        break;
      case 'encryption':
        subComponent = WsConstants.SET_ENCRYPTION_CONFIG;
        break;
      case 'management':
        subComponent = WsConstants.SET_MANAGEMENT_CONFIG;
        break;
    }

    for (const [key, control] of Object.entries(this.group.controls)) {
      // TODO(mfuller): need to validate this within the form
      if (typeof this.dataObj[key] === 'number') {
        this.dataObj[key] = +control.value;
      } else {
        this.dataObj[key] = control.value;
      }
    }

    const msg = new WsMessage(WsConstants.CTL, WsConstants.CONFIG, subComponent);
    msg.data = JSON.parse(JSON.stringify(this.dataObj));
    msg.name = this.dataObj.Name;

    this.websocketService.sendMessage(msg);
    this.dialogRef.close();
  }

  closeDialog(): void {
    this.dialogRef.close();
  }

  // annoying helper method because apparently I can't just test this natively
  // inside an *ngIf
  isBool(val: any): boolean {
    return typeof val === 'boolean';
  }

}
