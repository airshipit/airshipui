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

import { Component } from '@angular/core';
import { WsService } from 'src/services/ws/ws.service';
import { WsMessage, WsReceiver, WsConstants } from 'src/services/ws/ws.models';
import { Log } from 'src/services/log/log.service';
import { LogMessage } from 'src/services/log/log-message';
import { KustomNode, RunOptions } from './phase.models';
import { NestedTreeControl } from '@angular/cdk/tree';
import { MatTreeNestedDataSource } from '@angular/material/tree';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { PhaseViewerComponent } from './phase-viewer/phase-viewer.component';
import { PhaseRunnerComponent } from './phase-runner/phase-runner.component';

@Component({
  selector: 'app-phase',
  templateUrl: './phase.component.html',
  styleUrls: ['./phase.component.css']
})

export class PhaseComponent implements WsReceiver {
  className = this.constructor.name;
  statusMsg: string;
  loading: boolean;
  running: boolean;
  isOpen: boolean;
  phaseViewerRef: MatDialogRef<PhaseViewerComponent, any>;
  phaseRunnerRef: MatDialogRef<PhaseRunnerComponent, any>;

  type = WsConstants.CTL;
  component = WsConstants.PHASE;
  activeLink = 'overview';

  phaseTree: KustomNode[] = [];

  treeControl = new NestedTreeControl<KustomNode>(node => node.children);
  dataSource = new MatTreeNestedDataSource<KustomNode>();

  currentDocId: string;

  showEditor: boolean;
  saveBtnDisabled = true;
  editorOptions = {language: 'yaml', automaticLayout: true, value: '', theme: 'airshipTheme'};
  code: string;
  editorTitle: string;

  hasChild = (_: number, node: KustomNode) => !!node.children && node.children.length > 0;

  onInit(editor): void {
    editor.onDidChangeModelContent(() => {
      this.saveBtnDisabled = false;
    });
  }

  constructor(private websocketService: WsService, public dialog: MatDialog) {
    this.websocketService.registerFunctions(this);
    this.getPhaseTree(); // load the source first
  }

  public async receiver(message: WsMessage): Promise<void> {
    if (message.hasOwnProperty(WsConstants.ERROR)) {
      this.websocketService.printIfToast(message);
      this.loading = false;
      if (message.subComponent === WsConstants.RUN) {
        this.toggleNode(message.id);
      }
    } else {
      switch (message.subComponent) {
        case WsConstants.GET_PHASE_TREE:
          this.handleGetPhaseTree(message.data);
          break;
        case WsConstants.GET_PHASE:
          this.handleGetPhase(message);
          break;
        case WsConstants.GET_YAML:
          this.handleGetYaml(message);
          break;
        case WsConstants.GET_DOCUMENT_BY_SELECTOR:
          this.handleGetDocumentsBySelector(message);
          break;
        case WsConstants.GET_EXECUTOR_DOC:
          this.handleGetExecutorDoc(message);
          break;
        case WsConstants.YAML_WRITE:
          this.handleYamlWrite(message);
          break;
        case WsConstants.VALIDATE_PHASE:
          this.handleValidatePhase(message);
          break;
        case WsConstants.RUN:
          this.handleRunPhase(message);
          break;
        default:
          Log.Error(new LogMessage('Phase message sub component not handled', this.className, message));
          break;
      }
    }
  }

  handleValidatePhase(message: WsMessage): void {
    this.websocketService.printIfToast(message);
  }

  handleRunPhase(message: WsMessage): void {
    this.toggleNode(message.id);
  }

  handleGetPhaseTree(data: JSON): void {
    this.loading = false;
    Object.assign(this.phaseTree, data);
    this.dataSource.data = this.phaseTree;
  }

  handleGetPhase(message: WsMessage): void {
    this.loading = false;
    let yaml = '';
    if (message.yaml !== '' && message.yaml !== undefined) {
      yaml = atob(message.yaml);
    }
    this.phaseViewerRef = this.openPhaseDialog(message.id, message.name, message.details, yaml);
  }

  handleGetExecutorDoc(message: WsMessage): void {
    this.phaseViewerRef.componentInstance.executorYaml = atob(message.yaml);
  }

  handleGetDocumentsBySelector(message: WsMessage): void {
    this.phaseViewerRef.componentInstance.loading = false;
    Object.assign(this.phaseViewerRef.componentInstance.results, message.data);
  }

  handleGetYaml(message: WsMessage): void {
    if (message.message === 'rendered') {
      this.phaseViewerRef.componentInstance.yaml = atob(message.yaml);
    } else {
      this.changeEditorContents((message.yaml));
      this.editorTitle = message.name;
      this.showEditor = true;
      this.currentDocId = message.id;
    }
  }

  handleYamlWrite(message: WsMessage): void {
    this.changeEditorContents((message.yaml));
    this.editorTitle = message.name;
    this.currentDocId = message.id;
    this.websocketService.printIfToast(message);
  }

  changeEditorContents(yaml: string): void {
    this.code = atob(yaml);
  }

  saveYaml(): void {
    const websocketMessage = this.newMessage(WsConstants.YAML_WRITE);
    websocketMessage.id = this.currentDocId;
    websocketMessage.name = this.editorTitle;
    websocketMessage.yaml = btoa(this.code);
    this.websocketService.sendMessage(websocketMessage);
  }

  getPhaseTree(): void {
    this.loading = true;
    const websocketMessage = this.newMessage(WsConstants.GET_PHASE_TREE);
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
        id: node.phaseId,
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

  getPhase(id: object): void {
    this.loading = true;
    const msg = this.newMessage(WsConstants.GET_PHASE);
    msg.id = JSON.stringify(id);
    this.websocketService.sendMessage(msg);
  }

  getYaml(id: string): void {
    this.code = null;
    const msg = this.newMessage(WsConstants.GET_YAML);
    msg.id = id;
    msg.message = 'source';
    this.websocketService.sendMessage(msg);
  }

  getExecutorDoc(id: object): void {
    const msg = this.newMessage(WsConstants.GET_EXECUTOR_DOC);
    msg.id = JSON.stringify(id);
    this.websocketService.sendMessage(msg);
  }

  closeEditor(): void {
    this.code = null;
    this.showEditor = false;
  }

  // TODO(mfuller): we'll probably want to run / check phase validation
  // before actually running the phase
  runPhase(node: KustomNode, opts: RunOptions): void {
    node.running = true;
    const msg = this.newMessage(WsConstants.RUN);
    msg.id = JSON.stringify(node.phaseId);
    if (opts !== undefined) {
      msg.data = JSON.parse(JSON.stringify(opts));
    }
    this.websocketService.sendMessage(msg);
  }

  validatePhase(id: object): void {
    const msg = this.newMessage(WsConstants.VALIDATE_PHASE);
    msg.id = JSON.stringify(id);
    this.websocketService.sendMessage(msg);
  }

  newMessage(subComponent: string): WsMessage {
    return new WsMessage(this.type, this.component, subComponent);
  }

  findNode(node: KustomNode, id: string): KustomNode {
    if (node.id === id) {
      return node;
    }

    for (const child of node.children) {
      const c = this.findNode(child, id);
      if (c) {
        return c;
      }
    }
  }

  toggleNode(id: string): void {
    const phaseID = JSON.parse(id);
    for (const node of this.phaseTree) {
      const name = phaseID.Name as string;
      if (node.phaseId.Name === name) {
        node.running = false;
        return;
      }
    }
  }
}
