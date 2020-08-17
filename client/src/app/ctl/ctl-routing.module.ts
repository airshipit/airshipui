import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {DocumentComponent} from './document/document.component';
import {BaremetalComponent} from './baremetal/baremetal.component';

const routes: Routes = [{
    path: 'documents',
    component: DocumentComponent,
}, {
    path: 'baremetal',
    component: BaremetalComponent
}];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class CtlRoutingModule {}
