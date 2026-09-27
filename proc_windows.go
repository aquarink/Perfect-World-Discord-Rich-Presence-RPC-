//go:build windows

package main

import (
	"strings"
	"syscall"
	"unsafe"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	user32   = syscall.NewLazyDLL("user32.dll")

	procMessageBoxW              = user32.NewProc("MessageBoxW")
	procCreateMutexW             = kernel32.NewProc("CreateMutexW")
	procCloseHandle              = kernel32.NewProc("CloseHandle")
	procCreateToolhelp32Snapshot = kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32FirstW          = kernel32.NewProc("Process32FirstW")
	procProcess32NextW           = kernel32.NewProc("Process32NextW")
)

const (
	TH32CS_SNAPPROCESS   = 0x00000002
	ERROR_ALREADY_EXISTS = 183
)

type PROCESSENTRY32W struct {
	DwSize              uint32
	CntUsage            uint32
	Th32ProcessID       uint32
	Th32DefaultHeapID   uintptr
	Th32ModuleID        uint32
	CntThreads          uint32
	Th32ParentProcessID uint32
	PcPriClassBase      int32
	DwFlags             uint32
	SzExeFile           [syscall.MAX_PATH]uint16
}

func showAlert(title, message string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	msgPtr, _ := syscall.UTF16PtrFromString(message)
	procMessageBoxW.Call(0, uintptr(unsafe.Pointer(msgPtr)), uintptr(unsafe.Pointer(titlePtr)), 0x30)
}

func acquireSingleInstanceMutex(name string) (uintptr, bool) {
	mutexNamePtr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0, true
	}
	handle, _, callErr := procCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(mutexNamePtr)))
	if handle == 0 {
		return 0, true
	}
	if errno, ok := callErr.(syscall.Errno); ok && errno == ERROR_ALREADY_EXISTS {
		procCloseHandle.Call(handle)
		return 0, false
	}
	return handle, true
}

func releaseSingleInstanceMutex(handle uintptr) {
	if handle != 0 {
		procCloseHandle.Call(handle)
	}
}

func countRunningProcesses(targetNames ...string) int {
	hSnapshot, _, _ := procCreateToolhelp32Snapshot.Call(TH32CS_SNAPPROCESS, 0)
	if hSnapshot == uintptr(syscall.InvalidHandle) || hSnapshot == 0 {
		return 0
	}
	defer procCloseHandle.Call(hSnapshot)

	var entry PROCESSENTRY32W
	entry.DwSize = uint32(unsafe.Sizeof(entry))

	ret, _, _ := procProcess32FirstW.Call(hSnapshot, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		return 0
	}

	targets := make([]string, len(targetNames))
	for i, t := range targetNames {
		targets[i] = strings.ToLower(t)
	}

	count := 0
	for {
		exe := strings.ToLower(syscall.UTF16ToString(entry.SzExeFile[:]))
		for _, t := range targets {
			if exe == t {
				count++
				break
			}
		}

		ret, _, _ = procProcess32NextW.Call(hSnapshot, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
	}
	return count
}
