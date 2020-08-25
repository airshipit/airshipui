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
