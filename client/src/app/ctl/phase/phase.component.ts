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
import {KustomNode, RunOptions} from './phase.models';
import {NestedTreeControl} from '@angular/cdk/tree';
import {MatTreeNestedDataSource} from '@angular/material/tree';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { PhaseViewerComponent } from './phase-viewer/phase-viewer.component';
import {PhaseRunnerComponent} from './phase-runner/phase-runner.component';

@Component({
  selector: 'app-phase',
  templateUrl: './phase.component.html',
  styleUrls: ['./phase.component.css']
})

export class PhaseComponent implements WSReceiver {
  className = this.constructor.name;
  statusMsg: string;
  loading: boolean;
  running: boolean;
  isOpen: boolean;
  phaseViewerRef: MatDialogRef<PhaseViewerComponent, any>;
  phaseRunnerRef: MatDialogRef<PhaseRunnerComponent, any>;

  type = 'ctl';
  component = 'phase';
  activeLink = 'overview';

  targetPath: string;
  phaseTree: KustomNode[] = [];

  treeControl = new NestedTreeControl<KustomNode>(node => node.children);
  dataSource = new MatTreeNestedDataSource<KustomNode>();

  currentDocId: string;

  showEditor: boolean;
  saveBtnDisabled = true;
  editorOptions = {language: 'yaml', automaticLayout: true, value: '', theme: 'airshipTheme'};
  code: string;
  editorTitle: string;
  editorSubtitle: string;

  hasChild = (_: number, node: KustomNode) => !!node.children && node.children.length > 0;

  onInit(editor): void {
    editor.onDidChangeModelContent(() => {
      this.saveBtnDisabled = false;
    });
  }

  constructor(private websocketService: WebsocketService, public dialog: MatDialog) {
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
        case 'getPhase':
          this.handleGetPhase(message);
          break;
        case 'getYaml':
          this.handleGetYaml(message);
          break;
        case 'getDocumentsBySelector':
          this.handleGetDocumentsBySelector(message);
          break;
        case 'getExecutorDoc':
          this.handleGetExecutorDoc(message);
          break;
        case 'yamlWrite':
          this.handleYamlWrite(message);
          break;
        case 'validatePhase':
          this.handleValidatePhase(message);
          break;
        case 'run':
          this.handleRunPhase(message);
          break;
        default:
          Log.Error(new LogMessage('Phase message sub component not handled', this.className, message));
          break;
      }
    }
  }

  handleValidatePhase(message: WebsocketMessage): void {
    this.websocketService.printIfToast(message);
  }

  handleRunPhase(message: WebsocketMessage): void {
    this.running = false;
    this.websocketService.printIfToast(message);
  }

  handleGetPhaseTree(data: JSON): void {
    this.loading = false;
    Object.assign(this.phaseTree, data);
    this.dataSource.data = this.phaseTree;
  }

  handleGetPhase(message: WebsocketMessage): void {
    this.loading = false;
    let yaml = '';
    if (message.yaml !== '' && message.yaml !== undefined) {
      yaml = atob(message.yaml);
    }
    this.phaseViewerRef = this.openPhaseDialog(message.id, message.name, message.details, yaml);
  }

  handleGetExecutorDoc(message: WebsocketMessage): void {
    this.phaseViewerRef.componentInstance.executorYaml = atob(message.yaml);
  }

  handleGetDocumentsBySelector(message: WebsocketMessage): void {
    this.phaseViewerRef.componentInstance.loading = false;
    Object.assign(this.phaseViewerRef.componentInstance.results, message.data);
  }

  handleGetYaml(message: WebsocketMessage): void {
    if (message.message === 'rendered') {
      this.phaseViewerRef.componentInstance.yaml = atob(message.yaml);
    } else {
      this.changeEditorContents((message.yaml));
      this.setTitle(message.name);
      this.showEditor = true;
      this.currentDocId = message.id;
    }
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
    const websocketMessage = this.newMessage('yamlWrite');
    websocketMessage.id = this.currentDocId;
    websocketMessage.name = this.editorTitle;
    websocketMessage.yaml = btoa(this.code);
    this.websocketService.sendMessage(websocketMessage);
  }

  getPhaseTree(): void {
    this.loading = true;
    const websocketMessage = this.newMessage('getPhaseTree');
    this.websocketService.sendMessage(websocketMessage);
  }

  openPhaseDialog(id: string, name: string, details: string, yaml: string): MatDialogRef<PhaseViewerComponent, any> {
    const dialogRef = this.dialog.open(PhaseViewerComponent, {
      width: '80vw',
      height: '90vh',
    });

    dialogRef.componentInstance.id = id;
    dialogRef.componentInstance.name = name;
    dialogRef.componentInstance.yaml = yaml;

    if (details === '' || details === undefined) {
      details = '(Phase details not provided)';
    }

    dialogRef.componentInstance.phaseDetails = details;

    this.getExecutorDoc(JSON.parse(id));
    return dialogRef;
  }

  confirmRunPhase(node: KustomNode): void {
    const dialogRef = this.dialog.open(PhaseRunnerComponent, {
      width: '25vw',
      height: '30vh',
      data: {
        id: node.phaseid,
        name: node.name,
        options: new RunOptions()
      }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result !== undefined) {
        const runOpts: RunOptions = result.options;
        this.runPhase(node, runOpts);
      }
    });

  }

  getPhaseDetails(id: object): void {
    const msg = this.newMessage('getPhaseDetails');
    msg.id = JSON.stringify(id);
    this.websocketService.sendMessage(msg);
  }

  getPhase(id: object): void {
    this.loading = true;
    const msg = this.newMessage('getPhase');
    msg.id = JSON.stringify(id);
    this.websocketService.sendMessage(msg);
  }

  getYaml(id: string): void {
    this.code = null;
    const msg = this.newMessage('getYaml');
    msg.id = id;
    msg.message = 'source';
    this.websocketService.sendMessage(msg);
  }

  getExecutorDoc(id: object): void {
    const msg = this.newMessage('getExecutorDoc');
    msg.id = JSON.stringify(id);
    this.websocketService.sendMessage(msg);
  }

  closeEditor(): void {
    this.code = null;
    this.showEditor = false;
  }

  getTarget(): void {
    const websocketMessage = this.newMessage('getTarget');
    this.websocketService.sendMessage(websocketMessage);
  }

  // TODO(mfuller): we'll probably want to run / check phase validation
  // before actually running the phase
  runPhase(node: KustomNode, opts: RunOptions): void {
    this.running = true;
    const msg = this.newMessage('run');
    msg.id = JSON.stringify(node.phaseid);
    if (opts !== undefined) {
      msg.data = JSON.parse(JSON.stringify(opts));
    }
    this.websocketService.sendMessage(msg);
  }

  validatePhase(id: object): void {
    const msg = this.newMessage('validatePhase');
    msg.id = JSON.stringify(id);
    this.websocketService.sendMessage(msg);
  }

  newMessage(subComponent: string): WebsocketMessage {
    return new WebsocketMessage(this.type, this.component, subComponent);
  }
}
