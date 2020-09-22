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

import { Injectable } from '@angular/core';
import { LogLevel } from 'src/services/log/log.enum';
import { LogMessage } from 'src/services/log/log-message';

@Injectable({
  providedIn: 'root'
})
export class Log {
  static Level: LogLevel = LogLevel.Trace;
  LogWithDate = true;


  static Debug(message: LogMessage): void {
    this.writeToLog(LogLevel.Debug, message);
  }

  static Info(message: LogMessage): void {
    this.writeToLog(LogLevel.Info, message);
  }

  static Warn(message: LogMessage): void {
    this.writeToLog(LogLevel.Warn, message);
  }

  static Error(message: LogMessage): void {
    this.writeToLog(LogLevel.Error, message);
  }

  static Fatal(message: LogMessage): void {
    this.writeToLog(LogLevel.Fatal, message);
  }

  private static  writeToLog(level: LogLevel, message: LogMessage): void {
    if (level <= this.Level) {
      console.log(
        '[airshipui][' + LogLevel[level] + '] ' + new Date().toLocaleString() + ' - ' +
          message.className + ' - ' + message.message + ': ', message.logMessage);
    }
  }
}
