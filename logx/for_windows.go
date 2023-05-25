//go:build windows
// +build windows

package logx

import "syscall"

func colorPrint(s string, i int) {
	p := proc.(*syscall.LazyProc)
	c := closeHandle.(*syscall.LazyProc)
	handle, _, _ := p.Call(uintptr(syscall.Stdout), uintptr(i))
	print(s, "\n")
	c.Call(handle)
}
func initKernel32() {
	kernel32 := syscall.NewLazyDLL(`kernel32.dll`)
	proc = kernel32.NewProc(`SetConsoleTextAttribute`)
	closeHandle = kernel32.NewProc(`CloseHandle`)
}
