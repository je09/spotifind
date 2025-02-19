# Spotifind 🎧
Spotifind is a Go-based application designed to search Spotify playlists for contact info. 
It extracts contact information and musical styles from playlists, providing a comprehensive overview of the playlists' content.
Works pretty similar to paid services like PlaylistSupply and Distrokid's playlist engine, but free and open-source.

## Features
- **Search Playlists**: Search for playlists based on specific criteria.
- **Extract Contacts**: Extract contact information from playlist descriptions.
- **Analyze Styles**: Analyze and categorize musical styles from playlists.
- **Progress Tracking**: Track the progress of playlist processing.

# Installation
## CLI
### Windows
1. Download the latest release for Windows from the [releases page](https://github.com/je09/spotifind/releases) the name has to be "spotifind-cli-windows.exe.zip" or "spotifind-cli-windows-arm64.exe.zip" for ARM. If you do not recognize what arm64 is, you should download the first one.
2. Extract the downloaded file in your User's folder or any other folder you're comfortable with.
3. Download and save "[spotyfind2.yml](https://github.com/je09/spotifind)" file to the directory where the spotifind .exe file from the previous step is located.

### macOS
1. Download the latest release for macOS from the [releases page](https://github.com/je09/spotifind/releases) the name has to be "spotifind-cli-darwin.zip" or "spotifind-cli-darwin-arm64.zip" for Apple Silicon (M1,M2,etc.).
2. Extract the downloaded file in your User's folder or any other folder you're comfortable with.
3. Download and save "[spotyfind2.yml](https://github.com/je09/spotifind)" file to the directory where the spotifind binary file from the previous step is located.

### Linux
If you use GNU/Linux based OS, I'm not even sure if you need my guidance, but here it is anyway:
1. Download the latest release for Linux from the [releases page](https://github.com/je09/spotifind/releases) the name has to be "spotifind-cli-linux.zip" or "spotifind-cli-linux-arm64.zip" for ARM.
2. Extract the downloaded file in your User's folder or any other folder you're comfortable with.
3. Download and save "[spotyfind2.yml](https://github.com/je09/spotifind)" file to the directory where the spotifind binary file from the previous step is located.

## GUI
TODO

# Setting Up
If you're not experienced user, the process of the configuration and usage of the Spotifind can be a bit tricky, but don't worry, I'll guide you through the process. 
###### First of all, you need to setup your Spotify Developer account and create a new application to get your API credentials. *it needs to be done so spotifind can communicate with SpotifyAPI and therefore extract the data from the playlists.*.
1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/applications) and log in with your Spotify account (yes, it's safe, it's Spotify's site after all).
2. Then go to the dashboard (click on your username on the top right angle of the page).
3. Click on the "Create an App" button.
4. Choose any "App name" and "App description" you're comfortable with. For the "Redirect URIs" you can put `http://localhost:8080/callback` and then click "Add" button, which is just in the same line.
5. To answer "Which API/SDKs are you planning to use?" simply choose "Web API".
6. Read the "Developers Terms of Service" and "Spotify Developer Terms of Service" and click "I understand..." if you agree with them.
7. Click "Save".
8. Cool, you're doing great! Now click "Settings" on the right side of the page.
9. Click "View Client Secret" and then copy both "Client ID" and "Client Secret" to the "spotyfind2.yml" file we set up on the first step. Simply replace the placeholders with your credentials.
10. Save the file and you're ready to go!
11. You may include credentials from multiple applications, so if you encounter any issues, you can switch to another application. Note, that this is not intended use and can violate the Spotify API terms of service. 

**Note, do not share your credentials with anyone, as it can lead to the unauthorized access to your Spotify account!**

# Usage
## CLI
This option is for those who are comfortable with the command line interface.

### How to open and use the CLI (Command Line Interface)
#### Windows
todo

#### macOS
todo

#### Linux
todo


# Using in your own projects
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
