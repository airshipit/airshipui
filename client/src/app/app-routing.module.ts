import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {HomeComponent} from './home/home.component';
import {CtlComponent} from './ctl/ctl.component';


const routes: Routes = [{
    path: 'ctl',
    component: CtlComponent,
    loadChildren: './ctl/ctl.module#CtlModule',
}, {
    path: '',
    component: HomeComponent
}];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {

}
