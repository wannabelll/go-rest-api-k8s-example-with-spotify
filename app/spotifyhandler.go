package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type SpotifyArtistStats struct {
	Followers  uint   `json:"followers"` // Use uint instead of int to match the type of Followers.Count
	ArtistName string `json:"artistName"`
	Popularity int    `json:"popularity"`
	SomeJunk   string `json:"somejunk"`
}

func getArtistStats(apiToken string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Create oauth2 token from the API token string
		token := &oauth2.Token{
			AccessToken: apiToken,
		}

		// Create an OAuth2 token source
		tokenSource := oauth2.StaticTokenSource(token)

		// Create an http client using the token source
		client := spotify.NewClient(oauth2.NewClient(context.Background(), tokenSource))

		// Get the artist ID from the URL parameters
		artistID := ps.ByName("artistID")

		// Make sure the artist ID is valid (it should be a Spotify artist ID without the "spotify:artist:" prefix)
		// For example: "3MtohoQqvZFtmRTwzp0xSH"
		if artistID == "" {
			http.Error(w, "artistID is required", http.StatusBadRequest)
			return
		}

		// Get artist data from Spotify by calling GetArtist with the proper ID
		artist, err := client.GetArtist(spotify.ID(artistID))
		if err != nil {
			http.Error(w, fmt.Sprintf("could not get artist data: %v", err), http.StatusInternalServerError)
			return
		}

		// Create the response struct with artist statistics
		stats := SpotifyArtistStats{
			Followers:  artist.Followers.Count, // No need to cast, Followers.Count is of type uint
			ArtistName: artist.Name,
			Popularity: artist.Popularity,
			SomeJunk:   "just an extra field",
		}

		// Write the response as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
		}
	}
}
