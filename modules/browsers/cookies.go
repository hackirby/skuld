package browsers

import (
	"path/filepath"

	_ "modernc.org/sqlite"
)

func (c *Chromium) GetCookies(path string) (cookies []Cookie, err error) {
	db, err := GetDBConnection(filepath.Join(path, "Network", "Cookies"))
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT name, encrypted_value, host_key, path, expires_utc FROM cookies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name, host, path      string
			encryptedValue, value []byte
			expiresUtc            int64
		)
		if err = rows.Scan(&name, &encryptedValue, &host, &path, &expiresUtc); err != nil {
			continue
		}

		if name == "" || host == "" || path == "" || encryptedValue == nil {
			continue
		}

		cookie := Cookie{
			Name:       name,
			Host:       host,
			Path:       path,
			ExpireDate: expiresUtc,
		}

		value, err = c.Decrypt(encryptedValue)
		if err != nil {
			continue
		}
		cookie.Value = string(value)
		cookies = append(cookies, cookie)
	}

	return cookies, nil
}

func (g *Gecko) GetCookies(path string) (cookies []Cookie, err error) {
	db, err := GetDBConnection(filepath.Join(path, "cookies.sqlite"))
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT name, value, host, path, expiry FROM moz_cookies")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			name, host, path string
			value            []byte
			expiry           int64
		)
		if err = rows.Scan(&name, &value, &host, &path, &expiry); err != nil {
			continue
		}

		if name == "" || host == "" || path == "" || value == nil {
			continue
		}

		cookie := Cookie{
			Name:       name,
			Host:       host,
			Path:       path,
			ExpireDate: expiry,
			Value:      string(value),
		}
		cookies = append(cookies, cookie)
	}

	return cookies, nil
}
