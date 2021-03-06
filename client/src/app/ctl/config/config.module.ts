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

import { CUSTOM_ELEMENTS_SCHEMA, NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ConfigComponent } from './config.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { MatDividerModule } from '@angular/material/divider';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { ConfigContextComponent } from './config-context/config-context.component';
import { ConfigEncryptionComponent } from './config-encryption/config-encryption.component';
import { ConfigManagementComponent } from './config-management/config-management.component';
import { ConfigManifestComponent } from './config-manifest/config-manifest.component';
import { ConfigInitComponent } from './config-init/config-init.component';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatExpansionModule } from '@angular/material/expansion';
import { ConfigNewComponent } from './config-new/config-new.component';
import { MatSelectModule } from '@angular/material/select';
import { RepositoryComponent } from './config-manifest/repository/repository.component';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    MatCardModule,
    MatInputModule,
    MatDividerModule,
    MatButtonModule,
    MatIconModule,
    MatCheckboxModule,
    MatExpansionModule,
    MatSelectModule
  ],
  declarations: [
    ConfigComponent,
    ConfigContextComponent,
    ConfigEncryptionComponent,
    ConfigManagementComponent,
    ConfigManifestComponent,
    ConfigInitComponent,
    ConfigNewComponent,
    RepositoryComponent
  ],
  providers: [],
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class ConfigModule { }
