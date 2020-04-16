package artist

import (
	"git.jacobcasper.com/brackets/env"
	"git.jacobcasper.com/brackets/routes"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
)

func Add(env *env.Env) routes.Handler {
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

		artist, err := env.C.GetArtist(spotify.ID(artistId))
		if err != nil {
			log.Printf("Failed to retrieve artist %s: %s", artistId, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		env.Db.Mu.Lock()
		defer env.Db.Mu.Unlock()
		env.Db.Db.Exec("INSERT INTO ARTIST (ID, NAME) VALUES (?, ?)", artist.ID, artist.Name)

		for _, genre := range artist.Genres {
			result, err := env.Db.Db.Exec("REPLACE INTO GENRE (NAME) VALUES (?)", genre)
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

			env.Db.Db.Exec("INSERT INTO ARTIST_GENRE_XREF (ARTIST_ID, GENRE_ID) VALUES (?, ?)", artist.ID, genreId)
		}
		w.WriteHeader(http.StatusCreated)
	}
}
