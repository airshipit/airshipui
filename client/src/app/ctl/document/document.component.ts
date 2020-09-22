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

import {Component} from '@angular/core';
import {WebsocketService} from '../../../services/websocket/websocket.service';
import {WebsocketMessage, WSReceiver} from '../../../services/websocket/websocket.models';
import {Log} from '../../../services/log/log.service';
import {LogMessage} from '../../../services/log/log-message';
import {KustomNode} from './document.models';
import {NestedTreeControl} from '@angular/cdk/tree';
import {MatTreeNestedDataSource} from '@angular/material/tree';

@Component({
  selector: 'app-document',
  templateUrl: './document.component.html',
  styleUrls: ['./document.component.css']
})

export class DocumentComponent implements WSReceiver {
  className = this.constructor.name;
  statusMsg: string;
  loading: boolean;

  type = 'ctl';
  component = 'document';
  activeLink = 'overview';

  targetPath: string;
  phaseTree: KustomNode[] = [];
  cache = new Map<string, KustomNode[]>();

  treeControl = new NestedTreeControl<KustomNode>(node => node.children);
  dataSource = new MatTreeNestedDataSource<KustomNode>();

  currentDocId: string;

  showEditor: boolean;
  saveBtnDisabled = true;
  hideButtons = true;
  editorOptions = {language: 'yaml', automaticLayout: true, value: '', theme: 'airshipTheme'};
  code: string;
  editorTitle: string;
  editorSubtitle: string;
  docType: string;

  hasChild = (_: number, node: KustomNode) => !!node.children && node.children.length > 0;

  onInit(editor): void {
    editor.onDidChangeModelContent(() => {
      this.saveBtnDisabled = false;
    });
  }

  constructor(private websocketService: WebsocketService) {
    this.websocketService.registerFunctions(this);
    this.getTarget();
    this.getPhaseTree(); // load the source first
  }

  public async receiver(message: WebsocketMessage): Promise<void> {
    if (message.hasOwnProperty('error')) {
      this.websocketService.printIfToast(message);
      this.loading = false;
    } else {
      switch (message.subComponent) {
        case 'getTarget':
          this.targetPath = message.message;
          break;
        case 'docPull':
          this.statusMsg = 'Message pull was a ' + message.message;
          break;
        case 'getPhaseTree':
          this.handleGetPhaseTree(message.data);
          break;
        case 'getYaml':
          this.handleGetYaml(message);
          break;
        case 'yamlWrite':
          this.handleYamlWrite(message);
          break;
        default:
          Log.Error(new LogMessage('Document message sub component not handled', this.className, message));
          break;
      }
    }
  }

  handleGetPhaseTree(data: JSON): void {
    this.loading = false;
    Object.assign(this.phaseTree, data);
    this.dataSource.data = this.phaseTree;
  }

  handleGetYaml(message: WebsocketMessage): void {
    this.changeEditorContents((message.yaml));
    this.setTitle(message.name);
    this.showEditor = true;
    this.currentDocId = message.id;
  }

  handleYamlWrite(message: WebsocketMessage): void {
    this.changeEditorContents((message.yaml));
    this.setTitle(message.name);
    this.currentDocId = message.id;
    this.websocketService.printIfToast(message);
  }

  setTitle(name: string): void {
    this.editorSubtitle = name;
    const str = name.split('/');
    this.editorTitle = str[str.length - 1];
  }

  changeEditorContents(yaml: string): void {
    this.code = atob(yaml);
  }

  saveYaml(): void {
    const websocketMessage = this.constructDocumentWsMessage('yamlWrite');
    websocketMessage.id = this.currentDocId;
    websocketMessage.name = this.editorTitle;
    websocketMessage.yaml = btoa(this.code);
    this.websocketService.sendMessage(websocketMessage);
  }

  getPhaseTree(): void {
    this.loading = true;
    const websocketMessage = this.constructDocumentWsMessage('getPhaseTree');
    this.websocketService.sendMessage(websocketMessage);
  }

  viewPhaseDocs(id: string): void {
    // show document viewer
  }

  getYaml(id: string): void {
    this.code = null;
    const msg = new WebsocketMessage('ctl', 'document', 'getYaml');
    msg.id = id;
    this.websocketService.sendMessage(msg);
  }

  closeEditor(): void {
    this.code = null;
    this.showEditor = false;
  }

  getTarget(): void {
    const websocketMessage = this.constructDocumentWsMessage('getTarget');
    this.websocketService.sendMessage(websocketMessage);
  }

  constructDocumentWsMessage(subComponent: string): WebsocketMessage {
    return new WebsocketMessage(this.type, this.component, subComponent);
  }

  documentPull(): void {
    this.websocketService.sendMessage(new WebsocketMessage(this.type, this.component, 'pull'));
  }
}
