import {NgModule} from '@angular/core';
import {RouterModule} from '@angular/router';
import {CtlComponent} from './ctl.component';
import {DocumentModule} from './document/document.module';
import {BaremetalModule} from './baremetal/baremetal.module';
import {CtlRoutingModule} from './ctl-routing.module';

@NgModule({
  imports: [
    CtlRoutingModule,

    RouterModule,
    DocumentModule,
    BaremetalModule
  ],
  declarations: [ CtlComponent ],
  providers: []
})
export class CtlModule {}
