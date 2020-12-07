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
import { FormGroup, FormControl, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA} from '@angular/material/dialog';
import { WsConstants, WsMessage } from 'src/services/ws/ws.models';
import { WsService } from 'src/services/ws/ws.service';
import { ManifestOptions } from '../../config.models';

@Component({
  selector: 'app-repository',
  templateUrl: './repository.component.html',
  styleUrls: ['./repository.component.css']
})
export class RepositoryComponent implements OnInit{

  group: FormGroup;
  checkoutTypes = ['Branch', 'CommitHash', 'Tag'];
  checkoutType = 'Branch';

  constructor(
    public dialogRef: MatDialogRef<RepositoryComponent>,
    @Inject(MAT_DIALOG_DATA) public data: {
      name: string,
    },
    private ws: WsService) {}

  ngOnInit(): void {
    this.group = new FormGroup({
      repoName: new FormControl('', Validators.required),
      url: new FormControl('', Validators.required),
      checkoutReference: new FormControl('', Validators.required),
      force: new FormControl(false),
      isPhase: new FormControl(false)
    });
  }

  cancel(): void {
    this.dialogRef.close();
  }

  setRepo(): void {
    const msg = new WsMessage(WsConstants.CTL, WsConstants.CONFIG, WsConstants.SET_MANIFEST);
    msg.name = this.data.name;
    const opts: ManifestOptions = {
      Name: this.data.name,
      RepoName: this.group.controls.repoName.value,
      URL: this.group.controls.url.value,
      Branch: null,
      CommitHash: null,
      Tag: null,
      RemoteRef: null,
      Force: this.group.controls.force.value,
      IsPhase: this.group.controls.isPhase.value,
      TargetPath: null,
      MetadataPath: null
    };

    opts[this.checkoutType] = this.group.controls.checkoutReference.value;

    msg.data = JSON.parse(JSON.stringify(opts));
    this.ws.sendMessage(msg);

    this.dialogRef.close();
  }

}
