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

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ConfigComponent } from './config.component';
import { ToastrModule } from 'ngx-toastr';
import { CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { ConfigContextModule } from './config-context/config-context.module';
import { ConfigManagementModule } from './config-management/config-management.module';
import { ConfigManifestModule } from './config-manifest/config-manifest.module';
import { ConfigEncryptionModule } from './config-encryption/config-encryption.module';
import { ConfigManifestComponent } from './config-manifest/config-manifest.component';
import { ConfigManagementComponent } from './config-management/config-management.component';
import { ConfigEncryptionComponent } from './config-encryption/config-encryption.component';
import { ConfigContextComponent } from './config-context/config-context.component';
import { MatExpansionModule } from '@angular/material/expansion';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatDialogModule } from '@angular/material/dialog';

describe('ConfigComponent', () => {
  let component: ConfigComponent;
  let fixture: ComponentFixture<ConfigComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        ToastrModule.forRoot(),
        FormsModule,
        MatButtonModule,
        MatInputModule,
        MatCheckboxModule,
        ConfigContextModule,
        ConfigManagementModule,
        ConfigManifestModule,
        ConfigEncryptionModule,
        ReactiveFormsModule,
        MatExpansionModule,
        BrowserAnimationsModule,
        MatDialogModule
      ],
      declarations: [
        ConfigComponent,
        ConfigManifestComponent,
        ConfigManagementComponent,
        ConfigEncryptionComponent,
        ConfigContextComponent
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ConfigComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
