package hideconsole

import (
	"syscall"
)

func HideConsole1() {
	getWin := syscall.NewLazyDLL("kernel32.dll").NewProc("GetConsoleWindow")
	showWin := syscall.NewLazyDLL("user32.dll").NewProc("ShowWindow")
	hwnd, _, _ := getWin.Call()
	_, _, _ = showWin.Call(hwnd, 0)
}

func HideConsole2() {
	getWin := syscall.NewLazyDLL("user32.dll").NewProc("GetForegroundWindow")
	showWin := syscall.NewLazyDLL("user32.dll").NewProc("ShowWindow")
	hwnd, _, _ := getWin.Call()
	_, _, _ = showWin.Call(hwnd, 0)
}

func Run() {
	HideConsole1()
	HideConsole2()
}
