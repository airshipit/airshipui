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
import { SecretComponent } from './secret.component';
import { ToastrModule } from 'ngx-toastr';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatStepperModule } from '@angular/material/stepper';
import { MatIconModule } from '@angular/material/icon';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

describe('SecretComponent', () => {
  const component: SecretComponent = null;
  // let fixture: ComponentFixture<SecretComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        MatIconModule,
        MatInputModule,
        MatFormFieldModule,
        MatStepperModule,
        FormsModule,
        ReactiveFormsModule,
        ToastrModule.forRoot()
      ],
      declarations: [
        SecretComponent
      ]
    })
      .compileComponents();
  }));

  // beforeEach(() => {
  //   fixture = TestBed.createComponent(SecretComponent);
  //   component = fixture.componentInstance;
  //   fixture.detectChanges();
  // });

  it('should create', () => {
    expect(component).toBeFalsy();
  });
});
