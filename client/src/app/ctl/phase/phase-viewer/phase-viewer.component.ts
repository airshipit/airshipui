/*
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/

import { Component, OnInit } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { KustomNode } from '../phase.models';
import { WsMessage, WsConstants } from 'src/services/ws/ws.models';
import { WsService } from 'src/services/ws/ws.service';
import { FormControl, FormGroup } from '@angular/forms';

@Component({
    selector: 'app-phase-viewer',
    templateUrl: 'phase-viewer.component.html',
    styleUrls: ['./phase-viewer.component.css']

})

export class PhaseViewerComponent implements OnInit {
    editorOptions = {language: 'yaml', automaticLayout: true, readOnly: true, theme: 'airshipTheme'};
    bundleYaml: string;
    executorYaml: string;
    phaseDetails: string;
    loading: boolean;
    resultsMsg = '';

    results: KustomNode[] = [];
    id: string;
    name: string;
    yaml: string;

    filterOptions = new FormGroup({
        name: new FormControl(''),
        namespace: new FormControl(''),
        gvk: new FormControl(''),
        kind: new FormControl(''),
        label: new FormControl(''),
        annotation: new FormControl('')
    });

    constructor(
        public dialogRef: MatDialogRef<PhaseViewerComponent>,
        private websocketService: WsService) {}

    ngOnInit(): void {
        this.bundleYaml = this.yaml;
        if (this.bundleYaml !== '') {
            this.getDocumentsBySelector('{}');
        }
        this.yaml = this.phaseDetails;
    }

    onClose(): void {
        this.dialogRef.close();
        this.results = null;
    }

    setModel(val: string): void {
        switch (val) {
            case WsConstants.BUNDLE:
                this.yaml = this.bundleYaml;
                break;
            case WsConstants.EXECUTOR:
                this.yaml = this.executorYaml;
                break;
            case WsConstants.DETAILS:
                this.yaml = this.phaseDetails;
                break;
        }
    }

    getDocumentsBySelector(selector: string): void {
        const msg = new WsMessage(WsConstants.CTL, WsConstants.PHASE, WsConstants.GET_DOCUMENT_BY_SELECTOR);
        msg.message = selector;
        msg.id = this.id;
        this.websocketService.sendMessage(msg);
    }

    getYaml(id: string): void {
        this.yaml = null;
        const msg = new WsMessage(WsConstants.CTL, WsConstants.PHASE, WsConstants.GET_YAML);
        msg.id = id;
        msg.message = 'rendered';
        this.websocketService.sendMessage(msg);
      }

    onSubmit(data: any): void {
        this.loading = true;
        this.results = [];
        this.resultsMsg = '';
        const selector = {};
        Object.keys(this.filterOptions.controls).forEach(key => {
            if (this.filterOptions.controls[key].value !== '') {
                if (key === 'gvk') {
                    const str: string = this.filterOptions.controls[key].value;
                    const arr = str.split(' ');
                    selector[key] = {
                        group: arr[0],
                        version: arr[1],
                        kind: arr[2]
                    };
                } else {
                    selector[key] = this.filterOptions.controls[key].value;
                }
            }
        });
        this.getDocumentsBySelector(JSON.stringify(selector));
    }
}
