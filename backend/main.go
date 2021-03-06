package main

import (
	"git.jacobcasper.com/brackets/env"
	"git.jacobcasper.com/brackets/routes/artist"
	"git.jacobcasper.com/brackets/routes/genre"
	"git.jacobcasper.com/brackets/scrape/graph"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	env, err := env.New()
	if err != nil {
		log.Fatal("Could not set up Env: ", err.Error())
	}

	http.HandleFunc("/artist/", artist.Index(env))
	http.HandleFunc("/artist/genre", artist.ByGenre(env))
	http.HandleFunc("/artist/add", artist.Add(env))

	http.HandleFunc("/genre", genre.Index(env))

	go graph.Scrape(env)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
