package client

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"os"
)

func Get() (*spotify.Client, error) {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		return nil, fmt.Errorf("client: %w", err)
	}

	client := spotify.Authenticator{}.NewClient(token)
	return &client, nil
}
