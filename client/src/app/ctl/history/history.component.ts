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
import { WsService } from 'src/services/ws/ws.service';
import { WsMessage, WsReceiver, WsConstants } from 'src/services/ws/ws.models';
import { Log } from 'src/services/log/log.service';
import { LogMessage } from 'src/services/log/log-message';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { StatData } from './history.models';

@Component({
  selector: 'app-bare-metal',
  templateUrl: './history.component.html',
  styleUrls: ['./history.component.css']
})

export class HistoryComponent implements WsReceiver, OnInit {
  className = this.constructor.name;
  type = WsConstants.CTL;
  component = WsConstants.HISTORY;

  selectedHistory = 'baremetal';

  displayedColumns: string[] = ['subComponent', 'user', 'type', 'target', 'success', 'started', 'elapsed', 'stopped'];
  dataSources: Map<string, MatTableDataSource<StatData>> = new Map();
  dataSource: MatTableDataSource<StatData> = new MatTableDataSource();
  @ViewChild(MatSort) sort: MatSort;
  @ViewChild(MatPaginator) paginator: MatPaginator;

  constructor(private websocketService: WsService) {
    this.websocketService.registerFunctions(this);
  }

  async receiver(message: WsMessage): Promise<void> {
    if (message.hasOwnProperty(WsConstants.ERROR)) {
      this.websocketService.printIfToast(message);
    } else {
      switch (message.subComponent) {
        case WsConstants.GET_DEFAULTS:
          this.pushData(message.data);
          break;
        default:
          Log.Error(new LogMessage('History message sub component not handled', this.className, message));
          break;
      }
    }
  }

  ngOnInit(): void {
    this.refresh();
  }

  // Filters the table based on the user input
  // taken partly from the example: https://material.angular.io/components/table/overview
  applyFilter(event: Event): void {
    // get the filter value
    const filterValue = (event.target as HTMLInputElement).value;
    this.dataSource.filter = filterValue.trim().toLowerCase();

    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }

  // hide / show tables based on what's selected
  displayChange(displaying): void {
    // tag the replace string with the current context
    this.selectedHistory = displaying;

    // clear filters if they exist
    const filter = (document.getElementById('filterInput') as HTMLInputElement);
    if (filter.value.length > 0) {
      filter.value = '';
      this.dataSource.filter = '';
    }

    // if we have data show the table otherwise show nothing.  Do you hear me Lebowski?  NOTHING!
    if (displaying in this.dataSources) {
      this.dataSource = this.dataSources[displaying];
      this.dataSource.paginator = this.paginator;
      this.dataSource.sort = this.sort;

      document.getElementById('HistoryFound').removeAttribute('hidden');
      document.getElementById('LoadingHistory').setAttribute('hidden', 'true');
      document.getElementById('HistoryNotFound').setAttribute('hidden', 'true');
    } else {
      document.getElementById('HistoryFound').setAttribute('hidden', 'true');
      document.getElementById('LoadingHistory').setAttribute('hidden', 'true');
      document.getElementById('HistoryNotFound').removeAttribute('hidden');
    }
  }

  refresh(): void {
    document.getElementById('HistoryFound').setAttribute('hidden', 'true');
    document.getElementById('LoadingHistory').removeAttribute('hidden');
    document.getElementById('HistoryNotFound').setAttribute('hidden', 'true');

    const message = new WsMessage(this.type, this.component, WsConstants.GET_DEFAULTS);
    Log.Debug(new LogMessage('Attempting to ask for node data', this.className, message));
    this.websocketService.sendMessage(message);
  }

  // extract the data structure sent from the backend & render it to the table
  private pushData(data): void {
    Object.keys(data).forEach(key => {
      const recordSet: StatData[] = new Array();
      data[key].forEach(record => {
        // we do it this way instead of a straight convert because we want to format dates and success / fail messages
        // it could be argued that this should be done on the backend
        recordSet.push({
          subComponent: record.SubComponent,
          user: record.User,
          type: record.ActionType,
          target: record.Target,
          success: record.Success ? 'Succeeded' : 'Failed',
          started: new Date(record.Started).toString(),
          elapsed: record.Elapsed + (record.Elapsed < 1000 ? ' ms' : ' seconds'),
          stopped: new Date(record.Stopped).toString()
        });
      });

      this.dataSources[key] = new MatTableDataSource(recordSet);
    });

    const displaying = (document.getElementById('displaySelect') as HTMLInputElement).value;
    this.displayChange(displaying);
  }
}
