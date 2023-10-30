package hideconsole

import (
	"syscall"
)

func Run() {
	getWin := syscall.NewLazyDLL("kernel32.dll").NewProc("GetConsoleWindow")
	showWin := syscall.NewLazyDLL("user32.dll").NewProc("ShowWindow")
	hwnd, _, _ := getWin.Call()
	_, _, _ = showWin.Call(hwnd, 0)
}
