import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { DashboardsComponent } from './dashboards/dashboards.component';
import { AirshipComponent } from './airship/airship.component';
import { BareMetalComponent } from './airship/bare-metal/bare-metal.component';
import { DocumentComponent } from './airship/document/document.component';
import { DocumentOverviewComponent } from './airship/document/document-overview/document-overview.component';
import { DocumentPullComponent } from './airship/document/document-pull/document-pull.component';


const routes: Routes = [
  {
    path: 'airship',
    component: AirshipComponent,
    children: [
      {
        path: 'bare-metal',
        component: BareMetalComponent
      }, {
        path: 'documents',
        component: DocumentComponent,
        children: [
          {
            path: 'overview',
            component: DocumentOverviewComponent
          }, {
            path: 'pull',
            component: DocumentPullComponent
          }]
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
