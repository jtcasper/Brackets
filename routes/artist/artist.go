package artist

import (
	"database/sql"
	"encoding/json"
	"git.jacobcasper.com/brackets/env"
	"git.jacobcasper.com/brackets/routes"
	"git.jacobcasper.com/brackets/types"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
)

func Index(env *env.Env) routes.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		artistId := r.FormValue("id")
		if artistId != "" {
			artist := types.Artist{}
			row := env.Db.Db.QueryRow("SELECT ID, NAME FROM ARTIST WHERE ID = ?", artistId)
			if err := row.Scan(&artist.ID, &artist.Name); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			b, err := json.Marshal(artist)
			if err != nil {
				log.Print(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.Write(b)
			return
		}

		rows, err := env.Db.Db.Query("SELECT ID, NAME FROM ARTIST LIMIT 20")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		artists := make([]types.Artist, 0)
		for rows.Next() {
			artist := types.Artist{}
			if err := rows.Scan(&artist.ID, &artist.Name); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			artists = append(artists, artist)
		}
		if err = rows.Err(); err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(artists)
		if err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}

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
			var genreId int64
			row := env.Db.Db.QueryRow(`
SELECT ID
FROM GENRE
WHERE NAME = lower(?)
`,
				genre)

			err := row.Scan(&genreId)
			if err == sql.ErrNoRows {
				result, err := env.Db.Db.Exec("INSERT INTO GENRE (NAME) VALUES (?)", genre)
				if err != nil {
					log.Printf("Failed to insert genre %s: %s", genre, err.Error())
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				genreId, err = result.LastInsertId()
				if err != nil {
					log.Print("Failed to retrieve last insert id: ", err.Error())
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			}

			env.Db.Db.Exec("INSERT INTO ARTIST_GENRE_XREF (ARTIST_ID, GENRE_ID) VALUES (?, ?)", artist.ID, genreId)
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func ByGenre(env *env.Env) routes.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		genreName := r.FormValue("genre_name")
		if genreName != "" {
			rows, err := env.Db.Db.Query("SELECT ID, NAME FROM ARTIST WHERE ID IN (SELECT ARTIST_ID FROM ARTIST_GENRE_XREF WHERE GENRE_ID IN (SELECT ID FROM GENRE WHERE NAME = lower(?))) LIMIT 20", genreName)
			if err != nil {
				log.Print(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			artists := make([]types.Artist, 0)
			for rows.Next() {
				artist := types.Artist{}
				if err := rows.Scan(&artist.ID, &artist.Name); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				artists = append(artists, artist)
			}
			if err = rows.Err(); err != nil {
				log.Print(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			b, err := json.Marshal(artists)
			if err != nil {
				log.Print(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.Write(b)
		}
	}
}
