import {Component} from '@angular/core';
import {WebsocketService} from '../../../services/websocket/websocket.service';
import {WebsocketMessage, WSReceiver} from '../../../services/websocket/websocket.models';
import {NestedTreeControl} from '@angular/cdk/tree';
import {MatTreeNestedDataSource} from '@angular/material/tree';
import {KustomNode} from './document.models';

@Component({
  selector: 'app-document',
  templateUrl: './document.component.html',
  styleUrls: ['./document.component.css']
})

export class DocumentComponent implements WSReceiver {
  obby: string;

  type = 'ctl';
  component = 'document';

  activeLink = 'overview';

  obj: KustomNode[] = [];
  currentDocId: string;

  saveBtnDisabled = true;
  hideButtons = true;
  isRendered = false;

  editorOptions = {language: 'yaml', automaticLayout: true, value: ''};
  code: string;
  editorTitle: string;

  treeControl = new NestedTreeControl<KustomNode>(node => node.children);
  dataSource = new MatTreeNestedDataSource<KustomNode>();

  onInit(editor): void {
    editor.onDidChangeModelContent(() => {
      this.saveBtnDisabled = false;
    });
  }

  constructor(private websocketService: WebsocketService) {
    this.websocketService.registerFunctions(this);
    this.getSource(); // load the source first
  }

  hasChild = (_: number, node: KustomNode) => !!node.children && node.children.length > 0;

  public async receiver(message: WebsocketMessage): Promise<void> {
    if (message.hasOwnProperty('error')) {
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
          this.hideButtons = this.isRendered;
          break;
        case 'yamlWrite':
          this.changeEditorContents((message.yaml));
          this.editorTitle = message.name;
          this.currentDocId = message.message;
          break;
        case 'docPull':
          this.obby = 'Message pull was a ' + message.message;
          break;
        default:
          console.log('Document message sub component not handled: ', message);
          break;
      }
    }
  }

  getYaml(id: string): void {
    this.code = null;
    const websocketMessage = this.constructDocumentWsMessage('getYaml');
    websocketMessage.message = id;
    this.websocketService.sendMessage(websocketMessage);
  }

  changeEditorContents(yaml: string): void {
    this.code = atob(yaml);
  }

  saveYaml(): void {
    const websocketMessage = this.constructDocumentWsMessage('yamlWrite');
    websocketMessage.message = this.currentDocId;
    websocketMessage.name = this.editorTitle;
    websocketMessage.yaml = btoa(this.code);
    this.websocketService.sendMessage(websocketMessage);
  }

  getSource(): void {
    this.isRendered = false;
    const websocketMessage = this.constructDocumentWsMessage('getSource');
    this.websocketService.sendMessage(websocketMessage);
  }

  getRendered(): void {
    this.isRendered = true;
    const websocketMessage = this.constructDocumentWsMessage('getRendered');
    this.websocketService.sendMessage(websocketMessage);
  }

  constructDocumentWsMessage(subComponent: string): WebsocketMessage {
    return new WebsocketMessage(this.type, this.component, subComponent);
  }

  closeEditor(): void {
    this.code = null;
    this.editorTitle = '';
    this.hideButtons = true;
  }

  documentPull(): void {
    this.websocketService.sendMessage(new WebsocketMessage(this.type, this.component, 'docPull'));
  }
}

