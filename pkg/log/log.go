/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var (
	// LogLevel can specify what level the system runs at
	LogLevel = 6
	levels   = map[int]string{
		6: "TRACE",
		5: "DEBUG",
		4: "INFO",
		3: "WARN",
		2: "ERROR",
		1: "FATAL",
	}
	airshipLog = log.New(os.Stderr, "[airshipui] ", log.LstdFlags|log.Llongfile)
	writeMutex sync.Mutex
)

// Init initializes settings related to logging
func Init(levelSet int, out io.Writer) {
	LogLevel = levelSet
	airshipLog.SetOutput(out)
}

// Trace is a wrapper for log.Trace
func Trace(v ...interface{}) {
	writeLog(6, v...)
}

// Tracef is a wrapper for log.Tracef
func Tracef(format string, v ...interface{}) {
	writeLog(6, fmt.Sprintf(format, v...))
}

// Debug is a wrapper for log.Debug
func Debug(v ...interface{}) {
	writeLog(5, v...)
}

// Debugf is a wrapper for log.Debugf
func Debugf(format string, v ...interface{}) {
	writeLog(5, fmt.Sprintf(format, v...))
}

// Info is a wrapper for log.Info
func Info(v ...interface{}) {
	writeLog(4, v...)
}

// Infof is a wrapper for log.Infof
func Infof(format string, v ...interface{}) {
	writeLog(4, fmt.Sprintf(format, v...))
}

// Warn is a wrapper for log.Warn
func Warn(v ...interface{}) {
	writeLog(3, v...)
}

// Warnf is a wrapper for log.Warnf
func Warnf(format string, v ...interface{}) {
	writeLog(3, fmt.Sprintf(format, v...))
}

// Error is a wrapper for log.Error
func Error(v ...interface{}) {
	writeLog(2, v...)
}

// Errorf is a wrapper for log.Errorf
func Errorf(format string, v ...interface{}) {
	writeLog(2, fmt.Sprintf(format, v...))
}

// Fatal is a wrapper for log.Fatal
func Fatal(v ...interface{}) {
	writeLog(1, v...)
	os.Exit(-1)
}

// Fatalf is a wrapper for log.Fatalf
func Fatalf(format string, v ...interface{}) {
	writeLog(1, fmt.Sprintf(format, v...))
	os.Exit(-1)
}

// Writer returns log output writer object
func Writer() io.Writer {
	return airshipLog.Writer()
}

// Logger is used by things like net/http to overwrite their standard logging
func Logger() *log.Logger {
	return airshipLog
}

func writeLog(level int, v ...interface{}) {
	// determine if we need to display the logs
	if level <= LogLevel {
		writeMutex.Lock()
		defer writeMutex.Unlock()
		// the origionall caller of this is 3 steps back, the output will display who called it
		err := airshipLog.Output(3, fmt.Sprintf("[%s] %v", levels[level], fmt.Sprint(v...)))
		if err != nil {
			airshipLog.Print(v...)
			airshipLog.Print(err)
		}
	}
}
