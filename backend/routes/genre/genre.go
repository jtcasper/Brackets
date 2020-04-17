package genre

import (
	"encoding/json"
	"git.jacobcasper.com/brackets/env"
	"git.jacobcasper.com/brackets/routes"
	"git.jacobcasper.com/brackets/types"
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

		genreName := r.FormValue("name")
		if genreName != "" {
			genre := types.Genre{}
			row := env.Db.Db.QueryRow(`
SELECT ID, NAME
FROM GENRE
WHERE NAME = lower(?)`,
				genreName,
			)
			if err := row.Scan(&genre.ID, &genre.Name); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			b, err := json.Marshal(genre)
			if err != nil {
				log.Print(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.Write(b)
			return
		}

		rows, err := env.Db.Db.Query(`
SELECT ID, NAME
FROM GENRE
LIMIT 20`,
		)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		genres := make([]types.Genre, 0)
		for rows.Next() {
			genre := types.Genre{}
			if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			genres = append(genres, genre)
		}
		if err = rows.Err(); err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(genres)
		if err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}
