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

type Progress struct {
	Total int
	Done  int
}

type PlaylistInfo struct {
	Playlist Playlist
	Progress Progress
}
