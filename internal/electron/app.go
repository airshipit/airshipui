/*
Copyright (c) 2020 AT&T. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/
package electron

import (
	"log"
	"os/exec"
	"syscall"
)

// RunElectron executes the standalone electron app which serves up our web components
func RunElectron() {
	cmd := exec.Command("npm", "start", "--prefix", "web")

	// make sure the logs are coming out of electron app
	//cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Fatalf("Exit with error: %v", exitError)
		}
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			log.Fatalf("Electron app close detected, exiting")
		}
	}
}
