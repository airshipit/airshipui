import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {DocumentComponent} from './document/document.component';
import {BaremetalComponent} from './baremetal/baremetal.component';
import {AuthGuard} from 'src/services/auth-guard/auth-guard.service';

const routes: Routes = [{
    path: 'documents',
    canActivate: [AuthGuard],
    component: DocumentComponent,
  }, {
    path: 'baremetal',
    canActivate: [AuthGuard],
    component: BaremetalComponent
}];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class CtlRoutingModule {}
