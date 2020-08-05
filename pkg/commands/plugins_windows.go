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

package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/webservice"
)

// struct to store PID and process handle
type process struct {
	Pid    int
	Handle uintptr
}

// ProcessGrpCmd wraps an exec.Cmd with channels for monitoring
type ProcessGrpCmd struct {
	sigChan chan os.Signal
	ctx     context.Context
	waitgrp *sync.WaitGroup
	handle  windows.Handle
	*exec.Cmd
}

// NewProcessGrpCmd creates a new ProcessGrpCmd to monitor an exec.Cmd
func NewProcessGrpCmd(c context.Context, wg *sync.WaitGroup, s chan os.Signal, cmd *exec.Cmd) *ProcessGrpCmd {
	return &ProcessGrpCmd{
		ctx:     c,
		waitgrp: wg,
		sigChan: s,
		Cmd:     cmd,
	}
}

// NewProcessGroup returns a handle to a Windows Job object that we can
// add plugin processes to, so that any child processes will get terminated
// with the parent Job
func NewProcessGroup() (windows.Handle, error) {
	handle, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return 0, err
	}

	// ensure all processes associated with Job object get killed when
	// handle is closed
	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	if _, err := windows.SetInformationJobObject(
		handle,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info))); err != nil {
		return 0, err
	}

	return handle, nil
}

// AddProcess adds an os.Process to a Windows Job object
func (grp *ProcessGrpCmd) AddProcess(p *os.Process) error {
	return windows.AssignProcessToJobObject(
		grp.handle,
		windows.Handle((*process)(unsafe.Pointer(p)).Handle))
}

// Run modifies the behavior of cmd.Run to create a new Windows
// job object, listen for interrupts or context.Done events,
// and signal to the WaitGroup when the process is terminated
func (grp *ProcessGrpCmd) Run() error {
	handle, err := NewProcessGroup()
	if err != nil {
		grp.waitgrp.Done()
		return err
	}
	grp.handle = handle

	if err := grp.Cmd.Start(); err != nil {
		grp.waitgrp.Done()
		return err
	}

	if err := grp.AddProcess(grp.Cmd.Process); err != nil {
		grp.waitgrp.Done()
		return err
	}

	go func() {
		select {
		case <-grp.sigChan:
			grp.Terminate()
			return
		case <-grp.ctx.Done():
			grp.Terminate()
			return
		}
	}()

	return grp.Cmd.Wait()
}

// Terminate kills the running process, logs a message, and sends Done() to the WaitGroup
func (grp *ProcessGrpCmd) Terminate() {
	log.Printf("Terminating process '%s' with PID: %d", grp.Cmd.Path, grp.Cmd.Process.Pid)
	if err := windows.CloseHandle(grp.handle); err != nil {
		log.Printf("Error terminating PID %d: %s", grp.Cmd.Process.Pid, err)
	}
	grp.waitgrp.Done()
}

func RunBinaryWithOptions(ctx context.Context, cmd string, args []string, wg *sync.WaitGroup, s chan os.Signal) {
	command := exec.CommandContext(ctx, cmd, args...)

	// push executable's stdout / stderr
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	prGrp := NewProcessGrpCmd(ctx, wg, s, command)

	if err := prGrp.Run(); err != nil {
		log.Printf("'%s' exited with error: %v", cmd, err)

		// send error to UI
		webservice.SendAlert(
			configs.Error,
			fmt.Sprintf("Plugin '%s' failed to start: %v", cmd, err),
		)
	}
}
