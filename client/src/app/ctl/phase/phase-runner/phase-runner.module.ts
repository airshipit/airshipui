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

import {NgModule} from '@angular/core';
import {PhaseRunnerComponent} from '../phase-runner/phase-runner.component';
import { MatDialogModule } from '@angular/material/dialog';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import {MatCheckboxModule} from '@angular/material/checkbox';
import { CommonModule } from '@angular/common';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';

@NgModule({
  declarations: [
    PhaseRunnerComponent
  ],
  imports: [
    CommonModule,
    MatDialogModule,
    FormsModule,
    MatCheckboxModule,
    ReactiveFormsModule,
    MatInputModule,
    MatButtonModule
  ],
  providers: [],
})
export class PhaseRunnerModule {}
