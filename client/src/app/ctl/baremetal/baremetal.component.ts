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

import { Component, OnInit, ViewChild } from '@angular/core';
import { WebsocketService } from '../../../services/websocket/websocket.service';
import { WebsocketMessage, WSReceiver } from '../../../services/websocket/websocket.models';
import { Log } from '../../../services/log/log.service';
import { LogMessage } from '../../../services/log/log-message';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { SelectionModel } from '@angular/cdk/collections';
import { NodeData, PhaseData } from './baremetal.models';

@Component({
  selector: 'app-bare-metal',
  templateUrl: './baremetal.component.html',
  styleUrls: ['./baremetal.component.css']
})

export class BaremetalComponent implements WSReceiver, OnInit {
  className = this.constructor.name;
  // TODO (aschiefe): extract these strings to constants
  type = 'ctl';
  component = 'baremetal';

  nodeColumns: string[] = ['select', 'name', 'id', 'bmcAddress'];
  nodeDataSource: MatTableDataSource<NodeData> = new MatTableDataSource();
  nodeSelection = new SelectionModel<NodeData>(true, []);
  @ViewChild('nodeTableSort', { static: false }) nodeSort: MatSort;
  @ViewChild('nodePaginator', { static: false }) nodePaginator: MatPaginator;

  phaseColumns: string[] = ['select', 'name', 'generateName', 'namespace', 'clusterName'];
  phaseDataSource: MatTableDataSource<PhaseData> = new MatTableDataSource();
  phaseSelection = new SelectionModel<PhaseData>(true, []);
  @ViewChild('phaseTableSort', { static: false }) phaseSort: MatSort;
  @ViewChild('phasePaginator', { static: false }) phasePaginator: MatPaginator;

  constructor(private websocketService: WebsocketService) {
    this.websocketService.registerFunctions(this);
  }

  async receiver(message: WebsocketMessage): Promise<void> {
    if (message.hasOwnProperty('error')) {
      this.websocketService.printIfToast(message);
    } else {
      switch (message.subComponent) {
        case 'getDefaults':
          this.pushData(message.data);
          break;
        default:
          Log.Error(new LogMessage('Baremetal message sub component not handled', this.className, message));
          break;
      }
    }
  }

  ngOnInit(): void {
    const message = new WebsocketMessage(this.type, this.component, 'getDefaults');
    Log.Debug(new LogMessage('Attempting to ask for node data', this.className, message));
    this.websocketService.sendMessage(message);
  }

  // Filters the table based on the user input
  // taken partly from the example: https://material.angular.io/components/table/overview
  applyFilter(event: Event): void {
    // get the filter value
    const filterValue = (event.target as HTMLInputElement).value;
    const displaying = (document.getElementById('displaySelect') as HTMLInputElement).value;
    let datasource: MatTableDataSource<any>;
    if (displaying === 'node') {
      datasource = this.nodeDataSource;
    } else {
      datasource = this.phaseDataSource;
    }

    datasource.filter = filterValue.trim().toLowerCase();

    if (datasource.paginator) {
      datasource.paginator.firstPage();
    }
  }

  // Whether the number of selected elements matches the total number of rows
  // taken partly from the example: https://material.angular.io/components/table/overview
  isAllSelected(): any {
    const displaying = (document.getElementById('displaySelect') as HTMLInputElement).value;
    let numSelected: number;
    let numRows: number;
    if (displaying === 'node') {
      numSelected = this.nodeSelection.selected.length;
      numRows = this.nodeDataSource.data.length;
    } else {
      numSelected = this.phaseSelection.selected.length;
      numRows = this.phaseDataSource.data.length;
    }

    // enable / disable the action items
    const select = (document.getElementById('operationSelect') as HTMLInputElement);
    if (numSelected > 0) {
      select.removeAttribute('disabled');
    } else {
      select.setAttribute('disabled', 'disabled');
      select.value = 'none';
      this.operationChange('none');
    }

    return numSelected === numRows;
  }

  // Selects all rows if they are not all selected; otherwise clear selection.
  // taken partly from the example: https://material.angular.io/components/table/overview
  masterToggle(): void {
    const displaying = (document.getElementById('displaySelect') as HTMLInputElement).value;
    if (displaying === 'node') {
      this.isAllSelected() ?
        this.nodeSelection.clear() :
        this.nodeDataSource.data.forEach(row => this.nodeSelection.select(row));
    } else {
      this.isAllSelected() ?
        this.phaseSelection.clear() :
        this.phaseDataSource.data.forEach(row => this.phaseSelection.select(row));
    }
  }

  // The label for the checkbox on the passed row
  // taken partly from the example: https://material.angular.io/components/table/overview
  checkboxLabel(row?: any): string {
    if (!row) {
      return `${this.isAllSelected() ? 'select' : 'deselect'} all`;
    }
    const displaying = (document.getElementById('displaySelect') as HTMLInputElement).value;
    if (displaying === 'node') {
      return `${this.nodeSelection.isSelected(row) ? 'deselect' : 'select'} row ${row.name}`;
    } else {
      return `${this.phaseSelection.isSelected(row) ? 'deselect' : 'select'} row ${row.name}`;
    }
  }

  // hide / show tables based on what's selected
  displayChange(displaying): void {
    if (displaying === 'node') {
      document.getElementById('NodeDiv').removeAttribute('hidden');
      document.getElementById('PhaseDiv').setAttribute('hidden', 'true');
    } else {
      document.getElementById('PhaseDiv').removeAttribute('hidden');
      document.getElementById('NodeDiv').setAttribute('hidden', 'true');
    }

    // clear out the selections & filters on change
    (document.getElementById('operationSelect') as HTMLInputElement).value = 'none';
    this.nodeSelection.clear();
    this.nodeDataSource.filter = '';
    this.phaseSelection.clear();
    this.phaseDataSource.filter = '';
  }

  // control if the run button is enabled based on the select menu
  operationChange(value): void {
    const button = document.getElementById('runButton');
    value !== 'none' ? button.removeAttribute('disabled') : button.setAttribute('disabled', 'disabled');
  }

  actionRun(): void {
    // retrieve the action to be run
    const subComponent = (document.getElementById('operationSelect') as HTMLInputElement).value;

    // retrieve the targets to run the action against
    // create the websocket message & fire the request to the backend
    const message = new WebsocketMessage(this.type, this.component, subComponent);
    const displaying = (document.getElementById('displaySelect') as HTMLInputElement).value;
    const targets: string[] = new Array();
    if (displaying === 'node') {
      this.nodeSelection.selected.forEach(node => {
        targets.push(node.name);
      });
      message.actionType = 'direct';
    } else {
      this.phaseSelection.selected.forEach(phase => {
        targets.push(phase.name);
      });
      message.actionType = 'phase';
    }
    message.targets = targets;

    Log.Debug(new LogMessage('Attempting to perform action(s)', this.className, message));
    this.websocketService.sendMessage(message);
  }

  // extract the data structure sent from the backend & render it to the table
  private pushData(data): void {
    const nodeConvertible: NodeData[] = data.nodes;
    this.nodeDataSource = new MatTableDataSource(nodeConvertible);
    this.nodeDataSource.paginator = this.nodePaginator;
    this.nodeDataSource.sort = this.nodeSort;

    const phaseConvertible: PhaseData[] = data.phases;
    this.phaseDataSource = new MatTableDataSource(phaseConvertible);
    this.phaseDataSource.paginator = this.phasePaginator;
    this.phaseDataSource.sort = this.phaseSort;
  }
}
