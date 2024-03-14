package browsers

import (
	"path/filepath"

	_ "modernc.org/sqlite"
)

func (c *Chromium) GetHistory(path string) (history []History, err error) {
	db, err := GetDBConnection(filepath.Join(path, "History"))
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT url, title, visit_count, last_visit_time FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			url, title    string
			visitCount    int
			lastVisitTime int64
		)
		if err = rows.Scan(&url, &title, &visitCount, &lastVisitTime); err != nil {
			continue
		}

		if url == "" || title == "" {
			continue
		}

		history = append(history, History{
			URL:           url,
			Title:         title,
			VisitCount:    visitCount,
			LastVisitTime: lastVisitTime,
		})

	}

	return history, nil
}

func (g *Gecko) GetHistory(path string) (history []History, err error) {
	db, err := GetDBConnection(filepath.Join(path, "places.sqlite"))
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT url, title, visit_count, last_visit_date FROM moz_places")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			url, title    string
			visitCount    int
			lastVisitTime int64
		)
		if err = rows.Scan(&url, &title, &visitCount, &lastVisitTime); err != nil {
			continue
		}

		if url == "" || title == "" {
			continue
		}

		history = append(history, History{
			URL:           url,
			Title:         title,
			VisitCount:    visitCount,
			LastVisitTime: lastVisitTime,
		})

	}

	return history, nil
}
