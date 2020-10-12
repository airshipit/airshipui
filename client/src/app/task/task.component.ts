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
import { WebsocketService } from '../../services/websocket/websocket.service';
import { WSReceiver, WebsocketMessage } from '../../services/websocket/websocket.models';
import { Task, Progress } from './task.models';
import { Log } from '../../services/log/log.service';
import { LogMessage } from '../../services/log/log-message';

@Component({
  selector: 'app-task',
  templateUrl: './task.component.html',
  styleUrls: ['./task.component.css']
})
export class TaskComponent implements WSReceiver {
  className = this.constructor.name;
  type = 'ui';
  component = 'task';

  message: string;
  tasks: Task[] = [];
  isOpen = false;

  constructor(private websocketService: WebsocketService) {
    this.websocketService.registerFunctions(this);
  }

  public async receiver(message: WebsocketMessage): Promise<void> {
    if (message.hasOwnProperty('error')) {
      this.websocketService.printIfToast(message);
    } else {
      switch (message.subComponent) {
        case 'taskStart':
          this.handleTaskStart(message);
          break;
        case 'taskUpdate':
          this.handleTaskUpdate(message);
          break;
        case 'taskEnd':
          this.handleTaskEnd(message);
          break;
        default:
          Log.Error(new LogMessage('Task message sub component not handled', this.className, message));
          break;
      }
    }
  }

  handleTaskStart(message: WebsocketMessage): void {
    this.addTask(message);
    const msg = new WebsocketMessage(this.type, this.component, message.subComponent);
    msg.message = `${message.name} added to Running Tasks`;
    msg.sessionID = message.sessionID;
    this.websocketService.printIfToast(msg);
  }

  handleTaskUpdate(message: WebsocketMessage): void {
    const task = this.findTask(message.id);
    if (task !== null) {
      Object.assign(task.progress, message.data);
      if (task.progress.errors.length > 0) {
        task.running = false;
        task.progress.message = task.progress.errors.toString();
      }
    } else {
      const msg = new WebsocketMessage(this.type, this.component, message.subComponent);
      msg.sessionID = message.sessionID;
      msg.message = `Task with id ${message.id} not found`;
      this.websocketService.printIfToast(msg);
    }
  }

  handleTaskEnd(message: WebsocketMessage): void {
    const task = this.findTask(message.id);
    if (task !== null) {
      Object.assign(task.progress, message.data);
      task.running = false;
    } else {
      const msg = new WebsocketMessage(this.type, this.component, message.subComponent);
      msg.sessionID = message.sessionID;
      msg.message = `Task with id ${message.id} not found`;
      this.websocketService.printIfToast(msg);
    }
  }

  taskRemove(id: string): void {
    for (let i = 0; i < this.tasks.length; i++) {
      if (this.tasks[i].id === id) {
        this.tasks.splice(i, 1);
      }
    }
  }

  addTask(message: WebsocketMessage): void {
    const p = new Progress();
    Object.assign(p, message.data);

    const task: Task = {
      id: message.id,
      name: message.name,
      running: true,
      progress: p
    };
    this.tasks.push(task);
  }

  findTask(id: string): Task {
    for (const task of this.tasks) {
      if (task.id === id) {
        return task;
      }
    }
    return null;
  }

  // TODO(mfuller): this was intended to be used for tooltip content, but
  // I can't get the tooltip to show up on menu items, even with 'hello world'
  taskToString(task: Task): string {
    return `Name: ${task.name}
    Start Time: ${new Date(task.progress.startTime).toUTCString()}
    Last Updated: ${new Date(task.progress.lastUpdated).toUTCString()}
    End Time: ${new Date(task.progress.endTime).toUTCString()}
    Message: ${task.progress.message}`;
  }
}
