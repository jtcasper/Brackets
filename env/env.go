package env

import (
	"git.jacobcasper.com/brackets/client"
	"git.jacobcasper.com/brackets/db"
	"github.com/zmb3/spotify"
)

type Env struct {
	Db *db.DB
	C  *spotify.Client
}

func New() (*Env, error) {
	db, err := db.New()
	if err != nil {
		return nil, err
	}

	client, err := client.Get()
	if err != nil {
		return nil, err
	}
	return &Env{Db: db, C: client}, nil
}
