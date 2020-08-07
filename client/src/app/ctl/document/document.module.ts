import {NgModule, CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {MatTabsModule} from '@angular/material/tabs';
import {DocumentComponent} from './document.component';
import {MatTreeModule} from '@angular/material/tree';
import {MatButtonModule} from '@angular/material/button';
import {MatButtonToggleModule} from '@angular/material/button-toggle';
import {MatIconModule} from '@angular/material/icon';
import {MonacoEditorModule} from 'ngx-monaco-editor';
import {FormsModule} from '@angular/forms';
import {ToastrModule} from 'ngx-toastr';
import {CommonModule} from '@angular/common';
import {MatProgressBarModule} from '@angular/material/progress-bar';
import {MatCardModule} from '@angular/material/card';
import {MatTooltipModule} from '@angular/material/tooltip/';

@NgModule({
  declarations: [
    DocumentComponent,
  ],
  imports: [
    CommonModule,
    MatTabsModule,
    MatTreeModule,
    MatButtonModule,
    MatButtonToggleModule,
    MatIconModule,
    MonacoEditorModule,
    FormsModule,
    ToastrModule,
    MatProgressBarModule,
    MatCardModule,
    MatTooltipModule
  ],
  providers: [],
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class DocumentModule {}
