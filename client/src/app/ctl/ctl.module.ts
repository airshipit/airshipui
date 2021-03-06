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
import { RouterModule } from '@angular/router';
import { CtlComponent } from './ctl.component';
import { DocumentModule } from './document/document.module';
import { BaremetalModule } from './baremetal/baremetal.module';
import { ClusterModule } from './cluster/cluster.module';
import { CtlRoutingModule } from './ctl-routing.module';
import { PhaseModule } from './phase/phase.module';
import { SecretModule } from './secret/secret.module';
import { ConfigModule } from './config/config.module';
import { HistoryModule } from './history/history.module';
import { CommonModule } from '@angular/common';

@NgModule({
  imports: [
    CommonModule,
    CtlRoutingModule,
    ClusterModule,
    ConfigModule,
    RouterModule,
    DocumentModule,
    BaremetalModule,
    PhaseModule,
    SecretModule,
    HistoryModule
  ],
  declarations: [CtlComponent],
  providers: []
})
export class CtlModule { }
