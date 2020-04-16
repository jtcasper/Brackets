package graph

import (
	"git.jacobcasper.com/brackets/env"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
	"net/url"
	"time"
)

func Scrape(env *env.Env) {
infinite:
	for {
		time.Sleep(time.Second * 5)
		rows, err := env.Db.Db.Query(`
SELECT ID
FROM ARTIST
WHERE ID NOT IN (
  SELECT ARTIST_ID
	FROM SCRAPED_ARTIST
	WHERE SCRAPED == 1
)`)
		if err != nil {
			log.Print(err)
			continue infinite
		}
		defer rows.Close()

		var artistId string
		for rows.Next() {
			if err := rows.Scan(&artistId); err != nil {
				log.Print(err)
				continue infinite
			}

			artists, err := env.C.GetRelatedArtists(spotify.ID(artistId))
			if err != nil {
				log.Print(err)
				continue infinite
			}

			success := true
		postArtists:
			for _, artist := range artists {
				resp, err := http.PostForm("http://localhost:8080/artist/add", url.Values{"id": {string(artist.ID)}})
				if err != nil {
					log.Print(err)
					success = false
					continue postArtists
				}
				if resp.StatusCode != http.StatusCreated {
					success = false
				}
			}

			if success {
				env.Db.Mu.Lock()
				env.Db.Db.Exec(`
REPLACE INTO SCRAPED_ARTIST (ARTIST_ID, SCRAPED)
VALUES (?, 1)`,
					string(artistId))
				env.Db.Mu.Unlock()
			}
		}
	}
}
