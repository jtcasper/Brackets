package main

import (
	"git.jacobcasper.com/brackets/client"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"net/url"
)

func main() {

	client, err := client.Get()
	if err != nil {
		log.Fatal("Could not get client: ", err.Error())
	}

	_, page, err := client.FeaturedPlaylists()

	for _, playlist := range page.Playlists {
		tracks, err := client.GetPlaylistTracks(playlist.ID)
		if err != nil {
			log.Printf("Couldn't retrieve playlist %s.", string(playlist.ID))
			continue
		}
		for _, trackPage := range tracks.Tracks {
			for _, artist := range trackPage.Track.Artists {
				http.PostForm("http://localhost:8080/artist/add", url.Values{"id": {string(artist.ID)}})
			}
		}
	}
}
