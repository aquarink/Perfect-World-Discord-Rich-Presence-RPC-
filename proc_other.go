//go:build !windows

package main

import "fmt"

func showAlert(title, message string) {
	fmt.Printf("[%s] %s\n", title, message)
}

func acquireSingleInstanceMutex(name string) (uintptr, bool) {
	return 1, true
}

func releaseSingleInstanceMutex(handle uintptr) {
}

func countRunningProcesses(targetNames ...string) int {
	return 0
}
