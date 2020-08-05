import { Component, OnInit } from '@angular/core';
import {WebsocketService} from '../../../services/websocket/websocket.service';
import { WSReceiver } from '../../../services/websocket/websocket.models';
import {WebsocketMessage} from '../../../services/websocket/models/websocket-message/websocket-message';
import {KustomNode} from './kustom-node';
import {NestedTreeControl} from '@angular/cdk/tree';
import {MatTreeNestedDataSource} from '@angular/material/tree';

@Component({
  selector: 'app-document',
  templateUrl: './document.component.html',
  styleUrls: ['./document.component.css']
})

export class DocumentComponent implements WSReceiver {
  obby: string;

  type: string = 'ctl';
  component: string = 'document';

  activeLink = 'overview';

  obj: KustomNode[] = [];
  currentDocId: string;

  saveBtnDisabled: boolean = true;
  hideButtons: boolean = true;
  isRendered: boolean = false;

  editorOptions = {language: 'yaml', automaticLayout: true, value: ''};
  code: string;
  editorTitle: string;
  onInit(editor) {
    editor.onDidChangeModelContent(() => {
      this.saveBtnDisabled = false;
  });
  }

  treeControl = new NestedTreeControl<KustomNode>(node => node.children);
  dataSource = new MatTreeNestedDataSource<KustomNode>();

  constructor(private websocketService: WebsocketService) {
    this.websocketService.registerFunctions(this);
    this.getSource(); // load the source first
  }

  hasChild = (_: number, node: KustomNode) => !!node.children && node.children.length > 0;

  public async receiver(message: WebsocketMessage): Promise<void> {
    if (message.hasOwnProperty("error")) {
      this.websocketService.printIfToast(message);
    } else {
      switch (message.subComponent) {
        case 'getDefaults':
          Object.assign(this.obj, message.data);
          this.dataSource.data = this.obj;
          break;
        case 'getSource':
          this.closeEditor();
          Object.assign(this.obj, message.data);
          this.dataSource.data = this.obj;
          break;
        case 'getRendered':
          this.closeEditor();
          Object.assign(this.obj, message.data);
          this.dataSource.data = this.obj;
          break;
        case 'getYaml':
          this.changeEditorContents((message.yaml));
          this.editorTitle = message.name;
          this.currentDocId = message.message;
          if (!this.isRendered) {
            this.hideButtons = false;
          } else {
            this.hideButtons = true;
          }
          break;
        case 'yamlWrite':
          this.changeEditorContents((message.yaml));
          this.editorTitle = message.name;
          this.currentDocId = message.message;
          break;
        case 'docPull':
          this.obby = "Message pull was a " + message.message;
          break;
        default:
          console.log("Document message sub component not handled: ", message);
          break;
      }
    }
  }

  getYaml(id: string): void {
    this.code = null;
    const websocketMessage = this.constructDocumentWsMessage("getYaml");
    websocketMessage.message = id;
    this.websocketService.sendMessage(websocketMessage);
  }

  changeEditorContents(yaml: string): void {
    this.code = atob(yaml);
  }

  saveYaml(): void {
    const websocketMessage = this.constructDocumentWsMessage("yamlWrite");
    websocketMessage.message = this.currentDocId;
    websocketMessage.name = this.editorTitle;
    websocketMessage.yaml = btoa(this.code);
    this.websocketService.sendMessage(websocketMessage);
  }

  getSource(): void {
    this.isRendered = false;
    const websocketMessage = this.constructDocumentWsMessage("getSource");
    this.websocketService.sendMessage(websocketMessage);
  }

  getRendered(): void {
    this.isRendered = true;
    const websocketMessage = this.constructDocumentWsMessage("getRendered");
    this.websocketService.sendMessage(websocketMessage);
  }

  constructDocumentWsMessage(subComponent: string): WebsocketMessage {
    return new WebsocketMessage(this.type, this.component, subComponent);
  }

  closeEditor(): void {
    this.code = null;
    this.editorTitle = "";
    this.hideButtons = true;
  }

  documentPull(): void {
    this.websocketService.sendMessage(new WebsocketMessage(this.type, this.component, "docPull"));
  }
}

