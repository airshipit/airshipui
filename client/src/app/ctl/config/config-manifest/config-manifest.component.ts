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

import { Component, Input, OnInit } from '@angular/core';
import { Manifest, ManifestOptions, Repository } from '../config.models';
import { FormControl, FormBuilder, FormGroup, FormArray, Validators } from '@angular/forms';
import { WsService } from 'src/services/ws/ws.service';
import { WsMessage, WsConstants } from 'src/services/ws/ws.models';
import { RepositoryComponent } from './repository/repository.component';
import { MatDialog } from '@angular/material/dialog';


@Component({
  selector: 'app-config-manifest',
  templateUrl: './config-manifest.component.html',
  styleUrls: ['./config-manifest.component.css']
})
export class ConfigManifestComponent implements OnInit {
  @Input() manifest: Manifest;

  type = WsConstants.CTL;
  component = WsConstants.CONFIG;

  locked = true;
  group: FormGroup;
  repoArray = new FormArray([]);
  selectArray: string[] = [];
  selectedIndex = 0;
  checkoutTypes = ['Branch', 'CommitHash', 'Tag'];

  constructor(private websocketService: WsService,
              private fb: FormBuilder,
              public dialog: MatDialog) {
    this.group = this.fb.group({
      name: new FormControl({value: '', disabled: true}),
      repositories: this.repoArray,
      targetPath: new FormControl({value: '', disabled: true}, Validators.required),
      metadataPath: new FormControl({value: '', disabled: true}, Validators.required)
    });
  }

  addRepo(name: string, repo: Repository): void {
    const repoArray = this.group.controls.repositories as FormArray;
    const repoGroup = new FormGroup({
      repoName: new FormControl({value: name, disabled: true}),
      url: new FormControl({value: repo.url, disabled: true}, Validators.required),
      checkoutLabel: new FormControl(''),
      checkoutReference: new FormControl({value: '', disabled: true}, Validators.required),
      force: new FormControl({value: repo.checkout.force, disabled: true}),
      isPhase: new FormControl({value: this.manifest.manifest.phaseRepositoryName === name, disabled: true}),
    });

    const checkout = this.getCheckoutRef(repo);
    repoGroup.controls.checkoutLabel.setValue(checkout[0]);
    repoGroup.controls.checkoutReference.setValue(checkout[1]);
    repoArray.push(repoGroup);
    this.selectArray.push(name);
  }

  getCheckoutRef(repo: Repository): string[] {
    for (const t of this.checkoutTypes) {
      const key = t[0].toLowerCase() + t.substring(1);
      if (repo.checkout[key] !== null && repo.checkout[key] !== '') {
        return [t, repo.checkout[key]];
      }
    }
    return null;
  }

  newRepoDialog(): void {
    const dialogRef = this.dialog.open(RepositoryComponent, {
      width: '400px',
      height: '520px',
      data: {
        name: this.manifest.name,
      }
    });
  }

  ngOnInit(): void {
    this.group.controls.name.setValue(this.manifest.name);
    this.group.controls.targetPath.setValue(this.manifest.manifest.targetPath);
    this.group.controls.metadataPath.setValue(this.manifest.manifest.metadataPath);
    for (const [name, repo] of Object.entries(this.manifest.manifest.repositories)) {
      this.addRepo(name, repo);
    }
  }

  // once set to true, 'isPhase' and 'force' cannot be set to false using airshipctl's
  // setters, so those controls won't be enabled if true. 'isPhase' can only be
  // set to false by setting it to true for another repo
  toggleLock(): void {
    this.toggleControl(this.group.controls.targetPath as FormControl);
    this.toggleControl(this.group.controls.metadataPath as FormControl);
    for (const grp of this.repoArray.controls as FormGroup[]) {
      Object.keys(grp.controls).forEach(key => {
        this.toggleControl(grp.controls[key] as FormControl);
      });
    }
    this.locked = !this.locked;
  }

  toggleControl(ctrl: FormControl): void {
    if (ctrl.disabled && ctrl.value !== true) {
      ctrl.enable();
    } else {
      ctrl.disable();
    }
  }

  setManifest(index: number): void {
    const m = this.repoArray.at(index) as FormGroup;
    const controls = m.controls;
    if (controls !== undefined) {
      const msg = new WsMessage(this.type, this.component, WsConstants.SET_MANIFEST);
      msg.name = this.manifest.name;

      const opts: ManifestOptions = {
        Name: this.manifest.name,
        RepoName: controls.repoName.value,
        URL: controls.url.value,
        Branch: null,
        CommitHash: null,
        Tag: null,
        RemoteRef: null,
        Force: controls.force.value,
        IsPhase: controls.isPhase.value,
        TargetPath: this.group.controls.targetPath.value,
        MetadataPath: this.group.controls.metadataPath.value
      };

      opts[controls.checkoutLabel.value] = controls.checkoutReference.value;

      msg.data = JSON.parse(JSON.stringify(opts));
      this.websocketService.sendMessage(msg);
      this.toggleLock();
    }
  }
}
