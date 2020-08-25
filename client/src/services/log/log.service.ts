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
