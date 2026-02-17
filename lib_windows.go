//go:build windows

package kjarni

import "syscall"

func openLibrary(path string) (uintptr, error) {
	h, err := syscall.LoadLibrary(path)
	if err != nil {
		return 0, err
	}
	return uintptr(h), nil
}

func findSymbol(lib uintptr, name string) (uintptr, error) {
	proc, err := syscall.GetProcAddress(syscall.Handle(lib), name)
	if err != nil {
		return 0, err
	}
	return uintptr(proc), nil
}