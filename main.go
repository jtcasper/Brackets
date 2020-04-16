package main

import (
	"git.jacobcasper.com/brackets/client"
	"git.jacobcasper.com/brackets/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
)

type handler func(http.ResponseWriter, *http.Request)

func main() {
	lockedDb, err := db.New()
	if err != nil {
		log.Fatal("Could not open db: ", err.Error())
	}

	client, err := client.Get()
	if err != nil {
		log.Fatal("Could not get client: ", err.Error())
	}

	http.HandleFunc(
		"/artist/add",
		addArtistHandler(lockedDb, client),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addArtistHandler(db *db.DB, c *spotify.Client) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		r.ParseForm()
		artistId := r.PostForm.Get("id")

		if artistId == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		artist, err := c.GetArtist(spotify.ID(artistId))
		if err != nil {
			log.Printf("Failed to retrieve artist %s: %s", artistId, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		db.Mu.Lock()
		defer db.Mu.Unlock()
		db.Db.Exec("INSERT INTO ARTIST (ID, NAME) VALUES (?, ?)", artist.ID, artist.Name)

		for _, genre := range artist.Genres {
			result, err := db.Db.Exec("REPLACE INTO GENRE (NAME) VALUES (?)", genre)
			if err != nil {
				log.Printf("Failed to insert genre %s: %s", genre, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			genreId, err := result.LastInsertId()
			if err != nil {
				log.Print("Failed to retrieve last insert id: ", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			db.Db.Exec("INSERT INTO ARTIST_GENRE_XREF (ARTIST_ID, GENRE_ID) VALUES (?, ?)", artist.ID, genreId)
		}
		w.WriteHeader(http.StatusCreated)
	}
}
