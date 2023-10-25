package fakeerror

import (
	"syscall"
	"unsafe"
)

func Run() {
	var title, text *uint16
	title, _ = syscall.UTF16PtrFromString("Fatal Error")
	text, _ = syscall.UTF16PtrFromString("Error code: Windows_0x988958\nSomething gone wrong.")
	syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(0, uintptr(unsafe.Pointer(text)), uintptr(unsafe.Pointer(title)), 0)
}
