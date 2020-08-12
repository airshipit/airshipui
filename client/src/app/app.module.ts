import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

import { MatToolbarModule } from '@angular/material/toolbar';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatTableModule } from '@angular/material/table';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatExpansionModule } from '@angular/material/expansion';
import { RouterModule } from '@angular/router';
import { CTLComponent } from './ctl/ctl.component';
import { DashboardsComponent } from './dashboards/dashboards.component';
import { HomeComponent } from './home/home.component';
import { BareMetalComponent } from './ctl/baremetal/baremetal.component';
import { DocumentComponent } from './ctl/document/document.component';
import { HttpClientModule } from '@angular/common/http';
import { FlexLayoutModule } from '@angular/flex-layout';
import { MatTabsModule } from '@angular/material/tabs';
import { WebsocketService } from '../services/websocket/websocket.service';
import { ToastrModule } from 'ngx-toastr';
import {FormsModule} from '@angular/forms';
import {MonacoEditorModule, NgxMonacoEditorConfig } from 'ngx-monaco-editor';
import {MatTreeModule} from '@angular/material/tree';
import {MatButtonToggleModule} from '@angular/material/button-toggle';

@NgModule({
  declarations: [
    AppComponent,
    CTLComponent,
    DashboardsComponent,
    HomeComponent,
    BareMetalComponent,
    DocumentComponent,
  ],
  imports: [
    AppRoutingModule,
    BrowserModule,
    BrowserAnimationsModule,
    FlexLayoutModule,
    FormsModule,
    MatToolbarModule,
    MatSidenavModule,
    MatListModule,
    MatIconModule,
    MatButtonModule,
    MatTableModule,
    MatCheckboxModule,
    MatExpansionModule,
    HttpClientModule,
    RouterModule,
    MatTabsModule,
    ToastrModule.forRoot(),
    MonacoEditorModule.forRoot(),
    MatTreeModule,
    MatButtonToggleModule,
  ],
  providers: [WebsocketService],
  bootstrap: [AppComponent]
})
export class AppModule { }
