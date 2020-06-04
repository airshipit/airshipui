// +build !windows

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
	"syscall"

	"opendev.org/airship/airshipui/internal/webservice"
)

// ProcessGrpCmd wraps an exec.Cmd and a signal chan
type ProcessGrpCmd struct {
	sigChan chan os.Signal
	ctx     context.Context
	waitgrp *sync.WaitGroup
	*exec.Cmd
}

// NewProcessGrpCmd creates a new ProcessGrpCmd to monitor a os.Signal chan
func NewProcessGrpCmd(c context.Context, wg *sync.WaitGroup, s chan os.Signal, cmd *exec.Cmd) *ProcessGrpCmd {
	return &ProcessGrpCmd{
		ctx:     c,
		waitgrp: wg,
		sigChan: s,
		Cmd:     cmd,
	}
}

// Run sets the process group id attribute to ensure all child
// processes started by the cmd will inherit it and will get killed
// with the parent
func (grp *ProcessGrpCmd) Run() error {
	grp.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := grp.Cmd.Start(); err != nil {
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

// RunBinaryWithOptions runs the binary
func RunBinaryWithOptions(ctx context.Context, cmd string, args []string, wg *sync.WaitGroup, s chan os.Signal) {
	command := exec.CommandContext(ctx, cmd, args...)

	// push executable's stdout / stderr
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	monitoredCmd := NewProcessGrpCmd(ctx, wg, s, command)

	if err := monitoredCmd.Run(); err != nil {
		log.Printf("'%s' exited with error: %v", cmd, err)

		// send error to UI
		webservice.SendAlert(
			webservice.Error,
			fmt.Sprintf("Plugin '%s' failed to start: %v", cmd, err),
		)
	}
}

// Terminate kills the running process, logs a message, and sends Done() to the WaitGroup
func (grp *ProcessGrpCmd) Terminate() {
	log.Printf("Terminating process '%s' with PID: %d", grp.Cmd.Path, grp.Cmd.Process.Pid)
	if err := syscall.Kill(-grp.Cmd.Process.Pid, syscall.SIGKILL); err != nil {
		log.Printf("Error terminating PID %d: %s", grp.Cmd.Process.Pid, err)
	}
	grp.waitgrp.Done()
}
