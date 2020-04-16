package main

import (
	"git.jacobcasper.com/brackets/env"
	"git.jacobcasper.com/brackets/routes/artist"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	env, err := env.New()
	if err != nil {
		log.Fatal("Could not set up Env: ", err.Error())
	}

	http.HandleFunc(
		"/artist/",
		artist.Index(env),
	)

	http.HandleFunc(
		"/artist/add",
		artist.Add(env),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
