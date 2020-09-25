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
import {RouterModule, Routes} from '@angular/router';
import {HomeComponent} from './home/home.component';
import {CtlComponent} from './ctl/ctl.component';
import {LoginComponent} from './login/login.component';
import {AuthGuard} from 'src/services/auth-guard/auth-guard.service';

const routes: Routes = [{
    path: 'ctl',
    component: CtlComponent,
    canActivate: [AuthGuard],
    loadChildren: './ctl/ctl.module#CtlModule',
}, {
    path: '',
    canActivate: [AuthGuard],
    component: HomeComponent
}, {
    path: 'login',
    canActivate: [AuthGuard],
    component: LoginComponent
}];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})

export class AppRoutingModule {}
