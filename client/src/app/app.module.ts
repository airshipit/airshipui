import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {RouterModule} from '@angular/router';
import {HttpClientModule} from '@angular/common/http';
import {WebsocketService} from '../services/websocket/websocket.service';
import {ToastrModule} from 'ngx-toastr';
import {MonacoEditorModule} from 'ngx-monaco-editor';
import {MatSidenavModule} from '@angular/material/sidenav';
import {MatIconModule} from '@angular/material/icon';
import {MatExpansionModule} from '@angular/material/expansion';
import {FlexLayoutModule} from '@angular/flex-layout';
import {MatListModule} from '@angular/material/list';
import {MatToolbarModule} from '@angular/material/toolbar';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MatButtonModule} from '@angular/material/button';
import {MatTabsModule} from '@angular/material/tabs';
import {CtlModule} from './ctl/ctl.module';

@NgModule({
  imports: [
    AppRoutingModule,

    CtlModule,

    BrowserModule,
    BrowserAnimationsModule,
    FlexLayoutModule,
    HttpClientModule,
    MatButtonModule,
    MatSidenavModule,
    MatIconModule,
    MatExpansionModule,
    MatListModule,
    MatToolbarModule,
    RouterModule,
    MatTabsModule,
    ToastrModule.forRoot(),
    MonacoEditorModule.forRoot(),
  ],
  declarations: [AppComponent],
  providers: [WebsocketService],
  bootstrap: [AppComponent]
})
export class AppModule { }
