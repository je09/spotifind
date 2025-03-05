package spotifind

const SpotifyOwnerID = "spotify"

// searchCriteria is an indicator, that a playlist probably has a curator contacts.
// It can be a false positive, if an artist just leaves his handle in there.
// But it could be idea to contact them, anyway.
var searchCriteria = []string{
	"@",
	"submithub",
	"sbmt",
	"submit",
	"submissions",
	"ko-fi",
}

// SpotifindAuth - just a set of info that needed to authenticate with Spotify.
type SpotifindAuth struct {
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	RedirectURI  string `yaml:"redirectUri"`
}

type SpotifyOwner struct {
	ID          string
	DisplayName string
}

type Playlist struct {
	ID             string
	Name           string
	Collaborative  bool
	Description    string
	Styles         []string
	ExternalURLs   map[string]string
	Owner          SpotifyOwner
	TracksTotal    int
	FollowersTotal int
	Contacts       []string
	Region         string
}

// Progress - a struct that contains the total amount of playlists and the amount of playlists that have been scanned.
type Progress struct {
	Total int
	Done  int
}

// PlaylistInfo - a struct that contains a playlist and progress of the scan.
// Even though we have a separate channel for progress, in case you don't want
// to set up a separate goroutine for it, you can use this field.
// However, the data may be outdated, as it's not updated as fast as the channel.
type PlaylistInfo struct {
	Playlist Playlist
	Progress Progress
}
