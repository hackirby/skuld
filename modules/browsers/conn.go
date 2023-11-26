package browsers

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func GetDBConnection(database string) (*sql.DB, error) {
	connection, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=ro&immutable=1", database))
	if err != nil {
		return nil, err
	}

	return connection, nil
}
