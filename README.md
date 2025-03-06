# Spotifind 🎧
Spotifind is a Go-based library designed to search Spotify playlists for contact info.
It extracts contact information and musical styles from playlists, providing a comprehensive overview of the playlists' content.
Works pretty similar to paid services like PlaylistSupply and Distrokid's playlist engine, but free and open-source.

If you're looking for the app, you need to go [here](https://github.com/je09/spotifind-app)

## Features
- **Search Playlists**: Search for playlists based on specific criteria.
- **Extract Contacts**: Extract contact information from playlist descriptions.
- **Analyze Styles**: Analyze and categorize musical styles from playlists.
- **Progress Tracking**: Track the progress of playlist processing.
- 
# Use in your own projects
Spotifind is designed as a library, so you can use it in your own projects.
Todo so, you can use go get command:
```bash
go get github.com/je09/spotifind
```

Then you can import the library in your project:
```go
package main

import (
	"fmt"
	"github.com/je09/spotifind"
)

func main() {
	auth := spotifind.SpotifindAuth{
		ClientID:     "client_id",
		ClientSecret: "client_secret",
	}

	s, err := spotifind.NewSpotifind(auth, false)
	if err != nil {
		panic(err)
	}

	// ch - channel for search results
	ch := make(spotifind.SpotifindChan)

	// pCh - channel for progress of the scan
	pCh := make(spotifind.ProgressChan)

	// Search for playlists on popular markets
	go func() {
		err = s.SearchPlaylistPopular(
			ch,
			pCh,
			[]string{"liquidfunk", "autonomic", "microfunk"}, // your search queries, just like the ones you'd type in the Spotify search bar
			[]string{"techno", "metal", "punk"})              // ignore these strings in description and name of the playlist
		if err != nil {
			panic(err)
		}
	}()

	// Output the progress of the scan
	go func() {
		for progress := range pCh {
			fmt.Printf("Progress: %d out of %d\n", progress.Done, progress.Total)
		}
	}()

	// Output found playlists
	for playlist := range ch {
		fmt.Printf("Playlist: %s has contacts %v\n", playlist.Playlist.Name, playlist.Playlist.Contacts)
	}
}
```

### If you want to say thank you somehow, simply [listen to my music on Spotify or anywhere you like](https://syglit.xyz), it would mean a lot to me! Everything I do is because of the passion for the music I have! Thank you!

# Attention!
All the code in this repository is for educational purposes only.
It is not intended to be used for any other purpose, as it can violate the terms of service of the Spotify API.
Please consult current version of the Spotify API terms of service before using the Spotifind.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## Acknowledgements

- [Spotify Web API](https://developer.spotify.com/documentation/web-api/)
