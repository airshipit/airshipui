import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { DashboardsComponent } from './dashboards/dashboards.component';
import { CTLComponent } from './ctl/ctl.component';
import { BareMetalComponent } from './ctl/baremetal/baremetal.component';
import { DocumentComponent } from './ctl/document/document.component';

const routes: Routes = [
  {
    path: 'ctl',
    component: CTLComponent,
    children: [
      {
        path: 'baremetal',
        component: BareMetalComponent
      }, {
        path: 'documents',
        component: DocumentComponent
      }]
  }, {
    path: 'dashboard',
    component: DashboardsComponent
  }, {
    path: '',
    component: HomeComponent
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {

}
