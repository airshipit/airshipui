import {NgModule} from '@angular/core';
import {BaremetalComponent} from './baremetal.component';
import {MatButtonModule} from '@angular/material/button';

@NgModule({
  imports: [
    MatButtonModule
  ],
  declarations: [
    BaremetalComponent
  ],
  providers: []
})
export class BaremetalModule {}
