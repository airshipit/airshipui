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
import { FormControl } from '@angular/forms';
import { WsService } from 'src/services/ws/ws.service';
import { WsMessage, WsConstants } from 'src/services/ws/ws.models';


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
  Name = new FormControl({value: '', disabled: true});
  RepoName = new FormControl({value: '', disabled: true});
  URL = new FormControl({value: '', disabled: true});
  Branch = new FormControl({value: '', disabled: true});
  CommitHash = new FormControl({value: '', disabled: true});
  Tag = new FormControl({value: '', disabled: true});
  RemoteRef = new FormControl({value: '', disabled: true});
  Force = new FormControl({value: false, disabled: true});
  IsPhase = new FormControl({value: false, disabled: true});
  SubPath = new FormControl({value: '', disabled: true});
  TargetPath = new FormControl({value: '', disabled: true});
  MetadataPath = new FormControl({value: '', disabled: true});

  controlsArray = [
    this.Name,
    this.RepoName,
    this.URL,
    this.Branch,
    this.CommitHash,
    this.Tag,
    this.RemoteRef,
    this.Force,
    this.IsPhase,
    this.SubPath,
    this.TargetPath,
    this.MetadataPath
  ];

  constructor(private websocketService: WsService) { }

  ngOnInit(): void {
    this.Name.setValue(this.manifest.name);

    // TODO(mfuller): not sure yet how to handle multiple repositories,
    // so for now, I'm just showing the phase repository (primary)
    const repoName = this.manifest.manifest.phaseRepositoryName;
    this.RepoName.setValue(repoName);
    const primaryRepo: Repository = this.manifest.manifest.repositories[repoName];
    this.URL.setValue(primaryRepo.url);
    this.Branch.setValue(primaryRepo.checkout.branch);
    this.CommitHash.setValue(primaryRepo.checkout.commitHash);
    this.Tag.setValue(primaryRepo.checkout.tag);
    this.RemoteRef.setValue(primaryRepo.checkout.remoteRef);
    this.Force.setValue(primaryRepo.checkout.force);
    // TODO(mfuller): this value doesn't come from the config file, but if set to true,
    // it appears to set the phaseRepositoryName key, and since that's
    // the only repo I'm showing, set to true for now
    this.IsPhase.setValue(true);
    this.SubPath.setValue(this.manifest.manifest.subPath);
    this.TargetPath.setValue(this.manifest.manifest.targetPath);
    this.MetadataPath.setValue(this.manifest.manifest.metadataPath);
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

  setManifest(): void {
    const msg = new WsMessage(this.type, this.component, WsConstants.SET_MANIFEST);
    msg.name = this.manifest.name;

    // TODO(mfuller): since "Force" and "IsPhase" can only be set by passing in
    // CLI flags rather than passing in values, there doesn't appear to be a way
    // to unset them once they're true without manually editing the config file.
    // Open a bug for this? Or is this intentional? I may have to write a custom
    // setter to set the value directly in the Config struct
    const opts: ManifestOptions = {
      Name: this.Name.value,
      RepoName: this.RepoName.value,
      URL: this.URL.value,
      Branch: this.Branch.value,
      CommitHash: this.CommitHash.value,
      Tag: this.Tag.value,
      RemoteRef: this.RemoteRef.value,
      Force: this.Force.value,
      IsPhase: this.IsPhase.value,
      SubPath: this.SubPath.value,
      TargetPath: this.TargetPath.value,
      MetadataPath: this.MetadataPath.value
    };

    msg.data = JSON.parse(JSON.stringify(opts));
    this.websocketService.sendMessage(msg);
    this.toggleLock();
  }
}
