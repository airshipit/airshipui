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

import {NgModule, CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {MatTabsModule} from '@angular/material/tabs';
import {DocumentComponent} from './document.component';
import {MatTreeModule} from '@angular/material/tree';
import {MatButtonModule} from '@angular/material/button';
import {MatButtonToggleModule} from '@angular/material/button-toggle';
import {MatIconModule} from '@angular/material/icon';
import {MonacoEditorModule} from 'ngx-monaco-editor';
import {FormsModule} from '@angular/forms';
import {ToastrModule} from 'ngx-toastr';
import {CommonModule} from '@angular/common';
import {MatProgressBarModule} from '@angular/material/progress-bar';
import {MatCardModule} from '@angular/material/card';
import {MatTooltipModule} from '@angular/material/tooltip/';
import {MatMenuModule} from '@angular/material/menu';
import {DocumentViewerModule} from './document-viewer/document-viewer.module';
import { MatDialogModule } from '@angular/material/dialog';
import { MatListModule } from '@angular/material/list';
import { PhaseRunnerModule } from './phase-runner/phase-runner.module';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

@NgModule({
  declarations: [
    DocumentComponent,
  ],
  imports: [
    CommonModule,
    MatTabsModule,
    MatTreeModule,
    MatButtonModule,
    MatButtonToggleModule,
    MatIconModule,
    MonacoEditorModule,
    FormsModule,
    ToastrModule,
    MatProgressBarModule,
    MatCardModule,
    MatTooltipModule,
    MatMenuModule,
    DocumentViewerModule,
    MatDialogModule,
    MatListModule,
    PhaseRunnerModule,
    MatInputModule,
    MatProgressSpinnerModule
  ],
  providers: [],
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class DocumentModule {}
