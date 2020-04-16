package types

import "github.com/zmb3/spotify"

type Artist struct {
	ID         spotify.ID `json:"id"`
	Name       string     `json:"name"`
	Popularity int        `json:"popularity"`
}
