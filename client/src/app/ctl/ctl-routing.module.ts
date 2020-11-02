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

import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { BaremetalComponent } from './baremetal/baremetal.component';
import { ClusterComponent } from './cluster/cluster.component';
import { ConfigComponent } from './config/config.component';
import { DocumentComponent } from './document/document.component';
import { ImageComponent } from './image/image.component';
import { PhaseComponent } from './phase/phase.component';
import { SecretComponent } from './secret/secret.component';
import { WsConstants } from 'src/services/ws/ws.models';
import { AuthGuard } from 'src/services/auth-guard/auth-guard.service';

const routes: Routes = [{
  path: WsConstants.BAREMETAL,
  canActivate: [AuthGuard],
  component: BaremetalComponent
}, {
  path: WsConstants.CLUSTER,
  canActivate: [AuthGuard],
  component: ClusterComponent,
}, {
  path: WsConstants.CONFIG,
  canActivate: [AuthGuard],
  component: ConfigComponent,
}, {
  path: WsConstants.DOCUMENTS,
  canActivate: [AuthGuard],
  component: DocumentComponent,
}, {
  path: WsConstants.IMAGE,
  canActivate: [AuthGuard],
  component: ImageComponent,
}, {
  path: WsConstants.PHASE,
  canActivate: [AuthGuard],
  component: PhaseComponent,
}, {
  path: WsConstants.SECRET,
  canActivate: [AuthGuard],
  component: SecretComponent,
}];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class CtlRoutingModule { }
