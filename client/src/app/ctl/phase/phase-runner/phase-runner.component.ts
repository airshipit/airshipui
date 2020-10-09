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

import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { RunOptions } from '../phase.models';
import { RunnerDialogData } from './phase-runner.models';

@Component({
  selector: 'app-phase-runner',
  templateUrl: './phase-runner.component.html',
  styleUrls: ['./phase-runner.component.css']
})
export class PhaseRunnerComponent {
  name: string;
  runOpts: RunOptions;

  checkboxes = [
    { name: 'Debug', checked: false},
    { name: 'DryRun', checked: false}
  ];

  constructor(
    public dialogRef: MatDialogRef<PhaseRunnerComponent>,
    @Inject(MAT_DIALOG_DATA) public data: RunnerDialogData) {
      this.name = data.name;
      this.runOpts = data.options;
    }

  onNoClick(): void {
    this.dialogRef.close();
  }

  getChecked(): void {
    this.runOpts.Debug = this.checkboxes[0].checked;
    this.runOpts.DryRun = this.checkboxes[1].checked;
  }
}
