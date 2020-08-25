import { NgModule } from '@angular/core';
import { LoginComponent } from './login.component';
import {ToastrModule} from 'ngx-toastr';

@NgModule({
    imports: [
      ToastrModule
    ],
    declarations: [
      LoginComponent,
    ]
})

export class LoginModule { }
