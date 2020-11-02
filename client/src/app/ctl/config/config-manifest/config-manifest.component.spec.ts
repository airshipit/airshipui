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
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ToastrModule } from 'ngx-toastr';
import { CtlManifest, Manifest, RepoCheckout, Repository } from '../config.models';

import { ConfigManifestComponent } from './config-manifest.component';

describe('ConfigManifestComponent', () => {
  let component: ConfigManifestComponent;
  let fixture: ComponentFixture<ConfigManifestComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ConfigManifestComponent ],
      imports: [
        BrowserAnimationsModule,
        FormsModule,
        MatInputModule,
        MatIconModule,
        MatCheckboxModule,
        MatButtonModule,
        ReactiveFormsModule,
        ToastrModule.forRoot(),
        MatExpansionModule
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ConfigManifestComponent);
    component = fixture.componentInstance;

    component.manifest = new Manifest();
    component.manifest.manifest = new CtlManifest();
    const repoName = 'fakerepo';
    component.manifest.manifest.phaseRepositoryName = repoName;
    component.manifest.manifest.repositories = {};
    component.manifest.manifest.repositories[repoName] = new Repository();
    component.manifest.manifest.repositories[repoName].checkout = new RepoCheckout();

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
