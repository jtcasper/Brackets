package main

import (
	"git.jacobcasper.com/brackets/client"
	"git.jacobcasper.com/brackets/db"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {

	db, err := db.New()
	if err != nil {
		log.Fatal("Could not open db: ", err.Error())
	}

	client, err := client.Get()
	if err != nil {
		log.Fatal("Could not get client: ", err.Error())
	}

	_, page, err := client.FeaturedPlaylists()

	for _, playlist := range page.Playlists {
		tracks, err := client.GetPlaylistTracks(playlist.ID)
		if err != nil {
			log.Printf("Couldn't retrieve playlist %s.", string(playlist.ID))
		}
		for _, trackPage := range tracks.Tracks {
			for _, artist := range trackPage.Track.Artists {
				db.Mu.Lock()
				db.Db.Exec("INSERT INTO ARTIST (ID, NAME) VALUES (?, ?)", artist.ID, artist.Name)
				db.Mu.Unlock()
			}
		}
	}
}
