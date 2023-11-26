package browsers

import (
	"path/filepath"
	"regexp"

	_ "modernc.org/sqlite"
)

func (c *Chromium) GetDownloads(path string) (downloads []Download, err error) {
	db, err := GetDBConnection(filepath.Join(path, "History"))
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT tab_url, target_path FROM downloads")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			url, path string
		)
		if err = rows.Scan(&url, &path); err != nil {
			continue
		}

		if url == "" || path == "" {
			continue
		}

		downloads = append(downloads, Download{
			URL:        url,
			TargetPath: path,
		})

	}

	return downloads, nil
}

func (g *Gecko) GetDownloads(path string) (downloads []Download, err error) {
	db, err := GetDBConnection(filepath.Join(path, "places.sqlite"))
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT place_id, GROUP_CONCAT(content), url, dateAdded FROM (SELECT * FROM moz_annos INNER JOIN moz_places ON moz_annos.place_id=moz_places.id) t GROUP BY place_id")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			content, url       string
			placeID, dateAdded int64
		)

		if err = rows.Scan(&placeID, &content, &url, &dateAdded); err != nil {
			continue
		}
		if url == "" || path == "" {
			continue
		}
		re := regexp.MustCompile(`file:///(.*?),`)
		result := re.FindStringSubmatch(content)
		if len(result) == 0 {
			continue
		}
		downloads = append(downloads, Download{
			URL:        url,
			TargetPath: result[1],
		})
	}

	return downloads, nil
}
